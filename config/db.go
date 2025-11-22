package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set in environment")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf(" Failed to parse database URL: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf(" Unable to ping DB: %v", err)
	}

	DB = pool
	fmt.Println(" Connected to PostgreSQL")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println(" Database connection closed")
	}
}
