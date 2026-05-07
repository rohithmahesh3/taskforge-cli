package taskforge

import "testing"

func TestStateOutputFromIssueWithExpandedState(t *testing.T) {
	issue := Issue{
		State: FlexibleState{
			ID:   "state-1",
			Name: "Done",
		},
	}

	got := StateOutputFromIssue(issue)
	if got.ID != "state-1" {
		t.Fatalf("expected state id state-1, got %q", got.ID)
	}
	if got.Name != "Done" {
		t.Fatalf("expected state name Done, got %q", got.Name)
	}
}

func TestStateOutputFromIssueWithEmptyState(t *testing.T) {
	issue := Issue{
		State: FlexibleState{}, // zero value = no state
	}

	got := StateOutputFromIssue(issue)
	if got.ID != "" {
		t.Fatalf("expected empty state id, got %q", got.ID)
	}
	if got.Name != "" {
		t.Fatalf("expected empty state name, got %q", got.Name)
	}
}
