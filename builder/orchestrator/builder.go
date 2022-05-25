package orchestrator

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/builder"
	"github.com/JosephNaberhaus/naberhausj.com/builder/cache"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Builder struct {
	src, out string

	handlers []builder.Handler

	cache   cache.NodeSet
	current file.NodeSet
	built   map[*file.Node]builder.Artifact
}

func CreateBuilder(src, out string) (*Builder, error) {
	cache, err := readCache(out)
	if err != nil {
		return nil, fmt.Errorf("error loading cache: %w", err)
	}

	current, err := file.ReadNodes(src)
	if err != nil {
		return nil, fmt.Errorf("error loading current: %w", err)
	}

	return &Builder{
		src:     src,
		out:     out,
		cache:   cache,
		current: current,
		built:   map[*file.Node]builder.Artifact{},
	}, nil
}

func (b *Builder) AddHandler(handler builder.Handler) {
	b.handlers = append(b.handlers, handler)
}

func (b *Builder) GetNode(file string) (*cache.Node, *file.Node, error) {
	cacheNode, ok := b.cache.PathToNode[file]
	if !ok {
		return nil, nil, fmt.Errorf("error looking up cache node at '%s'", file)
	}

	node, ok := b.current.PathToNode[file]
	if !ok {
		return nil, nil, fmt.Errorf("error looking up node at: '%s'", file)
	}

	return cacheNode, node, nil
}

func (b *Builder) LoadDependency(node *file.Node, dependency string) (builder.Artifact, error) {
	dependencyNode, ok := b.current.PathToNode[dependency]
	if !ok {
		return nil, fmt.Errorf("dependency node not found: %s", dependency)
	}

	if node == dependencyNode {
		return nil, fmt.Errorf("node attempted to load itself as a dependency")
	}

	artifact, err := b.buildNode(dependencyNode)
	if err != nil {
		return nil, fmt.Errorf("error building dependency node: %w", err)
	}

	node.Dependencies = append(node.Dependencies, dependency)
	return artifact, nil
}

func (b *Builder) FindDependencies(
	node *file.Node,
	filterFunc builder.NodeSearchFunc,
) ([]builder.Artifact, error) {
	var artifacts []builder.Artifact
	for _, toCheck := range b.current.Nodes {
		// A node can't depend on itself
		if node == toCheck {
			continue
		}

		if filterFunc(toCheck) {
			artifact, err := b.LoadDependency(node, toCheck.File)
			if err != nil {
				return nil, fmt.Errorf("error loading a dependency: %w", err)
			}
			artifacts = append(artifacts, artifact)
		}
	}

	return artifacts, nil
}

func (b *Builder) AbsPath(node *file.Node) string {
	return filepath.Join(b.src, node.File)
}

func (b *Builder) Write(node *file.Node, file string, data []byte) error {
	outputFile := filepath.Join(b.out, file)

	err := os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory for file: %w", err)
	}

	err = ioutil.WriteFile(outputFile, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing file: %s", err)
	}

	if node != nil {
		node.WrittenFiles = append(node.WrittenFiles, file)
	}
	return nil
}

func (b *Builder) buildNode(node *file.Node) (builder.Artifact, error) {
	if artifact, isBuilt := b.built[node]; isBuilt {
		return artifact, nil
	}

	handler := b.findHandlerForNode(node)

	cacheNode, isCached := b.cache.PathToNode[node.File]
	if handler.CanCache() && isCached {
		if !b.ShouldRebuild(node) {
			b.built[node] = cacheNode.Artifact
			node.Dependencies = cacheNode.Dependencies
			node.WrittenFiles = cacheNode.WrittenFiles
			return cacheNode.Artifact, nil
		}

		err := b.deleteWrittenFile(cacheNode)
		if err != nil {
			return nil, fmt.Errorf("error deleting written files before rebuilding: %w", err)
		}
	}

	fmt.Printf("Building: %s\n", node.File)
	artifact, err := handler.Build(node)
	if err != nil {
		return nil, fmt.Errorf("error building node: %w", err)
	}

	b.built[node] = artifact
	return artifact, nil
}

func (b *Builder) findHandlerForNode(node *file.Node) builder.Handler {
	for _, handler := range b.handlers {
		if handler.DoesHandle(node) {
			return handler
		}
	}

	return nil
}

func (b *Builder) removeDeleted() error {
	for _, cacheNode := range b.cache.Nodes {
		if _, exists := b.current.PathToNode[cacheNode.File]; !exists {
			err := b.deleteWrittenFile(cacheNode)
			if err != nil {
				return fmt.Errorf("error removing files for a deleted node: %w", err)
			}
		}
	}

	return nil
}

func (b *Builder) Build() error {
	err := b.removeDeleted()
	if err != nil {
		return fmt.Errorf("error removing deleted: %w", err)
	}

	for _, node := range b.current.Nodes {
		_, err = b.buildNode(node)
		if err != nil {
			return fmt.Errorf("error building node: %w", err)
		}
	}

	for _, handler := range b.handlers {
		err = handler.Finalize()
		if err != nil {
			return fmt.Errorf("error finalizing handler: %w", err)
		}
	}

	err = b.writeNewCache()
	if err != nil {
		return fmt.Errorf("error writing new cache: %w", err)
	}

	return nil
}
