package psql

import (
	"context"
	"errors"

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
