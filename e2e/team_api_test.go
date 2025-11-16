package e2e_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const baseURL = "http://localhost:8080"

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string `json:"team_name"`
	Members  []User `json:"members"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func postJSON(t *testing.T, url string, payload interface{}) (*http.Response, []byte) {
	body, err := json.Marshal(payload)
	assert.NoError(t, err)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	return resp, respBody
}

func getJSON(t *testing.T, url string) (*http.Response, []byte) {
	resp, err := http.Get(url)
	assert.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	return resp, respBody
}

func TestTeamAPI(t *testing.T) {
	teamPayload := Team{
		TeamName: "payments",
		Members: []User{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
		},
	}

	resp, body := postJSON(t, baseURL+"/team/add", teamPayload)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdTeam struct {
		Team Team `json:"team"`
	}
	err := json.Unmarshal(body, &createdTeam)
	assert.NoError(t, err)
	assert.Equal(t, teamPayload.TeamName, createdTeam.Team.TeamName)
	assert.Len(t, createdTeam.Team.Members, 2)

	resp, body = getJSON(t, baseURL+"/team/get?team_name=payments")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fetchedTeam Team
	err = json.Unmarshal(body, &fetchedTeam)
	assert.NoError(t, err)
	assert.Equal(t, teamPayload.TeamName, fetchedTeam.TeamName)
	assert.Len(t, fetchedTeam.Members, 2)

	userPayload := map[string]interface{}{
		"user_id":   "u2",
		"is_active": false,
	}
	resp, body = postJSON(t, baseURL+"/users/setIsActive", userPayload)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedUser struct {
		User User `json:"user"`
	}
	err = json.Unmarshal(body, &updatedUser)
	assert.NoError(t, err)
	assert.Equal(t, false, updatedUser.User.IsActive)
	assert.Equal(t, "u2", updatedUser.User.UserID)

	resp, body = postJSON(t, baseURL+"/team/add", teamPayload)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var errResp ErrorResponse
	err = json.Unmarshal(body, &errResp)
	assert.NoError(t, err)
	assert.Equal(t, "TEAM_EXISTS", errResp.Error.Code)

	resp, body = getJSON(t, baseURL+"/team/get?team_name=nonexistent")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	err = json.Unmarshal(body, &errResp)
	assert.NoError(t, err)

	nonExistentUser := map[string]interface{}{
		"user_id":   "u999",
		"is_active": true,
	}
	resp, body = postJSON(t, baseURL+"/users/setIsActive", nonExistentUser)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	err = json.Unmarshal(body, &errResp)
	assert.NoError(t, err)
}
