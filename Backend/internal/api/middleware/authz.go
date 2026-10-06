package middleware

import (
	"context"

	"github.com/MohammadrezaNadirkhanloo/internal/api/response"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/authz"
	"github.com/MohammadrezaNadirkhanloo/pkg/appctx"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"github.com/gin-gonic/gin"
)

type AbilityResolver interface {
	AbilityForIdentity(
		ctx context.Context,
		userID int64,
		username string,
		roles []string,
	) (*authz.Ability, error)
}

func LoadAbility(resolver AbilityResolver) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userID, authenticated := appctx.UserID(ctx)

		var (
			username string
			roles    []string
		)

		if authenticated {
			username, _ = appctx.Username(ctx)
			roles, _ = appctx.Roles(ctx)
		} else {
			roles = []string{authz.RoleGuest}
		}

		ability, err := resolver.AbilityForIdentity(ctx, userID, username, roles)
		if err != nil {
			// log.Error(ctx, "failed to build access rules",
			// 	logging.Cat(logging.CategoryAuth),
			// 	logging.F("user_id", userID),
			// 	logging.Err(err),
			// )
			response.Fail(c, apperror.Internal(err))
			return
		}

		c.Request = c.Request.WithContext(authz.NewContext(ctx, ability))
		c.Next()
	}
}

func Authorize(action, subject string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if _, ok := appctx.UserID(ctx); !ok {
			response.Fail(c, apperror.Unauthorized("you must log in to perform this operation."))
			return
		}

		ability, ok := authz.FromContext(ctx)
		if !ok {
			response.Fail(c, apperror.Forbidden())
			return
		}

		decision := ability.Can(action, subject, nil)
		if !decision.Allowed {
			if decision.Reason != "" {
				response.Fail(c, apperror.New(apperror.CodeForbidden, decision.Reason))
				return
			}
			response.Fail(c, apperror.Forbidden())
			return
		}

		c.Next()
	}
}
