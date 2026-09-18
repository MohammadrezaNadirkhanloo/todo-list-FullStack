package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MohammadrezaNadirkhanloo/todo-list-FullStack/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() error {
	var err error
	DB, err = pgxpool.New(context.Background(), config.Config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.Ping(ctx); err != nil {
		return fmt.Errorf("error connecting to db: %w", err)
	}

	log.Println("Connected DB ✅")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}