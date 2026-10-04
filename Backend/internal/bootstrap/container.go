package bootstrap

import (
	"context"
	"fmt"

	"github.com/MohammadrezaNadirkhanloo/internal/api"
	"github.com/MohammadrezaNadirkhanloo/internal/api/handler"
	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/MohammadrezaNadirkhanloo/internal/infra/database"
)

type Container struct {
	Config *config.Config
	// Logger  logging.Logger
	// Metrics *metrics.Registry
	DB *database.DB
	// Cache   *cache.Client
	// DBLog لاگ ممیزی و اپلیکیشن در PostgreSQL. nil یعنی خاموش.
	// DBLog *dblog.Manager
	Server *api.Server

	// limiters *middleware.RateLimiters
}

func New(ctx context.Context, cfg *config.Config, version string) (*Container, error) {
	c := &Container{Config: cfg}

	db, err := database.Connect(ctx, cfg.Postgres)
	if err != nil {
		return nil, err
	}
	c.DB = db

	handlers := api.Handlers{
		Health: handler.NewHealthHandler(version),
		// Auth:     handler.NewAuthHandler(userUC, tokenUC, authzUC, csrfSigner, cfg),
		// Category: handler.NewCategoryHandler(categoryUC),
		// Product:  handler.NewProductHandler(productUC),
	}
	router, err := api.NewRouter(api.RouterDeps{
		Config: cfg,
		// Logger:   log,
		// Metrics:  c.Metrics,
		// Tokens:   tokenManager,
		// CSRF:     csrfSigner,
		// Limiters: c.limiters,
		// Ability:  authzUC,
		Handlers: handlers,
	})
	if err != nil {
		// c.closeQuietly(ctx)
		return nil, fmt.Errorf("bootstrap: ساخت روتر ناموفق بود: %w", err)
	}
	c.Server = api.NewServer(router, cfg)
	return c, nil
}
