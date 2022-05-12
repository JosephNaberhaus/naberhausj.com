package resource

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func MoveResources(sourceDirectory, outputDirectory string) error {
	resourcePaths, err := findResources(sourceDirectory)
	if err != nil {
		return err
	}

	for _, resourcePath := range resourcePaths {
		if !strings.HasPrefix(resourcePath, sourceDirectory) {
			return fmt.Errorf("file %s not located in source directory", resourcePaths)
		}

		newPath := strings.Replace(resourcePath, sourceDirectory, outputDirectory, 1)

		err := os.MkdirAll(filepath.Dir(newPath), 0777)
		if err != nil {
			return fmt.Errorf("error creating directories for resource: %w", err)
		}

		err = os.Link(resourcePath, newPath)
		if err != nil {
			return fmt.Errorf("error linking resource: %w", err)
		}
	}

	return nil
}
