package bootstrap

import (
	"context"
	"fmt"

	"github.com/MohammadrezaNadirkhanloo/internal/api"
	"github.com/MohammadrezaNadirkhanloo/internal/api/handler"
	"github.com/MohammadrezaNadirkhanloo/internal/api/middleware"
	"github.com/MohammadrezaNadirkhanloo/internal/api/validation"
	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/MohammadrezaNadirkhanloo/internal/infra/cache"
	"github.com/MohammadrezaNadirkhanloo/internal/infra/database"
	"github.com/MohammadrezaNadirkhanloo/internal/infra/repository"
	"github.com/MohammadrezaNadirkhanloo/internal/usecase"
	"github.com/MohammadrezaNadirkhanloo/pkg/security"
)

type Container struct {
	Config *config.Config
	DB     *database.DB
	Cache  *cache.Client
	Server *api.Server

	limiters *middleware.RateLimiters
}

func New(ctx context.Context, cfg *config.Config, version string) (*Container, error) {
	c := &Container{Config: cfg}

	if err := validation.Register(cfg.Password); err != nil {
		return nil, fmt.Errorf("bootstrap: %w", err)
	}

	db, err := database.Connect(ctx, cfg.Postgres)
	if err != nil {
		return nil, err
	}
	c.DB = db

	redis, err := cache.Connect(ctx, cfg.Redis)
	if err != nil {
		c.closeQuietly(ctx)
		return nil, err
	}
	c.Cache = redis

	tokenManager, err := security.NewTokenManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.Issuer,
		cfg.JWT.Audience,
		cfg.JWT.AccessTTL,
	)
	if err != nil {
		c.closeQuietly(ctx)
		return nil, fmt.Errorf("bootstrap: %w", err)
	}

	csrfSigner, err := security.NewCSRFSigner(cfg.JWT.AccessSecret)
	if err != nil {
		c.closeQuietly(ctx)
		return nil, fmt.Errorf("bootstrap: %w", err)
	}

	hasher := security.NewArgon2Hasher(security.Argon2Params{
		Memory:      cfg.Password.Argon2Memory,
		Iterations:  cfg.Password.Argon2Iterations,
		Parallelism: cfg.Password.Argon2Parallelism,
		SaltLength:  16,
		KeyLength:   32,
	})

	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)
	tokenStore := repository.NewTokenStore(redis)
	permissionRepo := repository.NewPermissionRepository(db)

	tokenUC := usecase.NewTokenUsecase(tokenManager, tokenStore, userRepo, cfg.JWT)
	todoUC := usecase.NewTodoUsecase(todoRepo, cfg.Pagination)
	authzUC := usecase.NewAuthzUsecase(permissionRepo, cfg.Authz)
	userUC := usecase.NewUserUsecase(userRepo, tokenUC, hasher)

	c.limiters = middleware.NewRateLimiters(cfg.RateLimit)

	handlers := api.Handlers{
		Health: handler.NewHealthHandler(version),
		Auth:   handler.NewAuthHandler(userUC, tokenUC, authzUC, csrfSigner, cfg),
		Todo:   handler.NewTodoHandler(todoUC),
	}

	router, err := api.NewRouter(api.RouterDeps{
		Config:   cfg,
		Tokens:   tokenManager,
		CSRF:     csrfSigner,
		Limiters: c.limiters,
		Ability:  authzUC,
		Handlers: handlers,
	})
	if err != nil {
		c.closeQuietly(ctx)
		return nil, fmt.Errorf("bootstrap: ساخت روتر ناموفق بود: %w", err)
	}

	c.Server = api.NewServer(router, cfg)
	return c, nil
}

func (c *Container) Close(ctx context.Context) error {
	var firstErr error

	record := func(_ string, err error) {
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if c.limiters != nil {
		c.limiters.Close()
	}
	if c.Cache != nil {
		record("redis", c.Cache.Close())
	}
	if c.DB != nil {
		record("postgres", c.DB.Close())
	}

	return firstErr
}

func (c *Container) closeQuietly(ctx context.Context) {
	_ = c.Close(ctx)
}
