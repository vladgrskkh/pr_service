package team

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vladgrskkh/pr_service/internal/domain"
	"github.com/vladgrskkh/pr_service/internal/handlers/team/mocks"
	"github.com/vladgrskkh/pr_service/internal/repository"
)

func TestNewGetTeamHandler(t *testing.T) {
	r := chi.NewRouter()

	service := mocks.NewTeamGetter(t)

	team := &domain.Team{Name: "team1"}
	users := []string{"u1"}
	service.On("GetTeam", "team1").Return(team, users, nil)

	service.On("GetTeam", "team2").Return(nil, nil, repository.ErrRecordNotFound)

	logger := slog.New(slog.DiscardHandler)
	r.Get("/team/get", NewGetTeamHandler(logger, service))

	t.Run("success", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/team/get?team_name=team1", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, rec.Code, http.StatusOK)

		var response struct {
			TeamName string   `json:"team_name"`
			Members  []string `json:"members"`
		}

		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "team1", response.TeamName)

	})

	t.Run("not found", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/team/get?team_name=team2", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		type response struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}

		var data map[string]response
		err = json.Unmarshal(rec.Body.Bytes(), &data)
		assert.NoError(t, err)

		assert.Equal(t, "NOT_FOUND", data["error"].Code)
	})
}
