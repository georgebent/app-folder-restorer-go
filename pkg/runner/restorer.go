package runner

import (
	"fmt"
	"github.com/georgebent/go-restorer/pkg/core"
	"github.com/georgebent/go-restorer/pkg/file_manager"
	"github.com/georgebent/go-restorer/pkg/io_manager"
	"strconv"
)

func Restore() error {
	backups := core.GetEnv("BACKUP_DIR")
	source := core.GetEnv("ORIGIN_DIR")

	folders, err := file_manager.ListFolders(backups)
	if err != nil {
		return err
	}

	options := map[string]string{}
	for i, folder := range folders {
		key := strconv.Itoa(i + 1)
		options[key] = folder
	}

	chosen := io_manager.Ask("Choose restore file", options)

	backupsTmpPath := fmt.Sprintf("%s/tmp", backups)
	backupPath := fmt.Sprintf("%s/%s", backups, options[chosen])

	err = restoreFromBackup(source, backupPath, backupsTmpPath)
	if err != nil {
		return err
	}

	io_manager.Write(fmt.Sprintf("Folder %s restored from %s", source, backupPath))

	return nil
}

func restoreFromBackup(source, backupPath, backupsTmpPath string) error {
	err := file_manager.ForceCopy(source, backupsTmpPath)
	if err != nil {
		return err
	}

	err = file_manager.ForceCopy(backupPath, source)
	if err == nil {
		return nil
	}

	rollbackErr := file_manager.ForceCopy(backupsTmpPath, source)
	if rollbackErr != nil {
		return fmt.Errorf("restore failed: %w; rollback failed: %v", err, rollbackErr)
	}

	return fmt.Errorf("restore failed and rolled back: %w", err)
}
