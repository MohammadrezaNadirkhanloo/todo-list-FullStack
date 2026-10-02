package middleware

import (
	"github.com/MohammadrezaNadirkhanloo/pkg/appctx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderRequestID = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(HeaderRequestID)
		if !isValidRequestID(id) {
			id = uuid.NewString()
		}

		ctx := appctx.WithRequestID(c.Request.Context(), id)
		
		c.Request = c.Request.WithContext(ctx)

		c.Header(HeaderRequestID, id)

		c.Next()
	}
}

func isValidRequestID(id string) bool {
	if len(id) == 0 || len(id) > 64 {
		return false
	}
	for _, r := range id {
		isAllowed := (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_'
		if !isAllowed {
			return false
		}
	}
	return true
}