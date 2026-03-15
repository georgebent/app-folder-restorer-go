package runner

import "testing"

func TestBuildBackupNameUsesMaxNumericPrefix(t *testing.T) {
	folders := []string{"1.first", "tmp", "3.quick_save", "notes", "10.release"}

	got := buildBackupName(folders, "manual")

	if got != "11.manual" {
		t.Fatalf("expected next backup name to use max numeric prefix, got %q", got)
	}
}
