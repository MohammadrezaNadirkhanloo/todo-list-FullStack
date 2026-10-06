package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/MohammadrezaNadirkhanloo/internal/api/response"
	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/MohammadrezaNadirkhanloo/pkg/appctx"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"github.com/MohammadrezaNadirkhanloo/pkg/security"
	"github.com/gin-gonic/gin"
)

func csrfFailed() error {
	return apperror.New(apperror.CodeCSRFFailed,
		"request was not verified. Please reload the page and try again.")
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	}
	return false
}

func CrossOriginGuard(cfg *config.Config) gin.HandlerFunc {
	if !cfg.CSRF.Enabled {
		return func(c *gin.Context) { c.Next() }
	}

	trusted := make(map[string]struct{}, len(cfg.Cors.AllowOrigins))
	for _, o := range cfg.Cors.AllowOrigins {
		if o == "*" {
			continue
		}
		trusted[strings.ToLower(o)] = struct{}{}
	}

	return func(c *gin.Context) {
		if isSafeMethod(c.Request.Method) || isSameOriginRequest(c.Request, trusted) {
			c.Next()
			return
		}
		response.Fail(c, csrfFailed())
	}
}

func isSameOriginRequest(r *http.Request, trusted map[string]struct{}) bool {
	origin := r.Header.Get("Origin")
	if origin != "" {
		if _, ok := trusted[strings.ToLower(origin)]; ok {
			return true
		}
	}

	switch r.Header.Get("Sec-Fetch-Site") {
	case "":
	case "same-origin", "none":
		return true
	default:
		return false
	}

	if origin == "" {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil || u.Host == "" {
		return false
	}
	return strings.EqualFold(u.Host, r.Host)
}

func RequireCSRFToken(signer *security.CSRFSigner, cfg *config.Config) gin.HandlerFunc {
	if !cfg.CSRF.Enabled {
		return func(c *gin.Context) { c.Next() }
	}
	cookieName := CSRFCookieName(cfg.Server.InsecureCookies)

	return func(c *gin.Context) {
		if isSafeMethod(c.Request.Method) {
			c.Next()
			return
		}

		ctx := c.Request.Context()
		if _, authenticated := appctx.UserID(ctx); !authenticated {
			c.Next()
			return
		}

		sessionID, ok := tokenIDFrom(ctx)
		if !ok {
			response.Fail(c, csrfFailed())
			return
		}

		header := c.GetHeader(CSRFHeaderName)
		cookie, err := SingleCookie(c.Request, cookieName)
		if header == "" || err != nil || cookie == "" ||
			!security.ConstantTimeEquals(header, cookie) ||
			signer.Verify(header, sessionID) != nil {
			response.Fail(c, csrfFailed())
			return
		}

		c.Next()
	}
}
