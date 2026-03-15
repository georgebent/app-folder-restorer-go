package file_manager

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Copy(source, destination string) error {
	_, err := os.Stat(destination)
	if err == nil {
		return fmt.Errorf("%s already exists", destination)
	}
	if !os.IsNotExist(err) {
		return err
	}

	return copyFolder(source, destination)
}

func ForceCopy(source, destination string) error {
	err := clearFolder(destination)
	if err != nil {
		return err
	}

	return copyFolder(source, destination)
}

func ListFolders(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var directories []string
	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, file.Name())
		}
	}

	return directories, nil
}

func copyFolder(src, dst string) error {
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceInfo.IsDir() {
		return fmt.Errorf("%s не є папкою", src)
	}

	err = os.MkdirAll(dst, sourceInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destinationPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyFolder(sourcePath, destinationPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(sourcePath, destinationPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) (err error) {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := sourceFile.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := destinationFile.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func clearFolder(path string) (err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dir.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	entries, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		err = os.RemoveAll(filepath.Join(path, entry.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}
