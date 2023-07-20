package css

import (
	"bytes"
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/builder"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const cssExtension = ".css"
const scopedCssExtension = ".scoped.css"
const outputFile = "styles.css"

type handler struct {
	orchestrator builder.Orchestrator
	cssBuilders  map[string]*bytes.Buffer
}

func CreateHandler(orchestrator builder.Orchestrator) builder.Handler {
	return &handler{
		orchestrator: orchestrator,
		cssBuilders:  map[string]*bytes.Buffer{},
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

	cssBuilderKey := ""
	if strings.HasSuffix(node.File, scopedCssExtension) {
		cssBuilderKey = filepath.Dir(node.File)
	}

	cssBuilder := h.cssBuilders[cssBuilderKey]
	if cssBuilder == nil {
		cssBuilder = new(bytes.Buffer)
		h.cssBuilders[cssBuilderKey] = cssBuilder
	}

	cssBuilder.Write(data)
	cssBuilder.WriteRune('\n')
	return nil, err
}

func (h *handler) Finalize() error {
	for dir, cssBuilder := range h.cssBuilders {
		err := h.orchestrator.Write(nil, filepath.Join(dir, outputFile), cssBuilder.Bytes())
		if err != nil {
			return fmt.Errorf("error writing css: %w", err)
		}
	}

	return nil
}
