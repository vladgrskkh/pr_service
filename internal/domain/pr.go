package domain

import "time"

type PR struct {
	ID                string     `json:"pull_request_id"`
	Name              string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status,omitempty"`
	AssignedReviewers []string   `json:"assigned_reviewers,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	MergedAt          *time.Time `json:"merged_at,omitempty"`
}
