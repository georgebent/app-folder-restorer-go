package file_manager

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const BackupExtension = ".zip"

func CreateArchive(sourceDir, archivePath string) error {
	sourceInfo, err := os.Stat(sourceDir)
	if err != nil {
		return err
	}
	if !sourceInfo.IsDir() {
		return fmt.Errorf("%s не є папкою", sourceDir)
	}

	_, err = os.Stat(archivePath)
	if err == nil {
		return fmt.Errorf("%s already exists", archivePath)
	}
	if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(archivePath), 0o755); err != nil {
		return err
	}

	tempFile, err := os.CreateTemp(filepath.Dir(archivePath), "*.zip.tmp")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()
	removeTemp := true
	defer func() {
		if removeTemp {
			_ = os.Remove(tempPath)
		}
	}()

	archiveWriter := zip.NewWriter(tempFile)
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		relativePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		if relativePath == "." {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(relativePath)
		if info.IsDir() {
			header.Name += "/"
			header.Method = zip.Store
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archiveWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		_, copyErr := io.Copy(writer, file)
		closeErr := file.Close()
		if copyErr != nil {
			return copyErr
		}

		return closeErr
	})
	if err != nil {
		_ = archiveWriter.Close()
		_ = tempFile.Close()
		return err
	}

	if err := archiveWriter.Close(); err != nil {
		_ = tempFile.Close()
		return err
	}

	if err := tempFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tempPath, archivePath); err != nil {
		return err
	}

	removeTemp = false
	return nil
}

func ExtractArchive(archivePath, destination string) error {
	archiveReader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = archiveReader.Close()
	}()

	if err := clearFolder(destination); err != nil {
		return err
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return err
	}

	destinationRoot := filepath.Clean(destination)
	for _, file := range archiveReader.File {
		targetPath := filepath.Join(destinationRoot, file.Name)
		cleanTargetPath := filepath.Clean(targetPath)
		if cleanTargetPath != destinationRoot && !strings.HasPrefix(cleanTargetPath, destinationRoot+string(os.PathSeparator)) {
			return fmt.Errorf("invalid archive entry: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(cleanTargetPath, file.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(cleanTargetPath), 0o755); err != nil {
			return err
		}

		archiveFile, err := file.Open()
		if err != nil {
			return err
		}

		destinationFile, err := os.OpenFile(cleanTargetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, file.Mode())
		if err != nil {
			_ = archiveFile.Close()
			return err
		}

		_, copyErr := io.Copy(destinationFile, archiveFile)
		closeArchiveErr := archiveFile.Close()
		closeDestinationErr := destinationFile.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeArchiveErr != nil {
			return closeArchiveErr
		}
		if closeDestinationErr != nil {
			return closeDestinationErr
		}
	}

	return nil
}

func ListBackups(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	backups := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), BackupExtension) {
			continue
		}

		backups = append(backups, strings.TrimSuffix(file.Name(), BackupExtension))
	}

	sort.Slice(backups, func(i, j int) bool {
		leftIndex, leftOk := backupArchiveIndex(backups[i])
		rightIndex, rightOk := backupArchiveIndex(backups[j])
		if leftOk && rightOk && leftIndex != rightIndex {
			return leftIndex < rightIndex
		}
		if leftOk != rightOk {
			return leftOk
		}

		return backups[i] < backups[j]
	})
	return backups, nil
}

func backupArchiveIndex(name string) (int, bool) {
	prefix, _, ok := strings.Cut(name, ".")
	if !ok {
		return 0, false
	}

	index, err := strconv.Atoi(prefix)
	if err != nil || index < 1 {
		return 0, false
	}

	return index, true
}
