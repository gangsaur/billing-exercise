package psql

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Id        int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Psql) GetUser(ctx context.Context, id int) (User, error) {
	sql := "SELECT id, created_at, updated_at FROM users WHERE id = $1"

	var user User
	err := p.pool.QueryRow(ctx, sql, id).Scan(
		&user.Id, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}

	return user, err
}
