package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindFilesWithSuffix(root, suffix string) (paths []string, err error) {
	paths = make([]string, 0)

	allPaths, err := FindFiles(root)
	if err != nil {
		return nil, err
	}

	for _, path := range allPaths {
		if strings.HasSuffix(path, suffix) {
			paths = append(paths, path)
		}
	}

	return paths, nil
}

func FindFiles(root string) (paths []string, err error) {
	paths = make([]string, 0)

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking project files")
		}

		if info.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("error finding relative path for file: %w", err)
		}

		paths = append(paths, relativePath)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return paths, nil
}
