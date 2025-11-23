package healthcheck

import (
	"log/slog"
	"net/http"

	"github.com/vladgrskkh/pr_service/internal/apierrors"
	"github.com/vladgrskkh/pr_service/pkg/helpers/json"
)

func New(logger *slog.Logger, env, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := json.Envelope{
			"status":  "available",
			"env":     env,
			"version": version,
		}

		err := json.Write(w, http.StatusOK, data, nil)
		if err != nil {
			logger.Error("error writing json", slog.Any("error", err))
			apierrors.ServerErrorResponse(logger, w, r, err)
		}
	}
}
