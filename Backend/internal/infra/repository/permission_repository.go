package repository

import (
	"context"

	"github.com/MohammadrezaNadirkhanloo/internal/domain/model"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/repository"
	"github.com/MohammadrezaNadirkhanloo/internal/infra/database"
)

type permissionRepository struct {
	*BaseRepository[model.PermissionProfile]
}

var _ repository.PermissionRepository = (*permissionRepository)(nil)

func NewPermissionRepository(
	db *database.DB,
	// log logging.Logger,
	// m *metrics.Registry,
) repository.PermissionRepository {
	return &permissionRepository{
		BaseRepository: NewBaseRepository[model.PermissionProfile](db, Options{
			Entity:       "PermissionProfile",
			Spec:         model.PermissionProfileSpec,
			Preloads:     []string{"Rules"},
			ListPreloads: []string{},
		}),
	}
}

func (r *permissionRepository) RulesForUser(
	ctx context.Context,
	userID int64,
) ([]model.PermissionRule, error) {
	var rules []model.PermissionRule

	err := r.db.WithContext(ctx).
		Table("permission_rules AS pr").
		Select("pr.*").
		Joins("JOIN permission_profiles AS pp ON pp.id = pr.profile_id").
		Joins("JOIN user_permission_profiles AS upp ON upp.profile_id = pp.id").
		Where("upp.user_id = ?", userID).
		Where("pp.enabled IS TRUE").
		Where("pr.deleted_at IS NULL").
		Where("pp.deleted_at IS NULL").
		Order("pp.id, pr.position, pr.id").
		Scan(&rules).Error

	if err != nil {
		return nil, database.TranslateError(err, r.opts.Entity)
	}
	return rules, nil
}

func (r *permissionRepository) FindByName(
	ctx context.Context,
	name string,
) (model.PermissionProfile, error) {
	return r.firstBy(ctx, r.opts.Preloads, "name = ?", name)
}

func (r *permissionRepository) AssignToUser(
	ctx context.Context,
	userID, profileID int64,
) error {
	err := r.db.WithContext(ctx).Exec(
		`INSERT INTO user_permission_profiles (user_id, profile_id)
		 VALUES (?, ?) ON CONFLICT DO NOTHING`,
		userID, profileID,
	).Error
	if err != nil {
		return database.TranslateError(err, r.opts.Entity)
	}
	return nil
}

func (r *permissionRepository) RevokeFromUser(
	ctx context.Context,
	userID, profileID int64,
) error {
	err := r.db.WithContext(ctx).Exec(
		`DELETE FROM user_permission_profiles
		 WHERE user_id = ? AND profile_id = ?`,
		userID, profileID,
	).Error
	if err != nil {
		return database.TranslateError(err, r.opts.Entity)
	}
	return nil
}
