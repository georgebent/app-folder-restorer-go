package runner

import (
	"fmt"
	"github.com/georgebent/go-restorer/pkg/core"
	"github.com/georgebent/go-restorer/pkg/file_manager"
	"github.com/georgebent/go-restorer/pkg/io_manager"
	"strconv"
	"strings"
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
	folders, err := file_manager.ListFolders(backups)
	if err != nil {
		return err
	}

	name = buildBackupName(folders, name)

	source := core.GetEnv("ORIGIN_DIR")
	backupPath := fmt.Sprintf("%s/%s", backups, name)

	err = file_manager.Copy(source, backupPath)
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
