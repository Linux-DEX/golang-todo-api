package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	// creating database context for database operation
	ctx := context.Background()

	// parsing database connection string into pgxpool configuration
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Printf("Unable to parse DATABASE_URL : %v", err)
		return nil, err
	}

	// Create a new PostgreSQL connection pool using the parsed configuration.
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Printf("Unable to create connection pool : %v", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		log.Printf("Unable to ping database : %v", err)
		pool.Close()
		return nil, err
	}

	log.Println("Successfully connected to PostgresSQL database")
	return pool, nil
}
