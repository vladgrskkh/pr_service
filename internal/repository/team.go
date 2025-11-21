package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vladgrskkh/pr_service/internal/domain"
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
		INSERT INTO teams (name, members)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.DB.Exec(ctx, query, team.Name, team.Members)
	if err != nil {
		return err
	}

	return nil
}

func (r *TeamRepository) Get(name string) (*domain.Team, error) {
	query := `
		SELECT id, name, members
		FROM teams
		WHERE name = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var team domain.Team

	err := r.DB.QueryRow(ctx, query, name).Scan(&team.ID, &team.Name, &team.Members)
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
