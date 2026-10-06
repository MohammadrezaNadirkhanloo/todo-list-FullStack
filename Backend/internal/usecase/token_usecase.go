package usecase

import (
	"context"

	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/authz"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/model"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/repository"
	"github.com/MohammadrezaNadirkhanloo/internal/usecase/dto"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"github.com/MohammadrezaNadirkhanloo/pkg/security"
)

const refreshTokenBytes = 32

type TokenUsecase struct {
	tokens *security.TokenManager
	store  repository.TokenStore
	users  repository.UserRepository
	cfg    config.JWTConfig
}

func NewTokenUsecase(
	tokens *security.TokenManager,
	store repository.TokenStore,
	users repository.UserRepository,
	cfg config.JWTConfig,
) *TokenUsecase {
	return &TokenUsecase{tokens: tokens, store: store, users: users, cfg: cfg}
}

func (u *TokenUsecase) Issue(ctx context.Context, user model.User) (dto.TokenPair, error) {
	var zero dto.TokenPair

	access, err := u.tokens.Issue(user.ID, user.Username, []string{authz.RoleUser})
	if err != nil {
		return zero, apperror.Internal(err)
	}

	refreshToken, err := security.RandomToken(refreshTokenBytes)
	if err != nil {
		return zero, apperror.Internal(err)
	}

	hash := security.HashToken(refreshToken)
	if err := u.store.Save(ctx, user.ID, hash, u.cfg.RefreshTTL); err != nil {
		return zero, err
	}

	return dto.TokenPair{
		AccessToken:   access.Token,
		AccessTokenID: access.ID,
		RefreshToken:  refreshToken,
		ExpiresAt:     access.ExpiresAt,
	}, nil
}

func (u *TokenUsecase) Refresh(ctx context.Context, refreshToken string) (dto.TokenPair, error) {
	var zero dto.TokenPair

	if refreshToken == "" {
		return zero, apperror.New(apperror.CodeTokenRequired,
			"توکن تازه‌سازی ارسال نشده است. لطفاً دوباره وارد شوید.")
	}

	userID, ok, err := u.store.Consume(ctx, security.HashToken(refreshToken))
	if err != nil {
		return zero, err
	}
	if !ok {
		return zero, apperror.New(apperror.CodeTokenInvalid,
			"نشست شما معتبر نیست. لطفاً دوباره وارد شوید.")
	}

	user, err := u.users.GetByID(ctx, userID)
	if err != nil {
		return zero, err
	}

	return u.Issue(ctx, user)
}

func (u *TokenUsecase) RevokeAll(ctx context.Context, userID int64) error {
	return u.store.RevokeAll(ctx, userID)
}

func (u *TokenUsecase) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	_, _, _ = u.store.Consume(ctx, security.HashToken(refreshToken))
	return nil
}

func (u *TokenUsecase) RefreshTTLSeconds() int {
	return int(u.cfg.RefreshTTL.Seconds())
}

func (u *TokenUsecase) AccessTTLSeconds() int {
	return int(u.cfg.AccessTTL.Seconds())
}