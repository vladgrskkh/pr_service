package repository

import (
	"context"
	"errors"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vladgrskkh/pr_service/internal/domain"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateUser  = errors.New("duplicate user")
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

	err := conn.QueryRow(ctx, query, isActive, id).Scan(&user.ID, &user.Name, &user.TeamName, &user.IsActive)
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

func (r *UsersRepo) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, name, team_name, is_active)
		VALUES ($1, $2, $3, $4)
	`

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	args := []interface{}{user.ID, user.Name, user.TeamName, user.IsActive}
	_, err := conn.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrDuplicateUser
			}
		}

		return err
	}

	return nil
}

func (r *UsersRepo) Get(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, name, team_name, is_active
		FROM users
		WHERE id = $1
	`

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	var user domain.User

	args := []interface{}{&user.ID, &user.Name, &user.TeamName, &user.IsActive}
	err := conn.QueryRow(ctx, query, id).Scan(args...)
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

func (r *UsersRepo) GetAllForTeam(ctx context.Context, teamName string) ([]*domain.User, error) {
	query := `
		SELECT id, is_active
		FROM users
		WHERE team_name = $1
	`

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	rows, err := conn.Query(ctx, query, teamName)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var members []*domain.User

	for rows.Next() {
		var member domain.User

		err := rows.Scan(&member.ID, &member.IsActive)
		if err != nil {
			return nil, err
		}

		members = append(members, &member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, nil

}

func (r *UsersRepo) UpdateDeactivateForTeam(ctx context.Context, teamName string, users []string) error {
	query := `
		UPDATE users
		SET is_active = false
		WHERE team_name = $1 AND id = ANY($2)	
	`

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	commandTag, err := conn.Exec(ctx, query, teamName, users)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return nil
}
