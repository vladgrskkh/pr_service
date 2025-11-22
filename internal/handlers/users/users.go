package users

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/vladgrskkh/pr_service/internal/apierrors"
	"github.com/vladgrskkh/pr_service/internal/domain"
	"github.com/vladgrskkh/pr_service/internal/repository"
	"github.com/vladgrskkh/pr_service/pkg/helpers/json"
)

type IsActiveSetter interface {
	SetIsActiveUser(id string, isActive bool) (*domain.User, error)
}

func NewPostSetIsActiveHandler(logger *slog.Logger, service IsActiveSetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			UserID   string `json:"user_id"`
			IsActive bool   `json:"is_active"`
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
				logger.Info("inside db error internal")
				apierrors.ServerErrorResponse(logger, w, r, err)
			}

			return
		}

		err = json.Write(w, http.StatusOK, json.Envelope{"user": user}, nil)
		if err != nil {
			logger.Info("inside json write error internal")
			apierrors.ServerErrorResponse(logger, w, r, err)
		}
	}
}

type ReviewsGetter interface {
	GetReviewByUser(id string) ([]*domain.PR, error)
}

func NewGetReviewsHandler(logger *slog.Logger, service ReviewsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")

		prs, err := service.GetReviewByUser(userID)
		if err != nil {
			apierrors.ServerErrorResponse(logger, w, r, err)
			return
		}

		err = json.Write(w, http.StatusOK, json.Envelope{"user_id": userID, "pull_requests": prs}, nil)
		if err != nil {
			apierrors.ServerErrorResponse(logger, w, r, err)
		}
	}
}
