package file

import (
	"path/filepath"
	"strings"
)

func ResolveHTMLPath(path, htmlFilePath, root string) string {
	if strings.HasPrefix(path, "/") {
		return filepath.Join(root, path)
	}

	htmlDirectory := filepath.Dir(htmlFilePath)
	return filepath.Join(htmlDirectory, path)
}
