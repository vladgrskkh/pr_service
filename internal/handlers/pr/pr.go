package pr

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/vladgrskkh/pr_service/internal/apierrors"
	"github.com/vladgrskkh/pr_service/internal/domain"
	"github.com/vladgrskkh/pr_service/internal/repository"
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
