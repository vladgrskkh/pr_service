package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vladgrskkh/pr_service/internal/domain"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type UsersRepo struct {
	db *pgxpool.Pool
}

func NewUsersRepo(db *pgxpool.Pool) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) SetIsActive(id int64) (*domain.User, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		UPDATE users
		SET is_active = true
		WHERE user_id = $1
		RETURNING user_id, username, is_active
	`

	var user domain.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.IsActive)
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
