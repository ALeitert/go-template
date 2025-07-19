package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"template/internal/config"
	"template/internal/database/querier"
)

var (
	pool  *pgxpool.Pool
	model *querier.Queries
)

const (
	QueryTimeout       = 10 * time.Second
	TransactionTimeout = 10 * time.Second
)

func Connect(ctx context.Context) error {
	var err error
	for i := 0; i < 5; i++ {
		pool, err = openDatabase(ctx)
		if err != nil {
			fmt.Printf("failed to connect to database %v, retry connecting to database...\n", err)
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	if err != nil {
		return err
	}

	model = querier.New(pool)

	return nil
}

func openDatabase(ctx context.Context) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	dsn := fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s connect_timeout=15`,
		config.C.Database.Host,
		config.C.Database.Port,
		config.C.Database.User,
		config.C.Database.Password,
		config.C.Database.Name,
	)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
