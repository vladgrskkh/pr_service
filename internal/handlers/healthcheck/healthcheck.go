package healthcheck

import (
	"log/slog"
	"net/http"

	"github.com/vladgrskkh/pr_service/pkg/helpers/json"
)

func errorResponse(logger *slog.Logger, w http.ResponseWriter, status int, message interface{}) {
	data := json.Envelope{
		"error": message,
	}

	err := json.Write(w, status, data, nil)
	if err != nil {
		logger.Error("error writing json", slog.Any("error", err))
		w.WriteHeader(500)
	}
}

func New(logger *slog.Logger, env, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := json.Envelope{
			"status":  "avaliable",
			"env":     env,
			"version": version,
		}

		err := json.Write(w, http.StatusOK, data, nil)
		if err != nil {
			logger.Error("error writing json", slog.Any("error", err))
			errorResponse(logger, w, http.StatusInternalServerError, "server encountered a problem and could not process your request")
		}
	}
}
