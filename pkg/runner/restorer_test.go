package runner

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/georgebent/go-restorer/pkg/file_manager"
)

func TestRestoreFromBackupRollsBackWhenRestoreFails(t *testing.T) {
	root := t.TempDir()
	sourceDir := filepath.Join(root, "origin")
	backupSourceDir := filepath.Join(root, "backup-source")
	backupPath := filepath.Join(root, "backups", "1.sample.zip")

	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatalf("mkdir source: %v", err)
	}
	if err := os.MkdirAll(backupSourceDir, 0o755); err != nil {
		t.Fatalf("mkdir backup source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "file.txt"), []byte("origin-state"), 0o644); err != nil {
		t.Fatalf("write source file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(backupSourceDir, "file.txt"), []byte("backup-state"), 0o644); err != nil {
		t.Fatalf("write backup file: %v", err)
	}
	if err := file_manager.CreateArchive(backupSourceDir, backupPath); err != nil {
		t.Fatalf("create archive: %v", err)
	}

	previousForceCopyDirectory := forceCopyDirectory
	forceCopyDirectory = func(src, dst string) error {
		if filepath.Base(src) == "extracted" {
			return errors.New("restore failed")
		}

		return file_manager.ForceCopy(src, dst)
	}
	defer func() {
		forceCopyDirectory = previousForceCopyDirectory
	}()

	err := restoreFromBackup(sourceDir, backupPath)
	if err == nil {
		t.Fatal("expected restore to fail")
	}

	content, readErr := os.ReadFile(filepath.Join(sourceDir, "file.txt"))
	if readErr != nil {
		t.Fatalf("read restored source file: %v", readErr)
	}
	if string(content) != "origin-state" {
		t.Fatalf("expected rollback to restore original source contents, got %q", string(content))
	}
}
