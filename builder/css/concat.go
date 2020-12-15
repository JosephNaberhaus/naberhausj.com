package css

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"io/ioutil"
	"strings"
)

func ConcatFiles(root string) (result string, err error) {
	cssFilePaths, err := file.FindFilesWithSuffix(root, ".css")
	if err != nil {
		return "", fmt.Errorf("error finding CSS files: %w", err)
	}

	resultBuilder := strings.Builder{}

	for _, path := range cssFilePaths {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("error loading css file: %w", err)
		}

		resultBuilder.Write(content)
		resultBuilder.WriteRune('\n')
	}

	return resultBuilder.String(), err
}
