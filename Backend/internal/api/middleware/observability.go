package middleware

import (
	"context"
	"time"

	"github.com/MohammadrezaNadirkhanloo/internal/api/response"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"github.com/gin-gonic/gin"
)

// func Metrics(m *metrics.Registry) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		start := time.Now()

// 		c.Next()

// 		path := c.FullPath()
// 		if path == "" {
// 			path = "unmatched"
// 		}

// 		status := strconv.Itoa(c.Writer.Status())
// 		method := c.Request.Method

// 		m.HTTPRequests.WithLabelValues(method, path, status).Inc()
// 		m.HTTPDuration.WithLabelValues(method, path, status).
// 			Observe(time.Since(start).Seconds())
// 	}
// }

func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if d <= 0 {
			c.Next()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if ctx.Err() != nil && !c.Writer.Written() {
			response.Fail(c, apperror.New(apperror.CodeTimeout,
				"processing the request took too long. Please try again."))
		}
	}
}
