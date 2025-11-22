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
	MergePullReq(id int64) (*domain.PR, error)
}

func NewPostMergeHandler(logger *slog.Logger, service PRMerger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			PullRequestID int64 `json:"pull_request_id"`
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
				apierrors.NotFoundResponse(logger, w, r, err)
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
	CreatePullReq(id int64, name string, authorID int64) (*domain.PR, error)
}

func NewPostPullReqHandler(logger *slog.Logger, service PullReqCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			PullReqId   int64  `json:"pull_request_id"`
			PullReqName string `json:"pull_request_name"`
			AuthodID    int64  `json:"author_id"`
		}

		err := json.Read(w, r, &input)
		if err != nil {
			apierrors.BadRequestResponse(logger, w, r, err)
			return
		}

		pr, err := service.CreatePullReq(input.AuthodID, input.PullReqName, input.AuthodID)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrDuplicatePullReqID):
				// change to 409
				apierrors.BadRequestResponse(logger, w, r, err)
			// not found team or author(need to change message)
			case errors.Is(err, repository.ErrRecordNotFound):
				apierrors.NotFoundResponse(logger, w, r, err)
			default:
				apierrors.ServerErrorResponse(logger, w, r, err)
			}
		}

		err = json.Write(w, http.StatusCreated, json.Envelope{"pull_request": pr}, nil)
		if err != nil {
			apierrors.ServerErrorResponse(logger, w, r, err)
		}
	}
}

type PullReqReassigner interface {
	ReassignReviewer(prID, userID int64) (*domain.PR, error)
}

func NewPostReassignHandler(logger *slog.Logger, service PullReqReassigner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			PullReqID int64 `json:"pull_request_id"`
			OldUserID int64 `json:"old_user_id"`
		}

		err := json.Read(w, r, &input)
		if err != nil {
			apierrors.BadRequestResponse(logger, w, r, err)
			return
		}

		pr, err := service.ReassignReviewer(input.PullReqID, input.OldUserID)
		if err != nil {
			// think about how i can make it more clear
			switch {
			case errors.Is(err, s.ErrMergedPRChange):
				apierrors.EditConflictResponse(logger, w, r, err)
			case errors.Is(err, s.ErrUserNotAssigned):
				apierrors.EditConflictResponse(logger, w, r, err)
			case errors.Is(err, s.ErrNoCandidate):
				apierrors.EditConflictResponse(logger, w, r, err)
			case errors.Is(err, repository.ErrRecordNotFound):
				apierrors.NotFoundResponse(logger, w, r, err)
			default:
				apierrors.ServerErrorResponse(logger, w, r, err)
			}

			return
		}

		err = json.Write(w, http.StatusOK, json.Envelope{"pull_request": pr, "replaced_by": input.OldUserID}, nil)
		if err != nil {
			apierrors.ServerErrorResponse(logger, w, r, err)
		}
	}
}
