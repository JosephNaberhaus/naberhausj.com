package resource

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/builder"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"github.com/JosephNaberhaus/naberhausj.com/builder/minifier"
	"io/ioutil"
	"path/filepath"
)

const javascriptFileExtension = ".js"

type handler struct {
	orchestrator builder.Orchestrator
}

func CreateHandler(orchestrator builder.Orchestrator) builder.Handler {
	return &handler{
		orchestrator: orchestrator,
	}
}

func (h *handler) DoesHandle(_ *file.Node) bool {
	return true
}

func (h *handler) CanCache() bool {
	return true
}

func (h *handler) Build(node *file.Node) (builder.Artifact, error) {
	data, err := ioutil.ReadFile(h.orchestrator.AbsPath(node))
	if err != nil {
		return nil, fmt.Errorf("error reading resource file: %w", err)
	}

	if filepath.Ext(node.File) == javascriptFileExtension {
		data, err = minifier.Global.Bytes("text/javascript", data)
		if err != nil {
			return nil, fmt.Errorf("error minifying javascript: %w", err)
		}
	}

	err = h.orchestrator.Write(node, node.File, data)
	if err != nil {
		return nil, fmt.Errorf("error writing resource node: %w", err)
	}

	return Artifact{
		File: node.File,
	}, nil
}

func (h *handler) Finalize() error {
	return nil
}
