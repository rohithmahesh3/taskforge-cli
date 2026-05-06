package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListIssueTypesUsesWorkItemTypesEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/workspaces/test-workspace/projects/test-project/work-item-types/", r.URL.Path)

		err := json.NewEncoder(w).Encode(struct {
			Results []taskforge.IssueType `json:"results"`
		}{
			Results: []taskforge.IssueType{{ID: "type-1", Name: "Bug"}},
		})
		require.NoError(t, err)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	types, err := client.ListIssueTypes("test-project")
	require.NoError(t, err)
	assert.Len(t, types, 1)
	assert.Equal(t, "type-1", types[0].ID)
}

func TestCreateIssueTypeUsesWorkItemTypesEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/workspaces/test-workspace/projects/test-project/work-item-types/", r.URL.Path)

		var req taskforge.CreateIssueTypeRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "Bug", req.Name)

		err = json.NewEncoder(w).Encode(taskforge.IssueType{ID: "type-1", Name: req.Name})
		require.NoError(t, err)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	issueType, err := client.CreateIssueType("test-project", taskforge.CreateIssueTypeRequest{Name: "Bug"})
	require.NoError(t, err)
	assert.Equal(t, "type-1", issueType.ID)
}
