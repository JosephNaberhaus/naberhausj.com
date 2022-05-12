package resource

import (
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"strings"
)

var nonResourceExtensions = []string{
	".html",
	".css",
	".png",
	".jpg",
	".jpeg",
}

func isResource(path string) bool {
	for _, nonResourceExtension := range nonResourceExtensions {
		if strings.HasSuffix(path, nonResourceExtension) {
			return false
		}
	}

	return true
}

func findResources(root string) (resourcePaths []string, err error) {
	resourcePaths = make([]string, 0)

	filePaths, err := file.FindFiles(root)
	if err != nil {
		return nil, err
	}

	for _, path := range filePaths {
		if isResource(path) {
			resourcePaths = append(resourcePaths, path)
		}
	}

	return resourcePaths, nil
}
