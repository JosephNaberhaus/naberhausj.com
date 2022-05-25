package css

import (
	"bytes"
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/builder"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"github.com/JosephNaberhaus/naberhausj.com/builder/minifier"
	"io/ioutil"
	"path/filepath"
)

const cssExtension = ".css"
const outputFile = "styles.css"

type handler struct {
	orchestrator builder.Orchestrator
	cssBuilder   bytes.Buffer
}

func CreateHandler(orchestrator builder.Orchestrator) builder.Handler {
	return &handler{
		orchestrator: orchestrator,
	}
}

func (h *handler) CanCache() bool {
	return false
}

func (h *handler) DoesHandle(node *file.Node) bool {
	return filepath.Ext(node.File) == cssExtension
}

func (h *handler) Build(node *file.Node) (interface{}, error) {
	data, err := ioutil.ReadFile(h.orchestrator.AbsPath(node))
	if err != nil {
		return nil, fmt.Errorf("error loading css file: %w", err)
	}

	h.cssBuilder.Write(data)
	h.cssBuilder.WriteRune('\n')
	return nil, err
}

func (h *handler) Finalize() error {
	minifiedCSS, err := minifier.Global.Bytes("text/css", h.cssBuilder.Bytes())
	if err != nil {
		return fmt.Errorf("error minifying css: %w", err)
	}

	err = h.orchestrator.Write(nil, outputFile, minifiedCSS)
	if err != nil {
		return fmt.Errorf("error writing css: %w", err)
	}

	return nil
}
