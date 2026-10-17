package middleware

import (
	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/gin-gonic/gin"
)

func SecureHeaders(cfg *config.Config) gin.HandlerFunc {
	isProd := cfg.IsProduction()

	return func(c *gin.Context) {
		h := c.Writer.Header()

		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; sandbox")

		if isProd {
			h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		c.Next()
	}
}
