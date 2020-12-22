package resources

import (
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"strings"
)

func findResources(root string) (resourcePaths []string, err error) {
	resourcePaths = make([]string, 0)

	filePaths, err := file.FindFiles(root)
	if err != nil {
		return nil, err
	}

	for _, path := range filePaths {
		if !strings.HasSuffix(path, ".html") && !strings.HasSuffix(path, ".css") {
			resourcePaths = append(resourcePaths, path)
		}
	}

	return resourcePaths, nil
}
