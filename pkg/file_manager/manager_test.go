package file_manager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFailsWhenDestinationExists(t *testing.T) {
	sourceDir := t.TempDir()
	destinationRoot := t.TempDir()
	destinationDir := filepath.Join(destinationRoot, "backup")

	if err := os.MkdirAll(filepath.Join(sourceDir, "nested"), 0o755); err != nil {
		t.Fatalf("mkdir source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "nested", "new.txt"), []byte("new"), 0o644); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	if err := os.MkdirAll(destinationDir, 0o755); err != nil {
		t.Fatalf("mkdir destination: %v", err)
	}
	if err := os.WriteFile(filepath.Join(destinationDir, "stale.txt"), []byte("stale"), 0o644); err != nil {
		t.Fatalf("write destination file: %v", err)
	}

	err := Copy(sourceDir, destinationDir)
	if err == nil {
		t.Fatal("expected copy to fail when destination already exists")
	}

	if _, err := os.Stat(filepath.Join(destinationDir, "stale.txt")); err != nil {
		t.Fatalf("expected existing destination contents to remain: %v", err)
	}
	if _, err := os.Stat(filepath.Join(destinationDir, "nested", "new.txt")); !os.IsNotExist(err) {
		t.Fatalf("expected copy to avoid merging new files into destination, got err=%v", err)
	}
}
