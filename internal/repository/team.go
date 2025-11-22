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
	ErrDuplicateTeamName = errors.New("duplicate team name")
)

type TeamRepository struct {
	DB *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{
		DB: db,
	}
}

func (r *TeamRepository) Insert(team *domain.Team) error {
	query := `
		INSERT INTO teams (name)
		VALUES ($1)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.DB.Exec(ctx, query, team.Name)
	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "teams_name_key" (SQLSTATE 23505)`:
			return ErrDuplicateTeamName
		default:
			return err
		}
	}

	return nil
}

func (r *TeamRepository) Get(name string) (*domain.Team, error) {
	query := `
		SELECT name
		FROM teams
		WHERE name = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var team domain.Team

	err := r.DB.QueryRow(ctx, query, name).Scan(&team.Name)
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
