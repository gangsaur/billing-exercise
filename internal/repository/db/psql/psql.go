package psql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Psql struct {
	pool *pgxpool.Pool
}

func NewPsql(dsn string) (*Psql, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &Psql{
		pool: pool,
	}, nil
}

func (p *Psql) CloseConnection() {
	p.pool.Close()
}

// Transaction
// Still leaky since we return the pgx.Tx, adjust later

func (p *Psql) Begin(ctx context.Context) (pgx.Tx, error) {
	return p.pool.Begin(ctx)
}

func (p *Psql) Commit(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (p *Psql) Rollback(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}
