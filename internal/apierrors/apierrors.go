package apierrors

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/vladgrskkh/pr_service/pkg/helpers/json"
)

func logError(logger *slog.Logger, r *http.Request, err error) {
	logger.Error(err.Error(),
		slog.String("request_method", r.Method),
		slog.String("request_url", r.URL.String()),
		slog.String("trace", string(debug.Stack())))
}

// ErrorResponse writes a JSON response with a provided status code and message
// to the http.ResponseWriter.
func ErrorResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	data := json.Envelope{
		"error": message,
	}

	err := json.Write(w, status, data, nil)
	if err != nil {
		logError(logger, r, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	logError(logger, r, err)

	message := "server encountered a problem and could not process your request"
	ErrorResponse(logger, w, r, http.StatusInternalServerError, message)
}
