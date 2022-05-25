package html

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/handlers/image"
	"path/filepath"
	"strings"
)

type htmlImage struct {
	class         string
	width, height int
	src           string
	files         []image.File
	alt           string
}

func (h htmlImage) toHtml(dir string) (string, error) {
	srcSetBuilder := strings.Builder{}
	for i, file := range h.files {
		relPath, err := filepath.Rel(dir, file.File)
		if err != nil {
			return "", fmt.Errorf("error creating relative path for image: %w", err)
		}

		srcSetBuilder.WriteString(fmt.Sprintf("%s %dw", relPath, file.Width))
		if i != len(h.files)-1 {
			srcSetBuilder.WriteRune(',')
		}
	}

	return fmt.Sprintf(
		"<img class=\"%s\" width=\"%d\" height=\"%d\" src=\"%s\" srcset=\"%s\" sizes=\"%dpx\" alt=\"%s\">",
		h.class,
		h.width,
		h.height,
		h.src,
		srcSetBuilder.String(),
		h.width,
		h.alt,
	), nil
}
