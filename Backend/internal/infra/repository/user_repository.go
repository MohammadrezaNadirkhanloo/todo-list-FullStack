package repository

import (
	"context"

	"github.com/MohammadrezaNadirkhanloo/internal/domain/model"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/repository"
	"github.com/MohammadrezaNadirkhanloo/internal/infra/database"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
)

type userRepository struct {
	*BaseRepository[model.User]
}

var _ repository.UserRepository = (*userRepository)(nil)

func NewUserRepository(db *database.DB) repository.UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository[model.User](db, Options{
			Entity: "User",
			Spec:   model.UserSpec,
		}),
	}
}


func (r *userRepository) FindByUsername(ctx context.Context, username string) (model.User, error) {
	return r.firstBy(ctx, nil, "username = ?", username)
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return r.existsBy(ctx, "username = ?", username)
}

func (r *userRepository) UpdatePasswordHash(ctx context.Context, userID int64, hash string) error {
	res := r.DB().WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("password_hash", hash)

	if res.Error != nil {
		return database.TranslateError(res.Error, r.Entity())
	}
	if res.RowsAffected == 0 {
		return apperror.NotFound(r.Entity())
	}
	return nil
}