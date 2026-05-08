package psql

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Loan struct {
	Id                int
	Duration          int
	PrincipalAmount   int
	OutstandingAmount int
	InterestRate      float32
	UserId            int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

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
