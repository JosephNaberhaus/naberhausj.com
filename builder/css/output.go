package css

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/minifier"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Output(outputDirectory, content string) error {
	err := os.MkdirAll(outputDirectory, 0777)
	if err != nil {
		return fmt.Errorf("couldn't create output directory: %w", err)
	}

	minifiedContent, err := minifier.Minifier.String("text/css", content)
	if err != nil {
		return fmt.Errorf("failed to minify CSS file: %w", err)
	}

	return ioutil.WriteFile(filepath.Join(outputDirectory, "styles.css"), []byte(minifiedContent), 0777)
}
