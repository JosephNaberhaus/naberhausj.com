package css

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Output(outputDirectory, cssContent string) error {
	err := os.MkdirAll(outputDirectory, 0777)
	if err != nil {
		return fmt.Errorf("couldn't create output directory: %w", err)
	}

	return ioutil.WriteFile(filepath.Join(outputDirectory, "styles.css"), []byte(cssContent), 0777)
}
