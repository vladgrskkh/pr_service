package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vladgrskkh/pr_service/internal/domain"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
)

var (
	ErrDuplicatePullReqID = errors.New("duplicate pull request id")
)

type PullRequestRepo struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewPullRequestRepo(db *pgxpool.Pool, c *trmpgx.CtxGetter) *PullRequestRepo {
	return &PullRequestRepo{
		db:     db,
		getter: c,
	}
}

func (r *PullRequestRepo) Insert(ctx context.Context, pr *domain.PR) error {
	query := `
		INSERT INTO pull_requests (id, name, author_id, status, assigned_reviewers)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	args := []interface{}{pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.AssignedReviewers}

	err := conn.QueryRow(ctx, query, args...).Scan(&pr.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError

		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23505":
			return ErrDuplicatePullReqID
		default:
			return err
		}
	}

	return nil
}

func (r *PullRequestRepo) UpdateReviewers(ctx context.Context, pr *domain.PR) error {
	query := `
		UPDATE pull_requests
		SET assigned_reviewers = $1
		WHERE id = $2
		RETURNING created_at, merged_at
	`

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	err := conn.QueryRow(ctx, query, pr.AssignedReviewers, pr.ID).Scan(&pr.CreatedAt, &pr.MergedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *PullRequestRepo) UpdateStatus(ctx context.Context, pr *domain.PR) error {
	query := `
		UPDATE pull_requests
		SET status = $1
		WHERE id = $2
		RETURNING name, author_id, assigned_reviewers, created_at, merged_at
	`

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	err := conn.QueryRow(ctx, query, pr.Status, pr.ID).Scan(&pr.Name, &pr.AuthorID, &pr.AssignedReviewers, &pr.CreatedAt, &pr.MergedAt)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (r *PullRequestRepo) GetAllForUser(ctx context.Context, userID string) ([]*domain.PR, error) {
	query := `
		SELECT id, name, author_id, status, assigned_reviewers
		FROM pull_requests
		WHERE author_id = $1
	`

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	rows, err := conn.Query(ctx, query, userID)
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

func (r *PullRequestRepo) GetByID(ctx context.Context, id string) (*domain.PR, error) {
	query := `
		SELECT id, name, author_id, status, assigned_reviewers
		FROM pull_requests
		WHERE id = $1
	`

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	var pr domain.PR

	err := conn.QueryRow(ctx, query, id).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.AssignedReviewers)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &pr, nil
}
