package domain

import "time"

type PR struct {
	ID                string    `json:"pull_request_id"`
	Name              string    `json:"pull_request_name"`
	AuthorID          string    `json:"author_id"`
	Status            string    `json:"status"`
	AssignedReviewers []string  `json:"assigned_reviewers"`
	CreatedAt         time.Time `json:"-"`
	MergedAt          time.Time `json:"-"`
}
