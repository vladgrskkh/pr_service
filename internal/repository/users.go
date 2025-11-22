package repository

import (
	"context"
	"errors"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vladgrskkh/pr_service/internal/domain"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type UsersRepo struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewUsersRepo(db *pgxpool.Pool, c *trmpgx.CtxGetter) *UsersRepo {
	return &UsersRepo{
		db:     db,
		getter: c,
	}
}

func (r *UsersRepo) SetIsActive(ctx context.Context, id string, isActive bool) (*domain.User, error) {
	query := `
		UPDATE users
		SET is_active = $1
		WHERE id = $2
		RETURNING id, name, team_name, is_active
	`

	var user domain.User

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	err := conn.QueryRow(ctx, query, id, isActive).Scan(&user.ID, &user.Name, &user.TeamName, &user.IsActive)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
