package usecase

import (
	"context"

	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/authz"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/filter"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/repository"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
)

type Mappers[E any, TCreate any, TUpdate any, TOut any] struct {
	ToEntity  func(TCreate) E
	ToChanges func(TUpdate) map[string]any
	ToOutput  func(E) TOut
}

type Guard[E any] struct {
	Subject string

	ToResource func(E) authz.Resource

	FieldColumns map[string]string
}

type BaseUsecase[E any, TCreate any, TUpdate any, TOut any] struct {
	repo    repository.CRUD[E]
	spec    filter.Spec
	limits  filter.PageLimits
	mappers Mappers[E, TCreate, TUpdate, TOut]
	guard   Guard[E]
	enforce bool
}

func NewBaseUsecase[E any, TCreate any, TUpdate any, TOut any](
	repo repository.CRUD[E],
	spec filter.Spec,
	pageCfg config.PaginationConfig,
	mappers Mappers[E, TCreate, TUpdate, TOut],
) *BaseUsecase[E, TCreate, TUpdate, TOut] {
	return &BaseUsecase[E, TCreate, TUpdate, TOut]{
		repo: repo,
		spec: spec,
		limits: filter.PageLimits{
			Default: pageCfg.DefaultPageSize,
			Max:     pageCfg.MaxPageSize,
		},
		mappers: mappers,
	}
}

func (u *BaseUsecase[E, TCreate, TUpdate, TOut]) WithGuard(
	g Guard[E],
	enforce bool,
) *BaseUsecase[E, TCreate, TUpdate, TOut] {
	u.guard = g
	u.enforce = enforce
	return u
}

func (u *BaseUsecase[E, TCreate, TUpdate, TOut]) Spec() filter.Spec { return u.spec }

func (u *BaseUsecase[E, TCreate, TUpdate, TOut]) PageLimits() filter.PageLimits { return u.limits }

func (u *BaseUsecase[E, TCreate, TUpdate, TOut]) guarded() bool {
	return u.enforce && u.guard.Subject != ""
}

func (u *BaseUsecase[E, TCreate, TUpdate, TOut]) resourceOf(e E) authz.Resource {
	if u.guard.ToResource == nil {
		return nil
	}
	return u.guard.ToResource(e)
}

func (u *BaseUsecase[E, TCreate, TUpdate, TOut]) authorize(
	ctx context.Context,
	action string,
	res authz.Resource,
) error {
	if !u.guarded() {
		return nil
	}

	ability, ok := authz.FromContext(ctx)
	if !ok {
		return apperror.Unauthorized("برای این عملیات باید وارد شوید.")
	}

	decision := ability.Can(action, u.guard.Subject, res)
	if decision.Allowed {
		return nil
	}
	if decision.Reason != "" {
		return apperror.New(apperror.CodeForbidden, decision.Reason)
	}
	return apperror.Forbidden()
}

func (u *BaseUsecase[E, TCreate, TUpdate, TOut]) scopeFor(
	ctx context.Context,
	action string,
) authz.Scope {
	if !u.guarded() {
		return authz.Unrestricted()
	}
	ability, ok := authz.FromContext(ctx)
	if !ok {
		return authz.Scope{}
	}
	return ability.ScopeFor(action, u.guard.Subject)
}

func (u *BaseUsecase[E, TCreate, TUpdate, TOut]) authorizeFields(
	ctx context.Context,
	res authz.Resource,
	changes map[string]any,
) error {
	ability, ok := authz.FromContext(ctx)
	if !ok {
		return apperror.Unauthorized("برای این عملیات باید وارد شوید.")
	}

	permitted := ability.PermittedFields(authz.ActionUpdate, u.guard.Subject, res)
	if permitted == nil {
		return nil
	}

	if len(u.guard.FieldColumns) == 0 {
		return apperror.Forbidden()
	}

	byColumn := make(map[string]string, len(u.guard.FieldColumns))
	for apiName, column := range u.guard.FieldColumns {
		byColumn[column] = apiName
	}

	for column := range changes {
		apiName, known := byColumn[column]
		if !known {
			return apperror.Forbidden()
		}
		if !containsString(permitted, apiName) {
			return apperror.New(apperror.CodeForbidden,
				"شما اجازه‌ی تغییر فیلد «"+apiName+"» را ندارید.")
		}
	}
	return nil
}

func containsString(list []string, v string) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}

func (u *BaseUsecase[E, TCreate, TUpdate, TOut]) Create(ctx context.Context, in TCreate) (TOut, error) {
	var zero TOut

	entity := u.mappers.ToEntity(in)

	if err := u.authorize(ctx, authz.ActionCreate, u.resourceOf(entity)); err != nil {
		return zero, err
	}
	if err := u.repo.Create(ctx, &entity); err != nil {
		return zero, err
	}
	return u.mappers.ToOutput(entity), nil
}

func (u *BaseUsecase[E, TCreate, TUpdate, TOut]) Update(ctx context.Context, id int64, in TUpdate) (TOut, error) {
	var zero TOut

	changes := u.mappers.ToChanges(in)
	if len(changes) == 0 {
		return zero, apperror.InvalidInput("هیچ فیلد قابل تغییری ارسال نشده است.")
	}

	if u.guarded() {
		current, err := u.repo.GetByID(ctx, id)
		if err != nil {
			return zero, err
		}
		res := u.resourceOf(current)

		if err := u.authorize(ctx, authz.ActionUpdate, res); err != nil {
			return zero, err
		}
		if err := u.authorizeFields(ctx, res, changes); err != nil {
			return zero, err
		}
	}

	if err := u.repo.Update(ctx, id, changes); err != nil {
		return zero, err
	}

	entity, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return zero, err
	}
	return u.mappers.ToOutput(entity), nil
}

func (u *BaseUsecase[E, TCreate, TUpdate, TOut]) Delete(ctx context.Context, id int64) error {
	if u.guarded() {
		current, err := u.repo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if err := u.authorize(ctx, authz.ActionDelete, u.resourceOf(current)); err != nil {
			return err
		}
	}
	return u.repo.Delete(ctx, id)
}

func (u *BaseUsecase[E, TCreate, TUpdate, TOut]) GetByID(ctx context.Context, id int64) (TOut, error) {
	var zero TOut

	entity, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return zero, err
	}
	if err := u.authorize(ctx, authz.ActionRead, u.resourceOf(entity)); err != nil {
		return zero, err
	}
	return u.mappers.ToOutput(entity), nil
}

func (u *BaseUsecase[E, TCreate, TUpdate, TOut]) GetByQuery(
	ctx context.Context,
	q filter.Query,
) (*filter.PagedList[TOut], error) {
	q.Normalize(u.limits)
	if err := q.Validate(u.spec); err != nil {
		return nil, err
	}

	page, err := u.repo.GetByQuery(ctx, q, u.scopeFor(ctx, authz.ActionRead))
	if err != nil {
		return nil, err
	}
	return filter.MapPagedList(page, u.mappers.ToOutput), nil
}