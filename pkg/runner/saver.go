package runner

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/georgebent/go-restorer/pkg/core"
	"github.com/georgebent/go-restorer/pkg/file_manager"
	"github.com/georgebent/go-restorer/pkg/io_manager"
)

func QuickSave() error {
	return saveByName("quick_save")
}

func Save() error {
	name := io_manager.Read("Enter backup name: ")

	return saveByName(name)
}

func saveByName(name string) error {
	backups := core.GetEnv("BACKUP_DIR")
	backupsList, err := file_manager.ListBackups(backups)
	if err != nil {
		return err
	}

	name = buildBackupName(backupsList, name)

	source := core.GetEnv("ORIGIN_DIR")
	backupPath := backupArchivePath(backups, name)

	err = file_manager.CreateArchive(source, backupPath)
	if err != nil {
		return err
	}

	io_manager.Write(fmt.Sprintf("Folder %s saved to %s", source, backupPath))

	return nil
}

func buildBackupName(folders []string, name string) string {
	return fmt.Sprintf("%d.%s", nextBackupIndex(folders), name)
}

func nextBackupIndex(folders []string) int {
	maxIndex := 0

	for _, folder := range folders {
		index, ok := backupIndex(folder)
		if ok && index > maxIndex {
			maxIndex = index
		}
	}

	return maxIndex + 1
}

func backupIndex(folder string) (int, bool) {
	prefix, _, ok := strings.Cut(folder, ".")
	if !ok {
		return 0, false
	}

	index, err := strconv.Atoi(prefix)
	if err != nil || index < 1 {
		return 0, false
	}

	return index, true
}

func backupArchivePath(backupsDir, backupName string) string {
	return filepath.Join(backupsDir, backupName+file_manager.BackupExtension)
}
