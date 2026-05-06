package api

import (
	"strings"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
)

// FilterWorkspaceMembers performs client-side filtering for workspace members.
// It does not imply any server-side search support in the TaskForge API.
func FilterWorkspaceMembers(members []taskforge.User, query string, exact bool, limit int) []taskforge.User {
	if len(members) == 0 {
		return nil
	}

	if strings.TrimSpace(query) == "" {
		return limitWorkspaceMembers(members, limit)
	}

	normalizedQuery := normalizeMemberSearchValue(query)
	filtered := make([]taskforge.User, 0, len(members))

	for _, member := range members {
		if memberMatchesQuery(member, normalizedQuery, exact) {
			filtered = append(filtered, member)
		}
	}

	return limitWorkspaceMembers(filtered, limit)
}

func memberMatchesQuery(member taskforge.User, normalizedQuery string, exact bool) bool {
	for _, candidate := range memberSearchFields(member) {
		if candidate == "" {
			continue
		}

		if exact {
			if candidate == normalizedQuery {
				return true
			}
			continue
		}

		if strings.Contains(candidate, normalizedQuery) {
			return true
		}
	}

	return false
}

func limitWorkspaceMembers(members []taskforge.User, limit int) []taskforge.User {
	if limit <= 0 || len(members) <= limit {
		return members
	}

	return members[:limit]
}

func memberFullName(member taskforge.User) string {
	parts := make([]string, 0, 2)
	if first := strings.TrimSpace(member.FirstName); first != "" {
		parts = append(parts, first)
	}
	if last := strings.TrimSpace(member.LastName); last != "" {
		parts = append(parts, last)
	}

	return strings.Join(parts, " ")
}

func memberSearchFields(member taskforge.User) []string {
	fullName := memberFullName(member)

	return []string{
		normalizeMemberSearchValue(member.DisplayName),
		normalizeMemberSearchValue(member.Email),
		normalizeMemberSearchValue(member.FirstName),
		normalizeMemberSearchValue(member.LastName),
		normalizeMemberSearchValue(fullName),
		normalizeMemberSearchValue(member.ID),
	}
}

func normalizeMemberSearchValue(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
