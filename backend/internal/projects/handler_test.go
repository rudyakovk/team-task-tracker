package projects

import "testing"

func TestNormalizeProjectKey(t *testing.T) {
	t.Parallel()

	got := normalizeProjectKey(" core ")
	if got != "CORE" {
		t.Fatalf("normalizeProjectKey() = %q, want %q", got, "CORE")
	}
}

func TestValidateProjectInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		project string
		wantErr bool
	}{
		{name: "valid", key: "CORE", project: "Core Platform"},
		{name: "too short", key: "C", project: "Core Platform", wantErr: true},
		{name: "too long", key: "CORETRACKER", project: "Core Platform", wantErr: true},
		{name: "invalid chars", key: "CORE-1", project: "Core Platform", wantErr: true},
		{name: "missing name", key: "CORE", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateProjectInput(tt.key, tt.project)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
