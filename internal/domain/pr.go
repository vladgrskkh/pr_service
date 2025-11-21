package domain

import "time"

type PR struct {
	ID                int64     `json:"pull_request_id"`
	Name              string    `json:"pull_request_name"`
	AuthorID          int64     `json:"author_id"`
	Status            string    `json:"status"`
	AssignedReviewers []User    `json:"assigned_reviewers"`
	CreatedAt         time.Time `json:"-"`
	MergedAt          time.Time `json:"-"`
}
