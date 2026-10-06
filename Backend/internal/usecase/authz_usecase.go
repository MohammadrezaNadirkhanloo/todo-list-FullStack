package usecase

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/authz"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/model"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/repository"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
)

type AuthzUsecase struct {
	perms repository.PermissionRepository
	cfg   config.AuthzConfig
}

func NewAuthzUsecase(perms repository.PermissionRepository, cfg config.AuthzConfig) *AuthzUsecase {
	return &AuthzUsecase{perms: perms, cfg: cfg}
}

func (u *AuthzUsecase) Enforced() bool { return u.cfg.Enforce }

func (u *AuthzUsecase) RulesFor(ctx context.Context, user model.User) ([]authz.Rule, error) {
	return u.RulesForIdentity(ctx, user.ID, user.Username, user.RoleNames())
}

func (u *AuthzUsecase) RulesForIdentity(
	ctx context.Context,
	userID int64,
	username string,
	roles []string,
) ([]authz.Rule, error) {
	rules := make([]authz.Rule, 0, 8)

	if u.cfg.Mode == config.AuthzRoles || u.cfg.Mode == config.AuthzBoth {
		rules = append(rules, authz.RulesForRoles(roles)...)
	}

	if userID > 0 && (u.cfg.Mode == config.AuthzRules || u.cfg.Mode == config.AuthzBoth) {
		dynamic, err := u.dynamicRules(ctx, userID, username)
		if err != nil {
			return nil, err
		}
		rules = append(rules, dynamic...)
	}

	return rules, nil
}

func (u *AuthzUsecase) AbilityFor(ctx context.Context, user model.User) (*authz.Ability, error) {
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

func (u *AuthzUsecase) dynamicRules(
	ctx context.Context,
	userID int64,
	username string,
) ([]authz.Rule, error) {
	records, err := u.perms.RulesForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	out := make([]authz.Rule, 0, len(records))
	for _, rec := range records {
		rule, err := u.toRule(rec, userID, username)
		if err != nil {
			continue
		}
		out = append(out, rule)
	}
	return out, nil
}

func (u *AuthzUsecase) toRule(
	rec model.PermissionRule,
	userID int64,
	username string,
) (authz.Rule, error) {
	rule := authz.Rule{
		Action:   authz.Strings(splitCSV(rec.Actions)),
		Subject:  authz.Strings(splitCSV(rec.Subjects)),
		Fields:   splitCSV(rec.Fields),
		Inverted: rec.Inverted,
		Reason:   rec.Reason,
	}

	if s := strings.TrimSpace(rec.Conditions); s != "" && s != "null" {
		var conds authz.Conditions
		if err := json.Unmarshal([]byte(s), &conds); err != nil {
			return authz.Rule{}, err
		}
		rule.Conditions = interpolate(conds, userID, username)
	}

	if err := rule.Validate(); err != nil {
		return authz.Rule{}, err
	}
	return rule, nil
}

const (
	phUserID       = "${user.id}"
	phUserUsername = "${user.username}"
)

func interpolate(c authz.Conditions, userID int64, username string) authz.Conditions {
	if len(c) == 0 {
		return c
	}

	out := make(authz.Conditions, len(c))
	for k, v := range c {
		out[k] = interpolateValue(v, userID, username)
	}
	return out
}

func interpolateValue(v any, userID int64, username string) any {
	switch typed := v.(type) {

	case string:
		switch typed {
		case phUserID:
			return userID
		case phUserUsername:
			return username
		default:
			return typed
		}

	case map[string]any:
		out := make(map[string]any, len(typed))
		for k, item := range typed {
			out[k] = interpolateValue(item, userID, username)
		}
		return out

	case []any:
		out := make([]any, len(typed))
		for i, item := range typed {
			out[i] = interpolateValue(item, userID, username)
		}
		return out

	default:
		return v
	}
}

func splitCSV(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func (u *AuthzUsecase) AssignProfile(ctx context.Context, userID, profileID int64) error {
	if userID <= 0 || profileID <= 0 {
		return apperror.InvalidInput("user ID and profile ID must be valid.")
	}
	return u.perms.AssignToUser(ctx, userID, profileID)
}

func (u *AuthzUsecase) RevokeProfile(ctx context.Context, userID, profileID int64) error {
	if userID <= 0 || profileID <= 0 {
		return apperror.InvalidInput("user ID and profile ID must be valid.")
	}
	return u.perms.RevokeFromUser(ctx, userID, profileID)
}