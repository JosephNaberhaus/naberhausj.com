package html

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Output(sourceDirectory, outputDirectory string, htmlFiles []*File) error {
	for _, file := range htmlFiles {
		if !strings.HasPrefix(file.path, sourceDirectory) {
			return fmt.Errorf("file %s not located in source directory", file.path)
		}

		newPath := strings.Replace(file.path, sourceDirectory, outputDirectory, 1)

		info, err := os.Stat(file.path)
		if err != nil {
			return fmt.Errorf("can't read file info: %w", err)
		}

		err = os.MkdirAll(filepath.Dir(newPath), info.Mode())
		if err != nil {
			return fmt.Errorf("error creating nested directories for %s: %w", newPath, err)
		}

		err = ioutil.WriteFile(newPath, []byte(file.content), info.Mode())
		if err != nil {
			return fmt.Errorf("failed to write html file: %w", err)
		}
	}

	return nil
}
