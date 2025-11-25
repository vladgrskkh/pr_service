package pr

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

type PRMerger interface {
	MergePullReq(id string) (*domain.PR, error)
}

func NewPostMergeHandler(logger *slog.Logger, service PRMerger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			PullRequestID string `json:"pull_request_id"`
		}

		err := json.Read(w, r, &input)
		if err != nil {
			apierrors.BadRequestResponse(logger, w, r, err)
			return
		}

		pr, err := service.MergePullReq(input.PullRequestID)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				apierrors.NotFoundResponse(logger, w, r)
			default:
				apierrors.ServerErrorResponse(logger, w, r, err)
			}

			return
		}

		err = json.Write(w, http.StatusOK, json.Envelope{"pull_request": pr}, nil)
		if err != nil {
			apierrors.ServerErrorResponse(logger, w, r, err)
		}
	}
}

type PullReqCreater interface {
	CreatePullReq(id, name, authorID string) (*domain.PR, error)
}

func NewPostPullReqHandler(logger *slog.Logger, service PullReqCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			PullReqID   string `json:"pull_request_id"`
			PullReqName string `json:"pull_request_name"`
			AuthodID    string `json:"author_id"`
		}

		err := json.Read(w, r, &input)
		if err != nil {
			apierrors.BadRequestResponse(logger, w, r, err)
			return
		}

		pr, err := service.CreatePullReq(input.PullReqID, input.PullReqName, input.AuthodID)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrDuplicatePullReqID):
				apierrors.PullReqExistsResponse(logger, w, r)
			case errors.Is(err, s.ErrValidatePR):
				apierrors.BadRequestResponse(logger, w, r, err)
			case errors.Is(err, repository.ErrRecordNotFound):
				apierrors.NotFoundResponse(logger, w, r)
			default:
				apierrors.ServerErrorResponse(logger, w, r, err)
			}

			return
		}

		err = json.Write(w, http.StatusCreated, json.Envelope{"pr": pr}, nil)
		if err != nil {
			apierrors.ServerErrorResponse(logger, w, r, err)
		}
	}
}

type PullReqReassigner interface {
	ReassignReviewer(prID, userID string) (*domain.PR, string, error)
}

func NewPostReassignHandler(logger *slog.Logger, service PullReqReassigner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			PullReqID string `json:"pull_request_id"`
			OldUserID string `json:"old_user_id"`
		}

		err := json.Read(w, r, &input)
		if err != nil {
			apierrors.BadRequestResponse(logger, w, r, err)
			return
		}

		pr, replacement, err := service.ReassignReviewer(input.PullReqID, input.OldUserID)
		if err != nil {
			switch {
			case errors.Is(err, s.ErrMergedPRChange):
				apierrors.PullReqMergedResponse(logger, w, r)
			case errors.Is(err, s.ErrUserNotAssigned):
				apierrors.UserNotAssignedResponse(logger, w, r)
			case errors.Is(err, s.ErrNoCandidate):
				apierrors.NoCandidateResponse(logger, w, r)
			case errors.Is(err, repository.ErrRecordNotFound):
				apierrors.NotFoundResponse(logger, w, r)
			default:
				apierrors.ServerErrorResponse(logger, w, r, err)
			}

			return
		}

		err = json.Write(w, http.StatusOK, json.Envelope{"pr": pr, "replaced_by": replacement}, nil)
		if err != nil {
			apierrors.ServerErrorResponse(logger, w, r, err)
		}
	}
}
