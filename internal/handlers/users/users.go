package users

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/vladgrskkh/pr_service/internal/apierrors"
	"github.com/vladgrskkh/pr_service/internal/domain"
	"github.com/vladgrskkh/pr_service/internal/repository"
	s "github.com/vladgrskkh/pr_service/internal/service"
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
				apierrors.NotFoundResponse(logger, w, r)
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

type PRsByAuthorGetter interface {
	GetPRsByAuthor(id string) ([]*domain.PR, error)
}

func NewGetPRsByAuthorHandler(logger *slog.Logger, service PRsByAuthorGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")

		prs, err := service.GetPRsByAuthor(userID)
		if err != nil {
			switch {
			case errors.Is(err, s.ErrNoPRsCreated):
				apierrors.NotFoundResponse(logger, w, r)
			default:
				apierrors.ServerErrorResponse(logger, w, r, err)
			}

			return
		}

		err = json.Write(w, http.StatusOK, json.Envelope{"user_id": userID, "pull_requests": prs}, nil)
		if err != nil {
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
			switch {
			case errors.Is(err, s.ErrNoPRsAssigned):
				apierrors.NotFoundResponse(logger, w, r)
			default:
				apierrors.ServerErrorResponse(logger, w, r, err)
			}

			return
		}

		err = json.Write(w, http.StatusOK, json.Envelope{"user_id": userID, "pull_requests": prs}, nil)
		if err != nil {
			apierrors.ServerErrorResponse(logger, w, r, err)
		}
	}
}

type MassDeactivater interface {
	MassDeactiveUsers(teamName string, users []string) error
}

func NewPostMassDeactivate(logger *slog.Logger, service MassDeactivater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			TeamName string   `json:"team_name"`
			Users    []string `json:"users"`
		}

		err := json.Read(w, r, &input)
		if err != nil {
			apierrors.BadRequestResponse(logger, w, r, err)
			return
		}

		err = service.MassDeactiveUsers(input.TeamName, input.Users)
		if err != nil {
			// poor error handling here
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				apierrors.NotFoundResponse(logger, w, r)
			case errors.Is(err, s.ErrNoActiveUsers):
				apierrors.NoCandidateResponse(logger, w, r)
			default:
				apierrors.ServerErrorResponse(logger, w, r, err)
			}

			return
		}

		err = json.Write(w, http.StatusOK, json.Envelope{"message": "success deactivated users", "users": input.Users}, nil)
		if err != nil {
			apierrors.ServerErrorResponse(logger, w, r, err)
		}
	}
}
