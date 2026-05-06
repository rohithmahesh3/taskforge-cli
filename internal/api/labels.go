package api

import (
	"fmt"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
)

// ListLabels retrieves all labels for a project
func (c *Client) ListLabels(projectID string) ([]taskforge.Label, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/labels/", c.Workspace, projectID)

	var response struct {
		Results []taskforge.Label `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// GetLabel retrieves a specific label by ID
func (c *Client) GetLabel(projectID, labelID string) (*taskforge.Label, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/labels/%s/", c.Workspace, projectID, labelID)

	var label taskforge.Label
	if err := c.Get(path, nil, &label); err != nil {
		return nil, err
	}

	return &label, nil
}

// CreateLabel creates a new label in a project
func (c *Client) CreateLabel(projectID string, req taskforge.CreateLabelRequest) (*taskforge.Label, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/labels/", c.Workspace, projectID)

	var label taskforge.Label
	if err := c.Post(path, req, &label); err != nil {
		return nil, err
	}

	return &label, nil
}

// UpdateLabel updates an existing label
func (c *Client) UpdateLabel(projectID, labelID string, req taskforge.UpdateLabelRequest) (*taskforge.Label, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/labels/%s/", c.Workspace, projectID, labelID)

	var label taskforge.Label
	if err := c.Patch(path, req, &label); err != nil {
		return nil, err
	}

	return &label, nil
}

// DeleteLabel removes a label from a project
func (c *Client) DeleteLabel(projectID, labelID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/labels/%s/", c.Workspace, projectID, labelID)
	return c.Delete(path)
}
