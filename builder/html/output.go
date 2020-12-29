package html

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/minifier"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Output(sourceDirectory, outputDirectory string, originalPath string, content string) error {
	if !strings.HasPrefix(originalPath, sourceDirectory) {
		return fmt.Errorf("file %s not located in source directory", originalPath)
	}

	newPath := strings.Replace(originalPath, sourceDirectory, outputDirectory, 1)

	info, err := os.Stat(originalPath)
	if err != nil {
		return fmt.Errorf("can't read file info: %w", err)
	}

	err = os.MkdirAll(filepath.Dir(newPath), info.Mode())
	if err != nil {
		return fmt.Errorf("error creating nested directories for %s: %w", newPath, err)
	}

	minifiedContent, err := minifier.Minifier.String("text/html", content)
	if err != nil {
		return fmt.Errorf("failed to minify HTML file: %w", err)
	}

	err = ioutil.WriteFile(newPath, []byte(minifiedContent), info.Mode())
	if err != nil {
		return fmt.Errorf("failed to write html file: %w", err)
	}

	return nil
}
