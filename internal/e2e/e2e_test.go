package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vladgrskkh/pr_service/internal/domain"
)

const baseURL = "http://localhost:8080"

func waitForAPI(t *testing.T) {
	timeout := time.After(30 * time.Second)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timeout:
			t.Fatal("API did not become ready in time")
		case <-tick:
			resp, err := http.Get(baseURL + "/healthcheck")
			if err == nil && resp.StatusCode == 200 {
				return
			}
		}
	}
}

func postJSON(t *testing.T, url string, payload interface{}, expectedCode int) []byte {
	body, err := json.MarshalIndent(payload, "", "\t")
	assert.NoError(t, err)

	resp, err := http.Post(baseURL+url, "application/json", bytes.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, expectedCode, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	return data
}

func TestTeamAndPRFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	waitForAPI(t)

	users := []domain.User{
		{
			ID:       "u1",
			Name:     "User 1",
			TeamName: "E2ETestTeam",
			IsActive: true,
		},
		{
			ID:       "u2",
			Name:     "User 2",
			TeamName: "E2ETestTeam",
			IsActive: true,
		},
		{
			ID:       "u3",
			Name:     "User 3",
			TeamName: "E2ETestTeam",
			IsActive: true,
		},
		{
			ID:       "u4",
			Name:     "User 4",
			TeamName: "E2ETestTeam",
			IsActive: true,
		},
		{
			ID:       "u5",
			Name:     "User 5",
			TeamName: "E2ETestTeam",
			IsActive: true,
		},
	}

	teamPayload := map[string]interface{}{
		"name":    "E2ETestTeam",
		"members": users,
	}

	postJSON(t, "/team/add", teamPayload, http.StatusCreated)

	prs := []domain.PR{
		{
			ID:       "pr1",
			Name:     "PR 1",
			AuthorID: "u1",
		},
		{
			ID:       "pr2",
			Name:     "PR 2",
			AuthorID: "u2",
		},
		{
			ID:       "pr3",
			Name:     "PR 3",
			AuthorID: "u3",
		},
		{
			ID:       "pr4",
			Name:     "PR 4",
			AuthorID: "u4",
		},
		{
			ID:       "pr5",
			Name:     "PR 5",
			AuthorID: "u5",
		},
	}

	for i := 0; i < 5; i++ {
		respData := postJSON(t, "/pullRequest/create", prs[i], http.StatusCreated)
		var prResp map[string]domain.PR
		err := json.Unmarshal(respData, &prResp)
		assert.NoError(t, err)

		assert.Equal(t, "OPEN", prResp["pr"].Status)
		assert.NotContains(t, prs[i].AuthorID, prResp["pr"].AssignedReviewers)
	}

	deactivatePayload := map[string]interface{}{
		"user_id":   "u3",
		"is_active": false,
	}
	postJSON(t, "/users/setIsActive", deactivatePayload, http.StatusOK)

	for _, pr := range prs {
		mergePayload := map[string]interface{}{
			"pull_request_id": pr.ID,
		}
		postJSON(t, "/pullRequest/merge", mergePayload, http.StatusOK)
	}
}
