package middleware

import (
	"golang.org/x/time/rate"

	"github.com/MohammadrezaNadirkhanloo/internal/api/response"
	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"github.com/MohammadrezaNadirkhanloo/pkg/ratelimit"
	"github.com/gin-gonic/gin"
)

type RateLimiters struct {
	Global *ratelimit.KeyLimiter
	Auth   *ratelimit.KeyLimiter
	Otp    *ratelimit.KeyLimiter

	enabled bool
}

func NewRateLimiters(cfg config.RateLimitConfig) *RateLimiters {
	build := func(rps float64, burst int) *ratelimit.KeyLimiter {
		return ratelimit.New(ratelimit.Options{
			Rate:    rate.Limit(rps),
			Burst:   burst,
			IdleTTL: cfg.IdleTTL,
			MaxKeys: cfg.MaxKeys,
		})
	}

	return &RateLimiters{
		Global:  build(cfg.GlobalRPS, cfg.GlobalBurst),
		Auth:    build(cfg.AuthRPS, cfg.AuthBurst),
		// Otp:     build(cfg.OtpRPS, cfg.OtpBurst),
		enabled: cfg.Enabled,
	}
}

func (r *RateLimiters) Close() {
	r.Global.Close()
	r.Auth.Close()
	r.Otp.Close()
}

func (r *RateLimiters) Limit(limiter *ratelimit.KeyLimiter) gin.HandlerFunc {
	if !r.enabled {
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		if !limiter.Allow(c.ClientIP()) {
			c.Header("Retry-After", "60")
			response.Fail(c, apperror.New(apperror.CodeRateLimited,
				"you have exceeded the allowed number of requests. Please try again later."))
			return
		}
		c.Next()
	}
}
