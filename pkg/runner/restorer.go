package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/georgebent/go-restorer/pkg/core"
	"github.com/georgebent/go-restorer/pkg/file_manager"
	"github.com/georgebent/go-restorer/pkg/io_manager"
)

var copyDirectory = file_manager.Copy
var forceCopyDirectory = file_manager.ForceCopy
var extractBackupArchive = file_manager.ExtractArchive

func Restore() error {
	backups := core.GetEnv("BACKUP_DIR")
	source := core.GetEnv("ORIGIN_DIR")

	backupsList, err := file_manager.ListBackups(backups)
	if err != nil {
		return err
	}

	options := map[string]string{}
	for i, backup := range backupsList {
		key := strconv.Itoa(i + 1)
		options[key] = backup
	}

	chosen := io_manager.Ask("Choose restore file", options)
	backupPath := backupArchivePath(backups, options[chosen])

	err = restoreFromBackup(source, backupPath)
	if err != nil {
		return err
	}

	io_manager.Write(fmt.Sprintf("Folder %s restored from %s", source, backupPath))

	return nil
}

func restoreFromBackup(source, backupPath string) error {
	tempRoot, err := os.MkdirTemp("", "go-restorer-restore-*")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tempRoot)
	}()

	snapshotPath := filepath.Join(tempRoot, "snapshot")
	extractedPath := filepath.Join(tempRoot, "extracted")

	err = copyDirectory(source, snapshotPath)
	if err != nil {
		return err
	}

	err = extractBackupArchive(backupPath, extractedPath)
	if err != nil {
		return err
	}

	err = forceCopyDirectory(extractedPath, source)
	if err == nil {
		return nil
	}

	rollbackErr := forceCopyDirectory(snapshotPath, source)
	if rollbackErr != nil {
		return fmt.Errorf("restore failed: %w; rollback failed: %v", err, rollbackErr)
	}

	return fmt.Errorf("restore failed and rolled back: %w", err)
}
