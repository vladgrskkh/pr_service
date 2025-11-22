package team

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/vladgrskkh/pr_service/internal/apierrors"
	"github.com/vladgrskkh/pr_service/internal/domain"
	"github.com/vladgrskkh/pr_service/internal/repository"
	"github.com/vladgrskkh/pr_service/pkg/helpers/json"
)

type TeamCreater interface {
	CreateTeam(name string, members []*domain.User) (*domain.Team, error)
}

func NewPostTeamHandler(logger *slog.Logger, service TeamCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name    string         `json:"name"`
			Members []*domain.User `json:"members"`
		}

		err := json.Read(w, r, &input)
		if err != nil {
			apierrors.BadRequestResponse(logger, w, r, err)
			return
		}

		team, err := service.CreateTeam(input.Name, input.Members)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrDuplicateTeamName) || errors.Is(err, repository.ErrDuplicateUser):
				apierrors.BadRequestResponse(logger, w, r, err)
			default:
				logger.Info("inside db error internal")
				apierrors.ServerErrorResponse(logger, w, r, err)
			}

			return
		}

		err = json.Write(w, http.StatusCreated, json.Envelope{"team": team.Name, "members": input.Members}, nil)
		if err != nil {
			logger.Info("inside json write error internal")
			apierrors.ServerErrorResponse(logger, w, r, err)
		}
	}
}

type TeamGetter interface {
	GetTeam(name string) (*domain.Team, []string, error)
}

func NewGetTeamHandler(logger *slog.Logger, service TeamGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamName := r.URL.Query().Get("team_name")

		team, members, err := service.GetTeam(teamName)
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

		err = json.Write(w, http.StatusOK, json.Envelope{"team_name": team.Name, "members": members}, nil)
		if err != nil {
			logger.Info("inside json write error internal")
			apierrors.ServerErrorResponse(logger, w, r, err)
		}
	}
}
