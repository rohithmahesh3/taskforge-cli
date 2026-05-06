package api

import (
	"fmt"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
)

// GetUserInfo retrieves the current authenticated user info
// This can be used to verify authentication and get user details
func (c *Client) GetUserInfo() (*taskforge.User, error) {
	path := "/users/me/"

	var user taskforge.User
	if err := c.Get(path, nil, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetWorkspaceMembers retrieves all members of the workspace
func (c *Client) GetWorkspaceMembers() ([]taskforge.User, error) {
	path := fmt.Sprintf("/workspaces/%s/members/", c.Workspace)

	var members []taskforge.User
	if err := c.Get(path, nil, &members); err != nil {
		return nil, err
	}

	return members, nil
}
