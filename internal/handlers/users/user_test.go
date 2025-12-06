package users

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vladgrskkh/pr_service/internal/domain"
	"github.com/vladgrskkh/pr_service/internal/handlers/users/mocks"
	"github.com/vladgrskkh/pr_service/internal/repository"
	s "github.com/vladgrskkh/pr_service/internal/service"
)

func TestNewPostSetIsActiveHandler(t *testing.T) {
	r := chi.NewRouter()

	service := mocks.NewIsActiveSetter(t)

	user := &domain.User{
		ID:       "user1",
		Name:     "User One",
		TeamName: "team1",
		IsActive: true,
	}
	service.On("SetIsActiveUser", "user1", true).Return(user, nil).Once()
	service.On("SetIsActiveUser", "user2", false).Return(nil, repository.ErrRecordNotFound).Once()

	logger := slog.New(slog.DiscardHandler)
	r.Post("/users/set-active", NewPostSetIsActiveHandler(logger, service))

	t.Run("success", func(t *testing.T) {
		input := `{"user_id": "user1", "is_active": true}`
		req, err := http.NewRequest("POST", "/users/set-active", bytes.NewReader([]byte(input)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response struct {
			User *domain.User `json:"user"`
		}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user1", response.User.ID)
		assert.True(t, response.User.IsActive)
	})

	t.Run("not found", func(t *testing.T) {
		input := `{"user_id": "user2", "is_active": false}`
		req, err := http.NewRequest("POST", "/users/set-active", bytes.NewReader([]byte(input)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		type errorResponse struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		var data map[string]errorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &data)
		assert.NoError(t, err)
		assert.Equal(t, "NOT_FOUND", data["error"].Code)
	})
}

func TestNewGetReviewsHandler(t *testing.T) {
	r := chi.NewRouter()

	service := mocks.NewReviewsGetter(t)

	pr := &domain.PR{
		ID: "pr1",
	}
	prs := []*domain.PR{pr}
	service.On("GetReviewByUser", "user1").Return(prs, nil).Once()
	service.On("GetReviewByUser", "user2").Return(nil, s.ErrNoPRsAssigned).Once()

	logger := slog.New(slog.DiscardHandler)
	r.Get("/users/reviews", NewGetReviewsHandler(logger, service))

	t.Run("success", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/reviews?user_id=user1", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response struct {
			UserID       string       `json:"user_id"`
			PullRequests []*domain.PR `json:"pull_requests"`
		}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user1", response.UserID)
		assert.Len(t, response.PullRequests, 1)
		assert.Equal(t, "pr1", response.PullRequests[0].ID)
	})

	t.Run("no prs assigned", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/reviews?user_id=user2", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		type errorResponse struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		var data map[string]errorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &data)
		assert.NoError(t, err)
		assert.Equal(t, "NOT_FOUND", data["error"].Code)
	})
}

func TestNewPostMassDeactivate(t *testing.T) {
	r := chi.NewRouter()

	service := mocks.NewMassDeactivater(t)

	users := []string{"user1", "user2"}
	service.On("MassDeactiveUsers", "team1", users).Return(nil).Once()
	service.On("MassDeactiveUsers", "team2", []string{"user3"}).Return(repository.ErrRecordNotFound).Once()

	logger := slog.New(slog.DiscardHandler)
	r.Post("/users/mass-deactivate", NewPostMassDeactivate(logger, service))

	t.Run("success", func(t *testing.T) {
		input := `{"team_name": "team1", "users": ["user1", "user2"]}`
		req, err := http.NewRequest("POST", "/users/mass-deactivate", bytes.NewReader([]byte(input)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response struct {
			Message string   `json:"message"`
			Users   []string `json:"users"`
		}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success deactivated users", response.Message)
		assert.Equal(t, users, response.Users)
	})

	t.Run("not found", func(t *testing.T) {
		input := `{"team_name": "team2", "users": ["user3"]}`
		req, err := http.NewRequest("POST", "/users/mass-deactivate", bytes.NewReader([]byte(input)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		type errorResponse struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		var data map[string]errorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &data)
		assert.NoError(t, err)
		assert.Equal(t, "NOT_FOUND", data["error"].Code)
	})
}
