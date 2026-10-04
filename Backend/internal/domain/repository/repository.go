package repository

import (
	"context"
	"time"

	"github.com/MohammadrezaNadirkhanloo/internal/domain/filter"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/model"
)

type Reader[T any] interface {
	GetByID(ctx context.Context, id int64) (T, error)

	GetByQuery(ctx context.Context, q filter.Query, scope authz.Scope) (*filter.PagedList[T], error)

	Exists(ctx context.Context, id int64) (bool, error)
}

type Writer[T any] interface {
	Create(ctx context.Context, entity *T) error

	Update(ctx context.Context, id int64, changes map[string]any) error

	Delete(ctx context.Context, id int64) error
}

type CRUD[T any] interface {
	Reader[T]
	Writer[T]
}

type TodoRepository interface {
	CRUD[model.Todo]
}

type UserRepository interface {
	Reader[model.User]

	Create(ctx context.Context, user *model.User) error

	FindByUsername(ctx context.Context, username string) (model.User, error)

	ExistsByUsername(ctx context.Context, username string) (bool, error)

	UpdatePasswordHash(ctx context.Context, userID int64, hash string) error
}

type TokenStore interface {
	Save(ctx context.Context, userID int64, tokenHash string, ttl time.Duration) error

	Consume(ctx context.Context, tokenHash string) (userID int64, ok bool, err error)

	RevokeAll(ctx context.Context, userID int64) error
}
