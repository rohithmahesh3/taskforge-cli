package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_RestoreIssuePayloadUsesVersion(t *testing.T) {
	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		require.Equal(t, "/api/v1/workspaces/ws/projects/p-1/work-items/i-1/restore/", r.URL.Path)
		raw, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(raw, &body))
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "i-1"})
	}))
	defer server.Close()

	client := &Client{HTTPClient: &http.Client{Timeout: DefaultTimeout}, BaseURL: server.URL, APIKey: "k", Workspace: "ws"}
	_, err := client.RestoreIssue("p-1", "i-1", 3)
	require.NoError(t, err)
	assert.EqualValues(t, 3, body["version"])
	_, hasLegacy := body["version_num"]
	assert.False(t, hasLegacy)
}

func TestClient_RestoreCommentPayloadUsesVersionNumber(t *testing.T) {
	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		require.Equal(t, "/api/v1/workspaces/ws/projects/p-1/work-items/i-1/comments/c-1/restore/", r.URL.Path)
		raw, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(raw, &body))
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "c-1", "comment_stripped": "ok"})
	}))
	defer server.Close()

	client := &Client{HTTPClient: &http.Client{Timeout: DefaultTimeout}, BaseURL: server.URL, APIKey: "k", Workspace: "ws"}
	_, err := client.RestoreComment("p-1", "i-1", "c-1", 2)
	require.NoError(t, err)
	assert.EqualValues(t, 2, body["version_number"])
	_, hasLegacy := body["version_num"]
	assert.False(t, hasLegacy)
}

func TestClient_UpdateIssuePayloadUsesListKeys(t *testing.T) {
	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "PATCH", r.Method)
		require.Equal(t, "/api/v1/workspaces/ws/projects/p-1/work-items/i-1/", r.URL.Path)
		raw, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(raw, &body))
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "i-1"})
	}))
	defer server.Close()

	client := &Client{HTTPClient: &http.Client{Timeout: DefaultTimeout}, BaseURL: server.URL, APIKey: "k", Workspace: "ws"}
	_, err := client.UpdateIssue("p-1", "i-1", taskforge.UpdateIssueRequest{
		Assignees: []string{"u-1"},
		Labels:    []string{"l-1"},
	})
	require.NoError(t, err)
	_, hasAssigneesLegacy := body["assignees"]
	_, hasLabelsLegacy := body["labels"]
	assert.False(t, hasAssigneesLegacy)
	assert.False(t, hasLabelsLegacy)
	_, hasAssigneesList := body["assignees_list"]
	_, hasLabelsList := body["labels_list"]
	assert.True(t, hasAssigneesList)
	assert.True(t, hasLabelsList)
}
