package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(port, host, dbName, sslmode, user string) (*pgxpool.Pool, error) {
	var ctx context.Context = context.Background()

	var config *pgxpool.Config
	var err error
	config, err = pgxpool.ParseConfig(fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s", user, host, port, dbName, sslmode))

	if err != nil {
		log.Printf("Unable to connect to the postgres Server")
		return nil, err
	}

	var pool *pgxpool.Pool
	pool, err = pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		log.Printf("Unable to create connection pool")
		return nil, err
	}

	err = pool.Ping(ctx)

	if err != nil {
		log.Printf("Unable to ping database")
		pool.Close()
		return nil, err
	}
	log.Printf("Succesfully connected to postgres")
	return pool, nil
}
