package api

import (
	"github.com/gin-gonic/gin"

	"github.com/MohammadrezaNadirkhanloo/internal/api/handler"
	"github.com/MohammadrezaNadirkhanloo/internal/api/middleware"
	"github.com/MohammadrezaNadirkhanloo/internal/api/response"
	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/authz"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"github.com/MohammadrezaNadirkhanloo/pkg/security"
)

type Handlers struct {
	Health *handler.HealthHandler
	Auth   *handler.AuthHandler
	Todo   *handler.TodoHandler
}

type RouterDeps struct {
	Config   *config.Config
	Tokens   *security.TokenManager
	CSRF     *security.CSRFSigner
	Limiters *middleware.RateLimiters
	Ability  middleware.AbilityResolver
	Handlers Handlers
}

func NewRouter(deps RouterDeps) (*gin.Engine, error) {
	gin.SetMode(deps.Config.Server.RunMode)

	r := gin.New()

	if err := r.SetTrustedProxies(deps.Config.Server.TrustedProxies); err != nil {
		return nil, err
	}

	r.Use(
		middleware.RequestID(),
		gin.Recovery(),
		gin.Logger(),
		middleware.SecureHeaders(deps.Config),
		middleware.CORS(deps.Config),
		middleware.CrossOriginGuard(deps.Config),
		middleware.MaxBodySize(deps.Config.Server.MaxBodyBytes),
		middleware.Timeout(deps.Config.Server.RequestTimeout),
	)

	registerSystemRoutes(r, deps)
	registerAPIRoutes(r, deps)

	r.NoRoute(func(c *gin.Context) {
		response.Fail(c, apperror.NotFound("مسیر"))
	})
	r.NoMethod(func(c *gin.Context) {
		response.Fail(c, apperror.New(apperror.CodeInvalidInput,
			"این متد برای مسیر مورد نظر پشتیبانی نمی‌شود."))
	})

	return r, nil
}

func registerSystemRoutes(r *gin.Engine, deps RouterDeps) {
	r.GET("/healthz", deps.Handlers.Health.Live)
}

type routeGroups struct {
	Anonymous *gin.RouterGroup
	Protected *gin.RouterGroup
}

func newRouteGroups(r *gin.Engine, deps RouterDeps) routeGroups {
	cfg := deps.Config

	base := r.Group(cfg.Server.APIBasePath)
	base.Use(deps.Limiters.Limit(deps.Limiters.Global))

	accessCookie := middleware.AccessCookieName(cfg.Server.InsecureCookies)

	authenticated := []gin.HandlerFunc{
		middleware.Authenticate(deps.Tokens, accessCookie),
		middleware.RequireCSRFToken(deps.CSRF, cfg),
		middleware.LoadAbility(deps.Ability),
	}

	return routeGroups{
		Anonymous: base.Group(""),
		Protected: base.Group("", authenticated...),
	}
}

func registerAPIRoutes(r *gin.Engine, deps RouterDeps) {
	g := newRouteGroups(r, deps)
	limiters := deps.Limiters
	h := deps.Handlers

	authOpen := g.Anonymous.Group("/auth")
	{
		strict := limiters.Limit(limiters.Auth)

		authOpen.POST("/register", strict, h.Auth.Register)
		authOpen.POST("/login", strict, h.Auth.Login)
		authOpen.POST("/refresh", strict, h.Auth.Refresh)
		authOpen.POST("/logout", h.Auth.Logout)
	}

	authClosed := g.Protected.Group("/auth")
	{
		authClosed.GET("/me", h.Auth.Me)
		authClosed.POST("/logout-all", h.Auth.LogoutAll)
	}

	g.Protected.GET("/me", h.Auth.Me)

	todos := g.Protected.Group("/todos")
	handler.RegisterCRUD(todos, h.Todo.CRUD, authz.SubjectTodo, middleware.Authorize)
}