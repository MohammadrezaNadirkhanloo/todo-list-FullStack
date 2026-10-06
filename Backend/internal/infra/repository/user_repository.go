package repository

import (
	"context"

	"github.com/MohammadrezaNadirkhanloo/internal/domain/model"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/repository"
	"github.com/MohammadrezaNadirkhanloo/internal/infra/database"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"gorm.io/gorm"
)

type userRepository struct {
	*BaseRepository[model.User]
}

var _ repository.UserRepository = (*userRepository)(nil)

func NewUserRepository(db *database.DB) repository.UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository[model.User](db, Options{
			Entity:   "User",
			Spec:     model.UserSpec,
			Preloads: []string{"Roles"}, // ← این خط رو اضافه کنید
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

func (r *userRepository) Create(ctx context.Context, user *model.User, roleNames []string) error {
	return r.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return database.TranslateError(err, r.Entity())
		}

		if len(roleNames) == 0 {
			return nil
		}

		var roles []model.Role
		if err := tx.Where("name IN ?", roleNames).Find(&roles).Error; err != nil {
			return database.TranslateError(err, "Role")
		}
		if len(roles) != len(roleNames) {
			return apperror.New(apperror.CodeInternal,
				"یکی از نقش‌های پیش‌فرض در دیتابیس تعریف نشده است.")
		}

		if err := tx.Model(user).Association("Roles").Append(roles); err != nil {
			return database.TranslateError(err, r.Entity())
		}
		return nil
	})
}
