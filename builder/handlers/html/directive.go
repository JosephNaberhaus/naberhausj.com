package html

import (
	"bytes"
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/builder"
	"github.com/JosephNaberhaus/naberhausj.com/builder/cache"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"github.com/JosephNaberhaus/naberhausj.com/builder/handlers/image"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type directiveHandlerFunc = func(*file.Node, map[string]string, builder.Orchestrator) ([]ContentNode, error)

var directiveHandlers = map[string]directiveHandlerFunc{
	"created-on": handleCreatedOnDirective,
	"img":        handleImageDirective,
}

func handleCreatedOnDirective(
	node *file.Node,
	_ map[string]string,
	orchestrator builder.Orchestrator,
) ([]ContentNode, error) {
	absPath := orchestrator.AbsPath(node)
	cmd := exec.Command("git", "log", "-1", "--diff-filter=A", "--format=%ad", "--date=iso-strict", absPath)
	output := new(bytes.Buffer)
	cmd.Stdout = output
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error getting file creation date: %w", err)
	}

	outputStr := strings.TrimSpace(output.String())
	if outputStr == "" {
		// The file hasn't been committed yet. Ignore it.
		return nil, nil
	}

	createdAt, err := time.Parse(time.RFC3339, outputStr)
	if err != nil {
		return nil, err
	}

	createdAtNode := ContentNode(createdAt.Format("January 2006"))
	return []ContentNode{createdAtNode}, nil
}

func handleImageDirective(
	node *file.Node,
	parameters map[string]string,
	orchestrator builder.Orchestrator,
) ([]ContentNode, error) {
	src, ok := parameters["src"]
	if !ok {
		return nil, fmt.Errorf("no src provided for image directive at: %s", node.File)
	}

	artifact, err := orchestrator.LoadDependency(node, src)
	if err != nil {
		return nil, fmt.Errorf("error loading image dependency: %w", err)
	}

	var imageArtifact image.Artifact
	err = cache.UnmarshalArtifact(artifact, &imageArtifact)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling image artifact: %w", err)
	}

	width, ok := parseIntParameter(parameters, "width")
	if !ok {
		return nil, fmt.Errorf("no valid width provided for image directive at: %s", node.File)
	}
	height := (width * imageArtifact.OriginalHeight) / imageArtifact.OriginalWidth

	html, err := htmlImage{
		class:  parameters["class"],
		width:  width,
		height: height,
		src:    src,
		files:  imageArtifact.Files,
		alt:    parameters["alt"],
	}.toHtml(filepath.Dir(node.File))
	if err != nil {
		return nil, fmt.Errorf("error creating html for image: %w", err)
	}

	return []ContentNode{ContentNode(html)}, nil
}

func parseIntParameter(parameters map[string]string, key string) (int, bool) {
	valueStr, ok := parameters[key]
	if !ok {
		return 0, false
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, false
	}

	return value, true
}
