package taskforge

import (
	"encoding/json"
	"testing"
)

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

func TestFlexibleIssueUnmarshalFromString(t *testing.T) {
	var fi FlexibleIssue
	if err := json.Unmarshal([]byte(`"issue-uuid-123"`), &fi); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fi.ID != "issue-uuid-123" {
		t.Fatalf("expected ID issue-uuid-123, got %q", fi.ID)
	}
	if !fi.IsUUID {
		t.Fatal("expected IsUUID=true for string input")
	}
}

func TestFlexibleIssueUnmarshalFromObject(t *testing.T) {
	var fi FlexibleIssue
	input := `{"id":"i-1","name":"Bug fix","identifier":"PROJ-42","priority":"high","sequence_id":42}`
	if err := json.Unmarshal([]byte(input), &fi); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fi.ID != "i-1" {
		t.Fatalf("expected ID i-1, got %q", fi.ID)
	}
	if fi.Name != "Bug fix" {
		t.Fatalf("expected name 'Bug fix', got %q", fi.Name)
	}
	if fi.Identifier != "PROJ-42" {
		t.Fatalf("expected identifier PROJ-42, got %q", fi.Identifier)
	}
	if fi.Priority != "high" {
		t.Fatalf("expected priority high, got %q", fi.Priority)
	}
	if fi.SequenceID != 42 {
		t.Fatalf("expected sequence_id 42, got %d", fi.SequenceID)
	}
	if fi.IsUUID {
		t.Fatal("expected IsUUID=false for object input")
	}
}

func TestFlexibleIssueMarshalStringUUID(t *testing.T) {
	fi := FlexibleIssue{ID: "uuid-1", IsUUID: true}
	data, err := json.Marshal(fi)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `"uuid-1"` {
		t.Fatalf("expected marshaled UUID string, got %s", data)
	}
}

func TestFlexibleIssueMarshalObject(t *testing.T) {
	fi := FlexibleIssue{ID: "i-1", Name: "Test Issue", IsUUID: false}
	data, err := json.Marshal(fi)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unexpected error unmarshaling result: %v", err)
	}
	if result["id"] != "i-1" {
		t.Fatalf("expected id i-1, got %v", result["id"])
	}
	if result["name"] != "Test Issue" {
		t.Fatalf("expected name 'Test Issue', got %v", result["name"])
	}
}

func TestFlexibleIssueToIssue(t *testing.T) {
	fi := FlexibleIssue{ID: "i-1", Name: "Test", Identifier: "P-1", SequenceID: 5, Priority: "medium"}
	issue := fi.ToIssue()
	if issue.ID != "i-1" || issue.Name != "Test" || issue.Identifier != "P-1" || issue.SequenceID != 5 || issue.Priority != "medium" {
		t.Fatalf("ToIssue conversion mismatch: %+v", issue)
	}
}

func TestIntakeIssueUnmarshalWithExpandedIssue(t *testing.T) {
	input := `{
		"id": "intake-1",
		"status": 1,
		"source": "email",
		"issue": {"id": "issue-1", "name": "Feature request", "priority": "high"}
	}`
	var ii IntakeIssue
	if err := json.Unmarshal([]byte(input), &ii); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ii.Issue.ID != "issue-1" {
		t.Fatalf("expected issue ID issue-1, got %q", ii.Issue.ID)
	}
	if ii.Issue.Name != "Feature request" {
		t.Fatalf("expected issue name 'Feature request', got %q", ii.Issue.Name)
	}
	if ii.Issue.IsUUID {
		t.Fatal("expected IsUUID=false for expanded issue object")
	}
}
