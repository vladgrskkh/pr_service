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

func (r *UsersRepo) SetIsActive(id int64, isActive bool) (*domain.User, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		UPDATE users
		SET is_active = $1
		WHERE user_id = $2
		RETURNING user_id, username, team_name, is_active
	`

	var user domain.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, query, id, isActive).Scan(&user.ID, &user.Name, &user.TeamName, &user.IsActive)
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
