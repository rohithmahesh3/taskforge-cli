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

func TestClient_PatchIssueUsesServerPatchSchema(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "/api/v1/workspaces/test-workspace/projects/proj-1/work-items/issue-1/", r.URL.Path)

		var req taskforge.PatchRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		require.Len(t, req.Patches, 1)
		assert.Equal(t, "replace", req.Patches[0].Op)
		assert.Equal(t, "description", req.Patches[0].Field)
		assert.Equal(t, "old", req.Patches[0].Old)
		assert.Equal(t, "new", req.Patches[0].New)

		_ = json.NewEncoder(w).Encode(taskforge.Issue{ID: "issue-1", SequenceID: 1})
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		Workspace:  "test-workspace",
	}

	_, err := client.PatchIssue("proj-1", "issue-1", []taskforge.PatchOp{
		{Op: "replace", Field: "description", Old: "old", New: "new"},
	})
	require.NoError(t, err)
}

func TestClient_PatchIssueConflictReturnsStructuredAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Patch conflicts",
			"conflicts": []map[string]any{
				{"field": "description", "message": "Precondition text not found"},
			},
		})
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		Workspace:  "test-workspace",
	}

	_, err := client.PatchIssue("proj-1", "issue-1", []taskforge.PatchOp{
		{Op: "replace", Field: "description", Old: "old", New: "new"},
	})
	require.Error(t, err)

	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, http.StatusConflict, apiErr.StatusCode)
	assert.Equal(t, "Patch conflicts", apiErr.Message)
	assert.Contains(t, string(apiErr.Conflicts), "Precondition text not found")
}
