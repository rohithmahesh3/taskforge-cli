package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
)

func TestListPageCategories(t *testing.T) {
	mockCategories := []taskforge.PageCategory{
		{ID: "pc-1", Name: "Guides", SortOrder: 1},
		{ID: "pc-2", Name: "API Reference", SortOrder: 2},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/page-categories/" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		response := struct {
			Results []taskforge.PageCategory `json:"results"`
		}{Results: mockCategories}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	categories, err := client.ListPageCategories("test-project")
	if err != nil {
		t.Fatalf("ListPageCategories failed: %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categories))
	}
	if categories[0].Name != "Guides" {
		t.Errorf("Expected first category name 'Guides', got '%s'", categories[0].Name)
	}
}

func TestCreatePageCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req taskforge.CreatePageCategoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Name != "Guides" {
			t.Errorf("Expected name 'Guides', got '%s'", req.Name)
		}

		category := taskforge.PageCategory{
			ID:   "pc-new",
			Name: req.Name,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(category)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	req := taskforge.CreatePageCategoryRequest{Name: "Guides"}
	category, err := client.CreatePageCategory("test-project", req)
	if err != nil {
		t.Fatalf("CreatePageCategory failed: %v", err)
	}

	if category.Name != "Guides" {
		t.Errorf("Expected category name 'Guides', got '%s'", category.Name)
	}
}

func TestDeletePageCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/page-categories/pc-123/" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	err := client.DeletePageCategory("test-project", "pc-123")
	if err != nil {
		t.Fatalf("DeletePageCategory failed: %v", err)
	}
}

func TestListPages(t *testing.T) {
	mockPages := []taskforge.Page{
		{ID: "pg-1", Title: "Getting Started", Slug: "getting-started", CategoryID: "pc-1"},
		{ID: "pg-2", Title: "API Reference", Slug: "api-reference", CategoryID: "pc-2"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/workspaces/test-workspace/projects/test-project/pages/"
		if r.URL.Path != expectedPath {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		// Verify search query is passed through
		q := r.URL.Query().Get("q")
		if q != "" && q != "getting" {
			t.Errorf("Unexpected search query: %s", q)
		}

		response := struct {
			Results []taskforge.Page `json:"results"`
		}{Results: mockPages}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	pages, err := client.ListPages("test-project", "", "")
	if err != nil {
		t.Fatalf("ListPages failed: %v", err)
	}

	if len(pages) != 2 {
		t.Errorf("Expected 2 pages, got %d", len(pages))
	}
	if pages[0].Title != "Getting Started" {
		t.Errorf("Expected first page title 'Getting Started', got '%s'", pages[0].Title)
	}
}

func TestListPagesWithSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q != "getting" {
			t.Errorf("Expected search query 'getting', got '%s'", q)
		}

		response := struct {
			Results []taskforge.Page `json:"results"`
		}{
			Results: []taskforge.Page{
				{ID: "pg-1", Title: "Getting Started", Slug: "getting-started"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	pages, err := client.ListPages("test-project", "getting", "")
	if err != nil {
		t.Fatalf("ListPages with search failed: %v", err)
	}

	if len(pages) != 1 {
		t.Errorf("Expected 1 page, got %d", len(pages))
	}
}

func TestGetPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/pages/pg-123/" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		page := taskforge.Page{
			ID:     "pg-123",
			Title:  "Getting Started",
			Slug:   "getting-started",
			Content: "Full content here",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(page)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	page, err := client.GetPage("test-project", "pg-123")
	if err != nil {
		t.Fatalf("GetPage failed: %v", err)
	}

	if page.Title != "Getting Started" {
		t.Errorf("Expected page title 'Getting Started', got '%s'", page.Title)
	}
	if page.Content != "Full content here" {
		t.Errorf("Expected full content, got '%s'", page.Content)
	}
}

func TestCreatePage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req taskforge.CreatePageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Title != "My Page" {
			t.Errorf("Expected title 'My Page', got '%s'", req.Title)
		}
		if req.CategoryID != "pc-1" {
			t.Errorf("Expected category_id 'pc-1', got '%s'", req.CategoryID)
		}

		page := taskforge.Page{
			ID:         "pg-new",
			Title:      req.Title,
			Slug:       req.Slug,
			CategoryID: req.CategoryID,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(page)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	req := taskforge.CreatePageRequest{
		Title:      "My Page",
		Slug:       "my-page",
		CategoryID: "pc-1",
	}
	page, err := client.CreatePage("test-project", req)
	if err != nil {
		t.Fatalf("CreatePage failed: %v", err)
	}

	if page.Title != "My Page" {
		t.Errorf("Expected page title 'My Page', got '%s'", page.Title)
	}
}

func TestDeletePage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/pages/pg-123/" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	err := client.DeletePage("test-project", "pg-123")
	if err != nil {
		t.Fatalf("DeletePage failed: %v", err)
	}
}
