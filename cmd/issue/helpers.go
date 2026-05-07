package issue

import (
	"strconv"
	"strings"

	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
)

func resolveIssue(client *api.Client, projectID, ref string) (*taskforge.Issue, error) {
	if seqID, err := strconv.Atoi(strings.TrimSpace(ref)); err == nil {
		return client.GetIssueBySequenceID(projectID, seqID)
	}

	if looksLikeUUID(ref) {
		return client.GetIssue(projectID, ref)
	}

	issue, err := client.GetIssue(projectID, ref)
	if err == nil {
		return issue, nil
	}

	return client.GetIssueByIdentifier(ref)
}

func resolveIssueContext(client *api.Client, projectID, ref string) (string, string, error) {
	issue, err := resolveIssue(client, projectID, ref)
	if err != nil {
		return "", "", err
	}

	resolvedProjectID := strings.TrimSpace(issue.ProjectID)
	if resolvedProjectID == "" {
		resolvedProjectID = projectID
	}

	return resolvedProjectID, issue.ID, nil
}

func looksLikeUUID(s string) bool {
	if len(s) != 36 {
		return false
	}

	for i, r := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if r != '-' {
				return false
			}
			continue
		}

		if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
			return false
		}
	}

	return true
}
