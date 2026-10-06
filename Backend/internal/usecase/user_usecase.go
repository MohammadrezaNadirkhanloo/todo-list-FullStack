package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/MohammadrezaNadirkhanloo/internal/domain/model"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/repository"
	"github.com/MohammadrezaNadirkhanloo/internal/usecase/dto"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"github.com/MohammadrezaNadirkhanloo/pkg/security"
)

var errInvalidCredentials = apperror.Unauthorized("نام کاربری یا رمز عبور نادرست است.")

type UserUsecase struct {
	users  repository.UserRepository
	tokens *TokenUsecase
	hasher *security.Argon2Hasher
}

func NewUserUsecase(
	users repository.UserRepository,
	tokens *TokenUsecase,
	hasher *security.Argon2Hasher,
) *UserUsecase {
	return &UserUsecase{users: users, tokens: tokens, hasher: hasher}
}

func (u *UserUsecase) Register(ctx context.Context, in dto.RegisterInput) (dto.UserOutput, error) {
	var zero dto.UserOutput

	username := normalize(in.Username)

	exists, err := u.users.ExistsByUsername(ctx, username)
	if err != nil {
		return zero, err
	}
	if exists {
		return zero, apperror.Conflict("این نام کاربری قبلاً ثبت شده است.")
	}

	hash, err := u.hasher.Hash(in.Password)
	if err != nil {
		return zero, apperror.Internal(err)
	}

	user := model.User{
		Username:     username,
		PasswordHash: hash,
		Enabled:      true, // ← پیشنهاد می‌کنم این خط رو هم اضافه کنی
	}
	if err := u.users.Create(ctx, &user, []string{model.RoleUser}); err != nil {
		return zero, err
	}
	return dto.ToUserOutput(user), nil
}

func (u *UserUsecase) Login(ctx context.Context, in dto.LoginInput) (dto.TokenPair, error) {
	var zero dto.TokenPair

	user, err := u.users.FindByUsername(ctx, normalize(in.Username))
	if err != nil {
		if apperror.CodeOf(err) == apperror.CodeNotFound {
			u.hasher.DummyVerify()
			return zero, errInvalidCredentials
		}
		return zero, err
	}

	needsRehash, err := u.hasher.Verify(in.Password, user.PasswordHash)
	if err != nil {
		if errors.Is(err, security.ErrMismatch) {
			return zero, errInvalidCredentials
		}
		return zero, apperror.Internal(err)
	}

	if needsRehash {
		u.upgradeHash(ctx, user.ID, in.Password)
	}

	return u.tokens.Issue(ctx, user)
}

func (u *UserUsecase) Profile(ctx context.Context, userID int64) (dto.UserOutput, error) {
	user, err := u.Find(ctx, userID)
	if err != nil {
		return dto.UserOutput{}, err
	}
	return dto.ToUserOutput(user), nil
}

func (u *UserUsecase) Find(ctx context.Context, userID int64) (model.User, error) {
	return u.users.GetByID(ctx, userID)
}

func (u *UserUsecase) upgradeHash(ctx context.Context, userID int64, plainPassword string) {
	newHash, err := u.hasher.Hash(plainPassword)
	if err != nil {
		return
	}
	_ = u.users.UpdatePasswordHash(ctx, userID, newHash)
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
