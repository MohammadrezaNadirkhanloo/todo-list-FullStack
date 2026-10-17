package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/gin-gonic/gin"
)

func CORS(cfg *config.Config) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(cfg.Cors.AllowOrigins))
	wildcard := false
	for _, o := range cfg.Cors.AllowOrigins {
		if o == "*" {
			wildcard = true
			continue
		}
		allowed[strings.ToLower(o)] = struct{}{}
	}

	allowHeaders := strings.Join([]string{
		"Content-Type", "Content-Length", "Accept", "Accept-Encoding",
		"Origin", "Cache-Control", "X-Requested-With",
		HeaderRequestID, CSRFHeaderName,
	}, ", ")

	allowMethods := strings.Join([]string{
		http.MethodGet, http.MethodPost, http.MethodPut,
		http.MethodPatch, http.MethodDelete, http.MethodOptions,
	}, ", ")

	maxAge := strconv.Itoa(int(cfg.Cors.MaxAge.Seconds()))

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if origin == "" {
			c.Next()
			return
		}

		_, isAllowed := allowed[strings.ToLower(origin)]

		switch {
		case isAllowed:
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			if cfg.Cors.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}

		case wildcard && !cfg.Cors.AllowCredentials:
			c.Header("Access-Control-Allow-Origin", "*")

		default:
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.Next()
			return
		}

		c.Header("Access-Control-Allow-Headers", allowHeaders)
		c.Header("Access-Control-Allow-Methods", allowMethods)
		c.Header("Access-Control-Expose-Headers", HeaderRequestID)
		c.Header("Access-Control-Max-Age", maxAge)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
