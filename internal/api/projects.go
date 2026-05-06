package api

import (
	"fmt"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
)

func (c *Client) ListProjects() ([]taskforge.Project, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/", c.Workspace)

	var response struct {
		Results []taskforge.Project `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

func (c *Client) GetProject(projectID string) (*taskforge.Project, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/", c.Workspace, projectID)

	var project taskforge.Project
	if err := c.Get(path, nil, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) CreateProject(req taskforge.CreateProjectRequest) (*taskforge.Project, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/", c.Workspace)

	var project taskforge.Project
	if err := c.Post(path, req, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) DeleteProject(projectID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/", c.Workspace, projectID)
	return c.Delete(path)
}

func (c *Client) GetProjectMembers(projectID string) ([]taskforge.User, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/members/", c.Workspace, projectID)

	var members []taskforge.User
	if err := c.Get(path, nil, &members); err != nil {
		return nil, err
	}

	return members, nil
}
