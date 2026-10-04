package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	http *http.Server
	// log  logging.Logger
	cfg *config.Config
}

func NewServer(handler *gin.Engine, cfg *config.Config) *Server {
	return &Server{
		http: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: handler,

			ReadTimeout:       cfg.Server.ReadTimeout,
			ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
			WriteTimeout:      cfg.Server.WriteTimeout,
			IdleTimeout:       cfg.Server.IdleTimeout,
			MaxHeaderBytes:    1 << 20,
		},
		// log: log,
		cfg: cfg,
	}
}

func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		// s.log.Info(ctx, "HTTP server is running",
		// 	logging.Cat(logging.CategoryGeneral),
		// 	logging.Sub(logging.SubStartup),
		// 	logging.F("addr", s.http.Addr),
		// 	logging.F("env", string(s.cfg.Env)),
		// )
		fmt.Println("HTTP server is running")
		if err := s.http.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("api: server failed to run: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return err

	case <-ctx.Done():
		return s.shutdown()
	}
}

func (s *Server) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.ShutdownTimeout)
	defer cancel()

	// s.log.Info(ctx, "graceful shutdown started",
	// 	logging.Cat(logging.CategoryGeneral),
	// 	logging.Sub(logging.SubShutdown),
	// 	logging.F("timeout", s.cfg.Server.ShutdownTimeout.String()),
	// )
	fmt.Println("graceful shutdown started")

	if err := s.http.Shutdown(ctx); err != nil {
		// s.log.Error(ctx, "graceful shutdown did not complete; forcing connections closed",
		// 	logging.Cat(logging.CategoryGeneral),
		// 	logging.Sub(logging.SubShutdown),
		// 	logging.Err(err),
		// )
		fmt.Println("graceful shutdown did not complete; forcing connections closed")
		return s.http.Close()
	}

	// s.log.Info(ctx, "server shut down successfully",
	// 	logging.Cat(logging.CategoryGeneral),
	// 	logging.Sub(logging.SubShutdown),
	// )
	fmt.Println("server shut down successfully")
	return nil
}
