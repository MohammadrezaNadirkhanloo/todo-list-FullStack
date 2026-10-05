package usecase

import (
	"context"

	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/authz"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/model"
)

type AuthzUsecase struct {
	cfg config.AuthzConfig
}

func NewAuthzUsecase(cfg config.AuthzConfig) *AuthzUsecase {
	return &AuthzUsecase{cfg: cfg}
}

func (u *AuthzUsecase) Enforced() bool { return u.cfg.Enforce }

func (u *AuthzUsecase) RulesFor(ctx context.Context, user model.UserModel) ([]authz.Rule, error) {
	return u.RulesForIdentity(ctx, user.ID, user.Username, nil)
}

func (u *AuthzUsecase) RulesForIdentity(
	_ context.Context,
	userID int64,
	_ string,
	_ []string,
) ([]authz.Rule, error) {
	if userID <= 0 {
		return []authz.Rule{}, nil
	}

	return []authz.Rule{
		authz.AllowMany(
			[]string{
				authz.ActionRead,
				authz.ActionCreate,
				authz.ActionUpdate,
				authz.ActionDelete,
			},
			authz.SubjectTodo,
		).Where(authz.Conditions{"userId": userID}),
	}, nil
}

func (u *AuthzUsecase) AbilityFor(ctx context.Context, user model.UserModel) (*authz.Ability, error) {
	rules, err := u.RulesFor(ctx, user)
	if err != nil {
		return nil, err
	}
	return authz.NewAbility(rules), nil
}

func (u *AuthzUsecase) AbilityForIdentity(
	ctx context.Context,
	userID int64,
	username string,
	roles []string,
) (*authz.Ability, error) {
	rules, err := u.RulesForIdentity(ctx, userID, username, roles)
	if err != nil {
		return nil, err
	}
	return authz.NewAbility(rules), nil
}