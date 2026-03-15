package file_manager

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCreateArchiveAndExtractArchiveRoundTrip(t *testing.T) {
	root := t.TempDir()
	sourceDir := filepath.Join(root, "source")
	archivePath := filepath.Join(root, "backups", "1.sample.zip")
	restoreDir := filepath.Join(root, "restore")

	if err := os.MkdirAll(filepath.Join(sourceDir, "nested", "empty"), 0o755); err != nil {
		t.Fatalf("mkdir source tree: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "root.txt"), []byte("root"), 0o644); err != nil {
		t.Fatalf("write root file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "nested", "file.txt"), []byte("nested"), 0o644); err != nil {
		t.Fatalf("write nested file: %v", err)
	}

	if err := CreateArchive(sourceDir, archivePath); err != nil {
		t.Fatalf("create archive: %v", err)
	}
	if err := ExtractArchive(archivePath, restoreDir); err != nil {
		t.Fatalf("extract archive: %v", err)
	}

	rootContent, err := os.ReadFile(filepath.Join(restoreDir, "root.txt"))
	if err != nil {
		t.Fatalf("read restored root file: %v", err)
	}
	if string(rootContent) != "root" {
		t.Fatalf("unexpected root file content: %q", string(rootContent))
	}

	nestedContent, err := os.ReadFile(filepath.Join(restoreDir, "nested", "file.txt"))
	if err != nil {
		t.Fatalf("read restored nested file: %v", err)
	}
	if string(nestedContent) != "nested" {
		t.Fatalf("unexpected nested file content: %q", string(nestedContent))
	}

	if info, err := os.Stat(filepath.Join(restoreDir, "nested", "empty")); err != nil || !info.IsDir() {
		t.Fatalf("expected empty directory to be restored, err=%v", err)
	}
}

func TestListBackupsReturnsZipArchivesOnly(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "2.quick_save.zip"), []byte("archive"), 0o644); err != nil {
		t.Fatalf("write archive: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "1.initial.zip"), []byte("archive"), 0o644); err != nil {
		t.Fatalf("write archive: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "notes.txt"), []byte("ignore"), 0o644); err != nil {
		t.Fatalf("write other file: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "folder-backup"), 0o755); err != nil {
		t.Fatalf("mkdir folder: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "10.release.zip"), []byte("archive"), 0o644); err != nil {
		t.Fatalf("write archive: %v", err)
	}

	backups, err := ListBackups(root)
	if err != nil {
		t.Fatalf("list backups: %v", err)
	}

	expected := []string{"1.initial", "2.quick_save", "10.release"}
	if !reflect.DeepEqual(backups, expected) {
		t.Fatalf("expected %v, got %v", expected, backups)
	}
}
