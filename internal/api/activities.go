package api

import (
	"fmt"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
)

// ListActivities retrieves all activities for an issue
func (c *Client) ListActivities(projectID, issueID string) ([]taskforge.Activity, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/activities/", c.Workspace, projectID, issueID)

	var response struct {
		Results []taskforge.Activity `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// GetActivity retrieves a specific activity by ID
func (c *Client) GetActivity(projectID, issueID, activityID string) (*taskforge.Activity, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/activities/%s/", c.Workspace, projectID, issueID, activityID)

	var activity taskforge.Activity
	if err := c.Get(path, nil, &activity); err != nil {
		return nil, err
	}

	return &activity, nil
}
