package repository

import (
	"context"

	"github.com/MohammadrezaNadirkhanloo/internal/domain/authz"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/filter"
	"github.com/MohammadrezaNadirkhanloo/internal/infra/database"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"gorm.io/gorm"
)

type Options struct {
	Entity       string
	Spec         filter.Spec
	Preloads     []string
	ListPreloads []string
}

type BaseRepository[T any] struct {
	db   *database.DB
	opts Options
}

func NewBaseRepository[T any](db *database.DB, opts Options) *BaseRepository[T] {
	if len(opts.ListPreloads) == 0 {
		opts.ListPreloads = opts.Preloads
	}
	return &BaseRepository[T]{db: db, opts: opts}
}

func (r *BaseRepository[T]) DB() *database.DB { return r.db }

func (r *BaseRepository[T]) Entity() string { return r.opts.Entity }

func applyPreloads(db *gorm.DB, preloads []string) *gorm.DB {
	for _, p := range preloads {
		db = db.Preload(p)
	}
	return db
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return database.TranslateError(err, r.opts.Entity)
	}
	return nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, id int64, changes map[string]any) error {
	if len(changes) == 0 {
		return apperror.InvalidInput("no fields provided for update.")
	}

	var entity T
	res := r.db.WithContext(ctx).
		Model(&entity).
		Where("id = ?", id).
		Updates(changes)

	if res.Error != nil {
		return database.TranslateError(res.Error, r.opts.Entity)
	}
	if res.RowsAffected == 0 {
		return apperror.NotFound(r.opts.Entity)
	}
	return nil
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id int64) error {
	var entity T
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity)

	if res.Error != nil {
		return database.TranslateError(res.Error, r.opts.Entity)
	}
	if res.RowsAffected == 0 {
		return apperror.NotFound(r.opts.Entity)
	}
	return nil
}

func (r *BaseRepository[T]) GetByID(ctx context.Context, id int64) (T, error) {
	var entity T

	db := applyPreloads(r.db.WithContext(ctx), r.opts.Preloads)
	if err := db.Where("id = ?", id).First(&entity).Error; err != nil {
		return entity, database.TranslateError(err, r.opts.Entity)
	}
	return entity, nil
}

func (r *BaseRepository[T]) Exists(ctx context.Context, id int64) (bool, error) {
	return r.existsBy(ctx, "id = ?", id)
}

func (r *BaseRepository[T]) GetByQuery(
	ctx context.Context,
	q filter.Query,
	scope authz.Scope,
) (*filter.PagedList[T], error) {
	var entity T

	scoped, err := database.ApplyScope(r.db.WithContext(ctx), scope, r.opts.Spec)
	if err != nil {
		return nil, err
	}

	countTx, err := database.ApplyFilters(scoped.Model(&entity), q, r.opts.Spec)
	if err != nil {
		return nil, err
	}

	var totalRows int64
	if err := countTx.Count(&totalRows).Error; err != nil {
		return nil, database.TranslateError(err, r.opts.Entity)
	}

	if totalRows == 0 {
		return filter.NewPagedList([]T{}, 0, q.Page, q.PageSize), nil
	}

	listTx, err := database.ApplyQuery(
		applyPreloads(scoped, r.opts.ListPreloads).Model(&entity),
		q, r.opts.Spec,
	)
	if err != nil {
		return nil, err
	}

	items := make([]T, 0, q.PageSize)
	if err := listTx.Find(&items).Error; err != nil {
		return nil, database.TranslateError(err, r.opts.Entity)
	}

	return filter.NewPagedList(items, totalRows, q.Page, q.PageSize), nil
}

func (r *BaseRepository[T]) existsBy(ctx context.Context, condition string, args ...any) (bool, error) {
	var entity T
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&entity).
		Where(condition, args...).
		Limit(1).
		Count(&count).Error; err != nil {
		return false, database.TranslateError(err, r.opts.Entity)
	}
	return count > 0, nil
}

func (r *BaseRepository[T]) firstBy(
	ctx context.Context,
	preloads []string,
	condition string,
	args ...any,
) (T, error) {
	var entity T

	db := applyPreloads(r.db.WithContext(ctx), preloads)
	if err := db.Where(condition, args...).First(&entity).Error; err != nil {
		return entity, database.TranslateError(err, r.opts.Entity)
	}
	return entity, nil
}
