package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vladgrskkh/pr_service/internal/domain"
)

type PullRequestRepo struct {
	DB *pgxpool.Pool
}

func NewPullRequestRepo(db *pgxpool.Pool) *PullRequestRepo {
	return &PullRequestRepo{
		DB: db,
	}
}

func (r *PullRequestRepo) Insert(pr *domain.PR) error {
	query := `
		INSERT INTO pull_requests (id, name, author_id, status, assigned_reviewers)
		VALUES ($1, $2, $3, $4, $5)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.AssignedReviewers}

	_, err := r.DB.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *PullRequestRepo) Update(pr *domain.PR) error {
	query := `
		UPDATE pull_requests
		SET status = $1
		WHERE id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	commandTag, err := r.DB.Exec(ctx, query, pr.Status, pr.ID)
	if err != nil {
		return err
	}

	rowsAffected := commandTag.RowsAffected()

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (r *PullRequestRepo) GetAllForUser(userID int64) ([]*domain.PR, error) {
	query := `
		SELECT id, name, author_id, status, assigned_reviewers
		FROM pull_requests
		WHERE author_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var prs []*domain.PR

	for rows.Next() {
		var pr domain.PR

		err := rows.Scan(
			&pr.ID,
			&pr.Name,
			&pr.AuthorID,
			&pr.Status,
			&pr.AssignedReviewers,
		)
		if err != nil {
			return nil, err
		}

		prs = append(prs, &pr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}
