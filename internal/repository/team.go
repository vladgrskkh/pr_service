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
	ErrDuplicateTeamName = errors.New("duplicate team name")
)

type TeamRepository struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewTeamRepository(db *pgxpool.Pool, c *trmpgx.CtxGetter) *TeamRepository {
	return &TeamRepository{
		db:     db,
		getter: c,
	}
}

func (r *TeamRepository) Insert(ctx context.Context, team *domain.Team) error {
	query := `
		INSERT INTO teams (name)
		VALUES ($1)
	`

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	_, err := conn.Exec(ctx, query, team.Name)
	if err != nil {
		var pgErr *pgconn.PgError

		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23505":
			return ErrDuplicateTeamName
		default:
			return err
		}
	}

	return nil
}

func (r *TeamRepository) Get(ctx context.Context, name string) (*domain.Team, error) {
	query := `
		SELECT name
		FROM teams
		WHERE name = $1
	`

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	var team domain.Team

	err := conn.QueryRow(ctx, query, name).Scan(&team.Name)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &team, nil
}
