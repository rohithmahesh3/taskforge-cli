package api

import (
	"fmt"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
)

func (c *Client) ListIssueDependencies(projectID, issueID string) (*taskforge.GroupedDependencies, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/dependencies/", c.Workspace, projectID, issueID)

	var deps taskforge.GroupedDependencies
	if err := c.Get(path, nil, &deps); err != nil {
		return nil, err
	}

	return &deps, nil
}

func (c *Client) AddIssueDependency(projectID, issueID, targetIssueID string, isBlocks bool) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/dependencies/bulk/", c.Workspace, projectID, issueID)
	
	payload := map[string][]string{}
	if isBlocks {
		payload["blocks"] = []string{targetIssueID}
	} else {
		payload["depends_on"] = []string{targetIssueID}
	}

	return c.Post(path, payload, nil)
}

func (c *Client) RemoveIssueDependency(projectID, issueID, dependencyID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/dependencies/%s/", c.Workspace, projectID, issueID, dependencyID)
	return c.Delete(path)
}
