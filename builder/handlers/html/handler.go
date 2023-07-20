package html

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/builder"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"io/ioutil"
	"strings"
)

const htmlExtension = ".html"
const componentExtension = ".component.html"

type handler struct {
	orchestrator builder.Orchestrator
}

func CreateHandler(orchestrator builder.Orchestrator) builder.Handler {
	return &handler{
		orchestrator: orchestrator,
	}
}

func (h *handler) DoesHandle(node *file.Node) bool {
	return strings.HasSuffix(node.File, htmlExtension) || strings.HasSuffix(node.File, componentExtension)
}

func (h *handler) CanCache() bool {
	return true
}

func (h *handler) Build(node *file.Node) (builder.Artifact, error) {
	data, err := ioutil.ReadFile(h.orchestrator.AbsPath(node))
	if err != nil {
		return nil, fmt.Errorf("error reading html file: %w", err)
	}

	content, err := substituteImports(node, h.orchestrator, toContent(string(data)))
	if err != nil {
		return nil, fmt.Errorf("error substituting content: %w", err)
	}

	if strings.HasSuffix(node.File, componentExtension) {
		return ComponentArtifact{
			Content: content,
		}, nil
	} else {
		content, err = substituteDirectives(node, h.orchestrator, content)
		if err != nil {
			return nil, fmt.Errorf("error substituting directives: %w", err)
		}

		err = h.orchestrator.Write(node, node.File, []byte(contentToString(content)))
		if err != nil {
			return nil, fmt.Errorf("error writing html file: %w", err)
		}

		return FileArtifact{
			File: node.File,
		}, nil
	}

}

func (h *handler) Finalize() error {
	return nil
}
