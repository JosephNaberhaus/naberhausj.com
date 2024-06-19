package image

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/builder"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"path/filepath"
	"strings"
)

var supportedImageExtensions = []string{
	".png",
	".jpg",
	".jpeg",
}

func createTargetWidths(originalWidth int) []int {
	baseTargetWidths := []int{100, 250, 500, 750, 1000, 1500, 2000}
	targetWidths := make([]int, 0, len(baseTargetWidths))
	for _, baseTargetWidth := range baseTargetWidths {
		if baseTargetWidth > originalWidth {
			targetWidths = append(targetWidths, originalWidth)
			break
		}

		targetWidths = append(targetWidths, baseTargetWidth)
	}
	return targetWidths
}

type handler struct {
	orchestrator builder.Orchestrator
}

func CreateHandler(orchestrator builder.Orchestrator) builder.Handler {
	return &handler{
		orchestrator: orchestrator,
	}
}

func (h *handler) CanCache() bool {
	return true
}

func (h *handler) DoesHandle(node *file.Node) bool {
	for _, supportedExtension := range supportedImageExtensions {
		if filepath.Ext(node.File) == supportedExtension {
			return true
		}
	}

	return false
}

func (h handler) Build(node *file.Node) (interface{}, error) {
	image, err := LoadImage(h.orchestrator.AbsPath(node))
	if err != nil {
		return nil, fmt.Errorf("error loading image file: %w", err)
	}

	originalWidth, originalHeight, err := image.Dimensions()
	if err != nil {
		return nil, fmt.Errorf("error getting image dimensions: %w", err)
	}

	var results []File
	for _, newWidth := range createTargetWidths(originalWidth) {
		if newWidth > originalWidth {
			continue
		}

		newHeight := (newWidth * originalHeight) / originalWidth
		result, err := image.resize(newWidth, newHeight)
		if err != nil {
			return nil, fmt.Errorf("error resizing image: %w", err)
		}

		newFile := fileName(node.File, result.extension, newWidth)
		err = h.orchestrator.Write(node, newFile, result.data)
		if err != nil {
			return nil, fmt.Errorf("error writing resized image: %w", err)
		}

		results = append(results, File{
			File:   newFile,
			Width:  newWidth,
			Height: newHeight,
		})
	}

	return Artifact{
		Files:          results,
		OriginalWidth:  originalWidth,
		OriginalHeight: originalHeight,
	}, nil
}

func fileName(originalFile, newExt string, width int) string {
	withoutExt := strings.TrimSuffix(originalFile, filepath.Ext(originalFile))
	return fmt.Sprintf("%s-%d%s", withoutExt, width, newExt)
}

func (h *handler) Finalize() error {
	return nil
}
