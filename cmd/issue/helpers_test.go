package issue

import "testing"

func TestLooksLikeUUID(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"550E8400-E29B-41D4-A716-446655440000", true},
		{"not-a-uuid", false},
		{"short", false},
		{"550e8400-e29b-41d4-a716", false},
		{"550e8400e29b41d4a716446655440000", false},
		{"", false},
	}

	for _, tt := range tests {
		got := looksLikeUUID(tt.input)
		if got != tt.want {
			t.Errorf("looksLikeUUID(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
