package usecase

import (
	"context"
	"strings"

	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/authz"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/model"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/repository"
	"github.com/MohammadrezaNadirkhanloo/internal/usecase/dto"
	"github.com/MohammadrezaNadirkhanloo/pkg/appctx"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
)

type TodoUsecase struct {
	*BaseUsecase[model.Todo, dto.CreateTodoInput, dto.UpdateTodoInput, dto.TodoOutput]
	repo repository.TodoRepository
}

func NewTodoUsecase(
	repo repository.TodoRepository,
	pageCfg config.PaginationConfig,
) *TodoUsecase {
	base := NewBaseUsecase(
		repository.CRUD[model.Todo](repo),
		model.TodoSpec,
		pageCfg,
		Mappers[model.Todo, dto.CreateTodoInput, dto.UpdateTodoInput, dto.TodoOutput]{
			ToEntity:  dto.ToTodoModel,
			ToChanges: dto.TodoChanges,
			ToOutput:  dto.ToTodoOutput,
		},
	)

	base.WithGuard(Guard[model.Todo]{
		Subject: authz.SubjectTodo,
		ToResource: func(t model.Todo) authz.Resource {
			return authz.Resource{"userId": t.UserID}
		},
	}, true)

	return &TodoUsecase{BaseUsecase: base, repo: repo}
}

func (u *TodoUsecase) Create(ctx context.Context, in dto.CreateTodoInput) (dto.TodoOutput, error) {
	var zero dto.TodoOutput

	id, ok := appctx.UserID(ctx)
	if !ok {
		return zero, apperror.Unauthorized("برای این عملیات باید وارد شوید.")
	}

	todo := dto.ToTodoModel(in)
	todo.Title = strings.TrimSpace(todo.Title)
	if todo.Title == "" {
		return zero, apperror.InvalidInput("عنوان کار نمی‌تواند خالی باشد.")
	}
	todo.UserID = id

	if err := u.authorize(ctx, authz.ActionCreate, u.resourceOf(todo)); err != nil {
		return zero, err
	}
	if err := u.repo.Create(ctx, &todo); err != nil {
		return zero, err
	}
	return dto.ToTodoOutput(todo), nil
}

func (u *TodoUsecase) Update(ctx context.Context, id int64, in dto.UpdateTodoInput) (dto.TodoOutput, error) {
	if in.Title != nil {
		t := strings.TrimSpace(*in.Title)
		if t == "" {
			return dto.TodoOutput{}, apperror.InvalidInput("عنوان کار نمی‌تواند خالی باشد.")
		}
		in.Title = &t
	}
	return u.BaseUsecase.Update(ctx, id, in)
}