package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
)

type IssueListOptions struct {
	State    string // State UUID or name (resolved by caller for consistent behavior)
	Assignee string // Assignee UUID or "me" (resolved by caller for consistent behavior)
	Limit    int
	Offset   int
}

func (c *Client) ListIssues(projectID string, opts IssueListOptions) ([]taskforge.Issue, *Pagination, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/", c.Workspace, projectID)

	query := url.Values{}
	query.Set("expand", "assignees,state")
	if opts.State != "" {
		query.Set("state", opts.State)
	}
	if opts.Assignee != "" {
		query.Set("assignee", opts.Assignee)
	}
	if opts.Limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Offset > 0 {
		query.Set("offset", fmt.Sprintf("%d", opts.Offset))
	}

	var response Response
	if err := c.Get(path, query, &response); err != nil {
		return nil, nil, err
	}

	var issues []taskforge.Issue
	if err := json.Unmarshal(response.Results, &issues); err != nil {
		return nil, nil, err
	}

	return issues, &response.Pagination, nil
}

func (c *Client) GetIssue(projectID, issueID string) (*taskforge.Issue, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/", c.Workspace, projectID, issueID)

	query := url.Values{}
	query.Set("expand", "assignees,state,labels")

	var issue taskforge.Issue
	if err := c.Get(path, query, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

func (c *Client) GetIssueByIdentifier(identifier string) (*taskforge.Issue, error) {
	path := fmt.Sprintf("/workspaces/%s/work-items/%s/", c.Workspace, url.PathEscape(identifier))

	query := url.Values{}
	query.Set("expand", "assignees,state,labels")

	var issue taskforge.Issue
	if err := c.Get(path, query, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

func (c *Client) GetIssueBySequenceID(projectID string, sequenceID int) (*taskforge.Issue, error) {
	project, err := c.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	identifier := strings.TrimSpace(project.Identifier)
	if identifier == "" {
		return nil, fmt.Errorf("project %s has no identifier", projectID)
	}

	return c.GetIssueByIdentifier(fmt.Sprintf("%s-%d", identifier, sequenceID))
}

func (c *Client) CreateIssue(projectID string, req taskforge.CreateIssueRequest) (*taskforge.Issue, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/", c.Workspace, projectID)

	var issue taskforge.Issue
	if err := c.Post(path, req, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

func (c *Client) UpdateIssue(projectID, issueID string, req taskforge.UpdateIssueRequest) (*taskforge.Issue, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/", c.Workspace, projectID, issueID)

	var issue taskforge.Issue
	if err := c.Patch(path, req, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

func (c *Client) DeleteIssue(projectID, issueID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/", c.Workspace, projectID, issueID)
	return c.Delete(path)
}

// SearchIssues searches for issues across the workspace
// Endpoint: GET /api/v1/workspaces/{workspace_slug}/work-items/search/
// Returns: {"issues": [...]}
func (c *Client) SearchIssues(query string) ([]taskforge.Issue, error) {
	path := fmt.Sprintf("/workspaces/%s/work-items/search/", c.Workspace)

	params := url.Values{}
	params.Set("search", query)

	// The search endpoint returns a different structure
	var response struct {
		Issues []taskforge.Issue `json:"issues"`
	}

	if err := c.Get(path, params, &response); err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && (apiErr.StatusCode == 400 || apiErr.StatusCode == 404) {
			return c.searchIssuesFallback(query)
		}
		return nil, err
	}

	return response.Issues, nil
}

func (c *Client) searchIssuesFallback(query string) ([]taskforge.Issue, error) {
	projects, err := c.ListProjects()
	if err != nil {
		return nil, err
	}

	needle := strings.ToLower(strings.TrimSpace(query))
	if needle == "" {
		return []taskforge.Issue{}, nil
	}

	var matches []taskforge.Issue
	for _, project := range projects {
		offset := 0
		for {
			issues, page, err := c.ListIssues(project.ID, IssueListOptions{Limit: 100, Offset: offset})
			if err != nil {
				return nil, err
			}
			for _, issue := range issues {
				hay := strings.ToLower(strings.Join([]string{issue.ID, issue.Identifier, issue.Name, issue.Description}, " "))
				if strings.Contains(hay, needle) {
					matches = append(matches, issue)
				}
			}
			if page == nil || page.Next == nil || len(issues) == 0 {
				break
			}
			offset += len(issues)
		}
	}

	return matches, nil
}

// PatchIssue applies RFC 6902-style JSON Patch operations to an issue
func (c *Client) PatchIssue(projectID, issueID string, patches []taskforge.PatchOp) (*taskforge.Issue, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/", c.Workspace, projectID, issueID)

	req := taskforge.PatchRequest{Patches: patches}

	var issue taskforge.Issue
	if err := c.Patch(path, req, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

// ListIssueVersions retrieves the version history for an issue
func (c *Client) ListIssueVersions(projectID, issueID string) ([]taskforge.VersionHistory, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/versions/", c.Workspace, projectID, issueID)

	var response Response
	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	var versions []taskforge.VersionHistory
	if err := json.Unmarshal(response.Results, &versions); err != nil {
		return nil, err
	}

	return versions, nil
}

// GetIssueVersion retrieves a specific version of an issue
func (c *Client) GetIssueVersion(projectID, issueID, field string, versionNum int) (*taskforge.VersionHistory, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/versions/%d/?field=%s", c.Workspace, projectID, issueID, versionNum, field)

	var version taskforge.VersionHistory
	if err := c.Get(path, nil, &version); err != nil {
		return nil, err
	}

	return &version, nil
}

// GetIssueVersionDiff retrieves the diff for a specific version of an issue
func (c *Client) GetIssueVersionDiff(projectID, issueID, field string, versionNum int) (string, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/versions/%d/diff/?field=%s", c.Workspace, projectID, issueID, versionNum, field)

	body, err := c.GetRaw(path, nil)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// RestoreIssue restores a deleted issue to a specific version
func (c *Client) RestoreIssue(projectID, issueID string, versionNum int) (*taskforge.Issue, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/restore/", c.Workspace, projectID, issueID)

	req := taskforge.IssueRestoreRequest{Version: versionNum}

	var issue taskforge.Issue
	if err := c.Post(path, req, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

// ListCommentVersions retrieves the version history for a comment
func (c *Client) ListCommentVersions(projectID, issueID, commentID string) ([]taskforge.VersionHistory, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/comments/%s/versions/", c.Workspace, projectID, issueID, commentID)

	var response Response
	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	var versions []taskforge.VersionHistory
	if err := json.Unmarshal(response.Results, &versions); err != nil {
		return nil, err
	}

	return versions, nil
}

// GetCommentVersion retrieves a specific version of a comment
func (c *Client) GetCommentVersion(projectID, issueID, commentID string, versionNum int) (*taskforge.VersionHistory, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/comments/%s/versions/%d/", c.Workspace, projectID, issueID, commentID, versionNum)

	var version taskforge.VersionHistory
	if err := c.Get(path, nil, &version); err != nil {
		return nil, err
	}

	return &version, nil
}

// GetCommentVersionDiff retrieves the diff for a specific version of a comment
func (c *Client) GetCommentVersionDiff(projectID, issueID, commentID string, versionNum int) (string, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/comments/%s/versions/%d/diff/", c.Workspace, projectID, issueID, commentID, versionNum)

	body, err := c.GetRaw(path, nil)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// RestoreComment restores a deleted comment to a specific version
func (c *Client) RestoreComment(projectID, issueID, commentID string, versionNum int) (*taskforge.Comment, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/comments/%s/restore/", c.Workspace, projectID, issueID, commentID)

	req := taskforge.CommentRestoreRequest{Version: versionNum}

	var comment taskforge.Comment
	if err := c.Post(path, req, &comment); err != nil {
		return nil, err
	}

	return &comment, nil
}
