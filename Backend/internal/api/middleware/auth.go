package middleware

import (
	"errors"

	"github.com/MohammadrezaNadirkhanloo/internal/api/response"
	"github.com/MohammadrezaNadirkhanloo/pkg/appctx"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"github.com/MohammadrezaNadirkhanloo/pkg/security"
	"github.com/gin-gonic/gin"
)

func Authenticate(tokens *security.TokenManager, cookieName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractAccessToken(c, cookieName)
		if err != nil {
			response.Fail(c, err)
			return
		}

		claims, err := tokens.Parse(token)
		if err != nil {
			switch {
			case errors.Is(err, security.ErrTokenExpired):
				response.Fail(c, apperror.New(apperror.CodeTokenExpired,
					"your session has expired. Please log in again."))
			default:
				response.Fail(c, apperror.New(apperror.CodeTokenInvalid,
					"the provided token is invalid."))
			}
			return
		}

		ctx := appctx.WithUser(c.Request.Context(), claims.UserID, claims.Username, claims.Roles)
		ctx = withTokenID(ctx, claims.ID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func RequireRoles(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if _, ok := appctx.UserID(ctx); !ok {
			response.Fail(c, apperror.Unauthorized("you must log in to perform this operation."))
			return
		}

		if !appctx.HasRole(ctx, allowed...) {
			response.Fail(c, apperror.Forbidden())
			return
		}

		c.Next()
	}
}

func OptionalAuth(tokens *security.TokenManager, cookieName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractAccessToken(c, cookieName)
		if err != nil {
			c.Next()
			return
		}

		claims, err := tokens.Parse(token)
		if err != nil {
			c.Next()
			return
		}

		ctx := appctx.WithUser(c.Request.Context(), claims.UserID, claims.Username, claims.Roles)
		ctx = withTokenID(ctx, claims.ID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func extractAccessToken(c *gin.Context, cookieName string) (string, error) {
	token, err := SingleCookie(c.Request, cookieName)
	switch {
	case errors.Is(err, ErrCookieDuplicate):
		return "", apperror.New(apperror.CodeTokenInvalid,
			"your session is invalid. Please log in again.")
	case err != nil, token == "":
		return "", apperror.New(apperror.CodeTokenRequired,
			"you must log in to perform this operation.")
	}
	return token, nil
}
