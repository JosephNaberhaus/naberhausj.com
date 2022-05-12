package orchestrator

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/build"
	"github.com/JosephNaberhaus/naberhausj.com/builder/cache"
	"github.com/JosephNaberhaus/naberhausj.com/builder/resource"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Builder struct {
	SourceDir, OutputDir string

	BuiltPaths               map[string]struct{}
	CacheNodeSet, CurNodeSet cache.NodeSet
	FileBuilders             []build.FileBuilder
}

func CreateBuilder(sourceDir, outputDir string) (*Builder, error) {
	builder := &Builder{
		SourceDir:  sourceDir,
		OutputDir:  outputDir,
		BuiltPaths: map[string]struct{}{},
	}

	builder.FileBuilders = []build.FileBuilder{
		&resource.Builder{ReaderWriter: builder},
	}

	err := builder.computeCurNodeSet()
	if err != nil {
		return nil, fmt.Errorf("error computing cur node set: %w", err)
	}

	return builder, nil
}

func (b *Builder) Read(n *cache.Node) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath.Join(b.SourceDir, n.Path))
	if err != nil {
		return nil, fmt.Errorf("error loading node content: %w", err)
	}

	return data, nil
}

func (b *Builder) Write(node *cache.Node, file string, data []byte) error {
	outputFile := filepath.Join(b.OutputDir, file)

	err := os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directories before writing: %w", err)
	}

	err = ioutil.WriteFile(outputFile, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing node: %w", err)
	}

	node.WrittenFiles = append(node.WrittenFiles, file)
	return nil
}

func (b *Builder) getNodeFileBuilder(n *cache.Node) build.FileBuilder {
	for _, fileBuilder := range b.FileBuilders {
		if fileBuilder.CanHandleFile(n) {
			return fileBuilder
		}
	}

	return nil
}

func (b *Builder) didNodeChange(node *cache.Node) bool {
	_, isBuilt := b.BuiltPaths[node.Path]
	if isBuilt {
		return false
	}

	cacheNode, ok := b.CacheNodeSet.PathToNode[node.Path]
	if !ok {
		return true
	}

	return node.Equals(cacheNode)
}
