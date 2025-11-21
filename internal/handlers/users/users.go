package users

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/vladgrskkh/pr_service/internal/apierrors"
	"github.com/vladgrskkh/pr_service/internal/domain"
	"github.com/vladgrskkh/pr_service/internal/repository"
	"github.com/vladgrskkh/pr_service/pkg/helpers/json"
)

type IsActiveSetter interface {
	SetIsActiveUser(id int64, isActive bool) (*domain.User, error)
}

func NewPostSetIsActiveHandler(logger *slog.Logger, service IsActiveSetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			UserID   int64 `json:"user_id"`
			IsActive bool  `json:"is_active"`
		}

		err := json.Read(w, r, &input)
		if err != nil {
			apierrors.BadRequestResponse(logger, w, r, err)
			return
		}

		user, err := service.SetIsActiveUser(input.UserID, input.IsActive)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				apierrors.NotFoundResponse(logger, w, r, err)
			default:
				apierrors.ServerErrorResponse(logger, w, r, err)
			}

			return
		}

		err = json.Write(w, http.StatusOK, json.Envelope{"user": user}, nil)
		if err != nil {
			apierrors.ServerErrorResponse(logger, w, r, err)
		}
	}
}

type ReviewsGetter interface {
	GetReviewByUser(id int64) ([]*domain.PR, error)
}

func NewGetReviewsHandler(logger *slog.Logger, service ReviewsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")

		id, err := strconv.ParseInt(userID, 10, 64)
		if err != nil || id < 1 {
			apierrors.BadRequestResponse(logger, w, r, fmt.Errorf("invalid id parameter"))
			return
		}

		prs, err := service.GetReviewByUser(id)
		if err != nil {
			apierrors.ServerErrorResponse(logger, w, r, err)
			return
		}

		err = json.Write(w, http.StatusOK, json.Envelope{"user_id": id, "pull_requests": prs}, nil)
		if err != nil {
			apierrors.ServerErrorResponse(logger, w, r, err)
		}
	}
}
