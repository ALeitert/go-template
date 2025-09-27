package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/risingwavelabs/eris"

	"template/internal/config"
	"template/internal/database/querier"
)

var (
	pool  *pgxpool.Pool
	model querier.Querier
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

func TxExec(ctx context.Context, tx func(ctx context.Context, model querier.Querier) error) (err error) {
	txCon, err := pool.Begin(ctx)
	if err != nil {
		return eris.Wrap(err, "failed to start transaction")
	}

	defer func() {
		rbErr := txCon.Rollback(ctx)
		if errors.Is(rbErr, pgx.ErrTxClosed) {
			rbErr = nil
		}
		err = errors.Join(err, rbErr)
	}()

	err = tx(ctx, querier.New(txCon))
	if err != nil {
		return err
	}

	return txCon.Commit(ctx)
}

func TxQuery[T any](ctx context.Context, tx func(ctx context.Context, model querier.Querier) (T, error)) (result T, err error) {
	var zeroT T
	txCon, err := pool.Begin(ctx)
	if err != nil {
		return zeroT, eris.Wrap(err, "failed to start transaction")
	}

	defer func() {
		rbErr := txCon.Rollback(ctx)
		if errors.Is(rbErr, pgx.ErrTxClosed) {
			rbErr = nil
		}
		err = errors.Join(err, rbErr)

		if err != nil {
			result = zeroT
		}
	}()

	result, err = tx(ctx, querier.New(txCon))
	if err != nil {
		return zeroT, err
	}

	err = txCon.Commit(ctx)
	if err != nil {
		return zeroT, err
	}

	return result, nil
}
