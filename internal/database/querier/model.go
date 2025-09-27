package querier

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Model struct {
	Queries
}

func NewModel(db DBTX) *Model {
	return &Model{Queries{db}}
}

func (m *Model) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return m.db.Exec(ctx, sql, args...)
}

func (m *Model) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return m.db.Query(ctx, sql, args...)
}

func (m *Model) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return m.db.QueryRow(ctx, sql, args...)
}
