package api

import (
	"fmt"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
)

// ListStates retrieves all states for a project
func (c *Client) ListStates(projectID string) ([]taskforge.State, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/states/", c.Workspace, projectID)

	var response struct {
		Results []taskforge.State `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// GetState retrieves a specific state by ID
func (c *Client) GetState(projectID, stateID string) (*taskforge.State, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/states/%s/", c.Workspace, projectID, stateID)

	var state taskforge.State
	if err := c.Get(path, nil, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// CreateState creates a new state in a project
func (c *Client) CreateState(projectID string, req taskforge.CreateStateRequest) (*taskforge.State, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/states/", c.Workspace, projectID)

	var state taskforge.State
	if err := c.Post(path, req, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// UpdateState updates an existing state
func (c *Client) UpdateState(projectID, stateID string, req taskforge.UpdateStateRequest) (*taskforge.State, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/states/%s/", c.Workspace, projectID, stateID)

	var state taskforge.State
	if err := c.Patch(path, req, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// DeleteState removes a state from a project
func (c *Client) DeleteState(projectID, stateID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/states/%s/", c.Workspace, projectID, stateID)
	return c.Delete(path)
}
