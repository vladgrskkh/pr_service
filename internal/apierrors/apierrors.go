package apierrors

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/vladgrskkh/pr_service/pkg/helpers/json"
)

type errorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

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

func BadRequestResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(logger, w, r, http.StatusBadRequest, err.Error())
}

func ServerErrorResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	logError(logger, r, err)

	message := "server encountered a problem and could not process your request"
	ErrorResponse(logger, w, r, http.StatusInternalServerError, message)
}

func NotFoundResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	errMessage := &errorMessage{
		Code:    "NOT_FOUND",
		Message: "resource not found",
	}
	ErrorResponse(logger, w, r, http.StatusNotFound, errMessage)
}

func TeamExistsResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	errMessage := &errorMessage{
		Code:    "TEAM_EXISTS",
		Message: "team_name already exists",
	}
	ErrorResponse(logger, w, r, http.StatusBadRequest, errMessage)
}

func PullReqExistsResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	errMessage := &errorMessage{
		Code:    "PR_EXISTS",
		Message: "PR id already exists",
	}
	ErrorResponse(logger, w, r, http.StatusConflict, errMessage)
}

func PullReqMergedResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	errMessage := &errorMessage{
		Code:    "PR_MERGED",
		Message: "cannot reassign on merged PR",
	}
	ErrorResponse(logger, w, r, http.StatusConflict, errMessage)
}

func UserNotAssignedResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	errMessage := &errorMessage{
		Code:    "NOT_ASSIGNED",
		Message: "reviewer is not assigned to this PR",
	}
	ErrorResponse(logger, w, r, http.StatusConflict, errMessage)
}

func NoCandidateResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	errMessage := &errorMessage{
		Code:    "NO_CANDIDATE",
		Message: "no active replacement candidate in team",
	}
	ErrorResponse(logger, w, r, http.StatusConflict, errMessage)
}

func UserExistsResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	errMessage := &errorMessage{
		Code:    "USER_EXISTS",
		Message: "user with this id already exists",
	}
	ErrorResponse(logger, w, r, http.StatusBadRequest, errMessage)
}
