package api

import (
	"github.com/MohammadrezaNadirkhanloo/internal/api/handler"
	"github.com/MohammadrezaNadirkhanloo/internal/api/middleware"
	"github.com/MohammadrezaNadirkhanloo/internal/api/response"
	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Health *handler.HealthHandler
}

type RouterDeps struct {
	Config   *config.Config
	Handlers Handlers
}

func NewRouter(dep RouterDeps) (*gin.Engine, error) {
	gin.SetMode(dep.Config.Server.RunMode)
	r := gin.New()

	if err := r.SetTrustedProxies(dep.Config.Server.TrustedProxies); err != nil {
		return nil, err
	}

	r.Use(
		middleware.RequestID(),
		gin.Recovery(),
		gin.Logger(),
	)

	r.GET("/healthz", dep.Handlers.Health.Live)

	r.NoRoute(func(c *gin.Context) {
		response.Fail(c, apperror.NotFound("مسیر"))
	})
	r.NoMethod(func(c *gin.Context) {
		response.Fail(c, apperror.New(apperror.CodeInvalidInput,
			"این متد برای مسیر مورد نظر پشتیبانی نمی‌شود."))
	})
	return r, nil
}
