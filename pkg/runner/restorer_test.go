package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRestoreFromBackupRollsBackWhenRestoreFails(t *testing.T) {
	root := t.TempDir()
	sourceDir := filepath.Join(root, "origin")
	backupRoot := filepath.Join(root, "backups")
	backupPath := filepath.Join(backupRoot, "missing")
	tmpPath := filepath.Join(backupRoot, "tmp")

	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatalf("mkdir source: %v", err)
	}
	if err := os.MkdirAll(backupRoot, 0o755); err != nil {
		t.Fatalf("mkdir backups: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "file.txt"), []byte("origin-state"), 0o644); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	err := restoreFromBackup(sourceDir, backupPath, tmpPath)
	if err == nil {
		t.Fatal("expected restore to fail for missing backup path")
	}

	content, readErr := os.ReadFile(filepath.Join(sourceDir, "file.txt"))
	if readErr != nil {
		t.Fatalf("read restored source file: %v", readErr)
	}
	if string(content) != "origin-state" {
		t.Fatalf("expected rollback to restore original source contents, got %q", string(content))
	}
}
