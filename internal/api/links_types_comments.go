package api

import (
	"fmt"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
)

// ListLinks retrieves all links for an issue
func (c *Client) ListLinks(projectID, issueID string) ([]taskforge.Link, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/links/", c.Workspace, projectID, issueID)

	var response struct {
		Results []taskforge.Link `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// GetLink retrieves a specific link
func (c *Client) GetLink(projectID, issueID, linkID string) (*taskforge.Link, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/links/%s/", c.Workspace, projectID, issueID, linkID)

	var link taskforge.Link
	if err := c.Get(path, nil, &link); err != nil {
		return nil, err
	}

	return &link, nil
}

// CreateLink adds a new link to an issue
func (c *Client) CreateLink(projectID, issueID string, req taskforge.CreateLinkRequest) (*taskforge.Link, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/links/", c.Workspace, projectID, issueID)

	var link taskforge.Link
	if err := c.Post(path, req, &link); err != nil {
		return nil, err
	}

	return &link, nil
}

// UpdateLink updates an existing link
func (c *Client) UpdateLink(projectID, issueID, linkID string, req taskforge.UpdateLinkRequest) (*taskforge.Link, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/links/%s/", c.Workspace, projectID, issueID, linkID)

	var link taskforge.Link
	if err := c.Patch(path, req, &link); err != nil {
		return nil, err
	}

	return &link, nil
}

// DeleteLink removes a link from an issue
func (c *Client) DeleteLink(projectID, issueID, linkID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/links/%s/", c.Workspace, projectID, issueID, linkID)
	return c.Delete(path)
}

// ListIssueTypes retrieves all issue types for a project
func (c *Client) ListIssueTypes(projectID string) ([]taskforge.IssueType, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-item-types/", c.Workspace, projectID)

	var response struct {
		Results []taskforge.IssueType `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// GetIssueType retrieves a specific issue type
func (c *Client) GetIssueType(projectID, typeID string) (*taskforge.IssueType, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-item-types/%s/", c.Workspace, projectID, typeID)

	var issueType taskforge.IssueType
	if err := c.Get(path, nil, &issueType); err != nil {
		return nil, err
	}

	return &issueType, nil
}

// CreateIssueType creates a new issue type
func (c *Client) CreateIssueType(projectID string, req taskforge.CreateIssueTypeRequest) (*taskforge.IssueType, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-item-types/", c.Workspace, projectID)

	var issueType taskforge.IssueType
	if err := c.Post(path, req, &issueType); err != nil {
		return nil, err
	}

	return &issueType, nil
}

// UpdateIssueType updates an existing issue type
func (c *Client) UpdateIssueType(projectID, typeID string, req taskforge.UpdateIssueTypeRequest) (*taskforge.IssueType, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-item-types/%s/", c.Workspace, projectID, typeID)

	var issueType taskforge.IssueType
	if err := c.Patch(path, req, &issueType); err != nil {
		return nil, err
	}

	return &issueType, nil
}

// DeleteIssueType removes an issue type
func (c *Client) DeleteIssueType(projectID, typeID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-item-types/%s/", c.Workspace, projectID, typeID)
	return c.Delete(path)
}

// ListComments retrieves all comments for an issue
func (c *Client) ListComments(projectID, issueID string) ([]taskforge.Comment, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/comments/", c.Workspace, projectID, issueID)

	var response struct {
		Results []taskforge.Comment `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// GetComment retrieves a specific comment
func (c *Client) GetComment(projectID, issueID, commentID string) (*taskforge.Comment, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/comments/%s/", c.Workspace, projectID, issueID, commentID)

	var comment taskforge.Comment
	if err := c.Get(path, nil, &comment); err != nil {
		return nil, err
	}

	return &comment, nil
}

// CreateComment adds a new comment to an issue
func (c *Client) CreateComment(projectID, issueID string, req taskforge.CreateCommentRequest) (*taskforge.Comment, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/comments/", c.Workspace, projectID, issueID)

	var comment taskforge.Comment
	if err := c.Post(path, req, &comment); err != nil {
		return nil, err
	}

	return &comment, nil
}

// UpdateComment updates an existing comment
func (c *Client) UpdateComment(projectID, issueID, commentID string, req taskforge.UpdateCommentRequest) (*taskforge.Comment, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/comments/%s/", c.Workspace, projectID, issueID, commentID)

	var comment taskforge.Comment
	if err := c.Patch(path, req, &comment); err != nil {
		return nil, err
	}

	return &comment, nil
}

// DeleteComment removes a comment from an issue
func (c *Client) DeleteComment(projectID, issueID, commentID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/comments/%s/", c.Workspace, projectID, issueID, commentID)
	return c.Delete(path)
}
