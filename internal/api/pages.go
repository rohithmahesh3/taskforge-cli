package api

import (
	"fmt"
	"net/url"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
)

// ── Page Categories ──

// ListPageCategories retrieves all page categories for a project.
func (c *Client) ListPageCategories(projectID string) ([]taskforge.PageCategory, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/page-categories/", c.Workspace, projectID)

	var response struct {
		Results []taskforge.PageCategory `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// CreatePageCategory creates a new page category.
func (c *Client) CreatePageCategory(projectID string, req taskforge.CreatePageCategoryRequest) (*taskforge.PageCategory, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/page-categories/", c.Workspace, projectID)

	var category taskforge.PageCategory
	if err := c.Post(path, req, &category); err != nil {
		return nil, err
	}

	return &category, nil
}

// UpdatePageCategory updates an existing page category.
func (c *Client) UpdatePageCategory(projectID, categoryID string, req taskforge.UpdatePageCategoryRequest) (*taskforge.PageCategory, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/page-categories/%s/", c.Workspace, projectID, categoryID)

	var category taskforge.PageCategory
	if err := c.Patch(path, req, &category); err != nil {
		return nil, err
	}

	return &category, nil
}

// DeletePageCategory removes a page category.
func (c *Client) DeletePageCategory(projectID, categoryID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/page-categories/%s/", c.Workspace, projectID, categoryID)
	return c.Delete(path)
}

// ── Pages ──

// ListPages retrieves all pages for a project with optional search and category filter.
func (c *Client) ListPages(projectID string, searchQuery, categoryID string) ([]taskforge.Page, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/pages/", c.Workspace, projectID)

	query := url.Values{}
	if searchQuery != "" {
		query.Set("q", searchQuery)
	}
	if categoryID != "" {
		query.Set("category_id", categoryID)
	}

	var response struct {
		Results []taskforge.Page `json:"results"`
	}

	var q url.Values
	if len(query) > 0 {
		q = query
	}

	if err := c.Get(path, q, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// GetPage retrieves a specific page by ID.
func (c *Client) GetPage(projectID, pageID string) (*taskforge.Page, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/pages/%s/", c.Workspace, projectID, pageID)

	var page taskforge.Page
	if err := c.Get(path, nil, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// CreatePage creates a new page.
func (c *Client) CreatePage(projectID string, req taskforge.CreatePageRequest) (*taskforge.Page, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/pages/", c.Workspace, projectID)

	var page taskforge.Page
	if err := c.Post(path, req, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// UpdatePage updates an existing page.
func (c *Client) UpdatePage(projectID, pageID string, req taskforge.UpdatePageRequest) (*taskforge.Page, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/pages/%s/", c.Workspace, projectID, pageID)

	var page taskforge.Page
	if err := c.Patch(path, req, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// DeletePage removes a page.
func (c *Client) DeletePage(projectID, pageID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/pages/%s/", c.Workspace, projectID, pageID)
	return c.Delete(path)
}
