package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/vladgrskkh/pr_service/internal/apierrors"
)

func RecoverPanic(logger *slog.Logger) func(http.Handler) http.Handler {
	// in this hour this was the only idea how to pass logger that came to my mind :)
	// if I have time will read about best practices for this problem
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connection", "Close")

					apierrors.ServerErrorResponse(logger, w, r, fmt.Errorf("%s", err))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
