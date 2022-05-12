package orchestrator

import (
	"crypto/md5"
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/cache"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"io/ioutil"
)

func (b *Builder) computeCurNodeSet() error {
	filePaths, err := file.FindFiles(b.SourceDir)
	if err != nil {
		return fmt.Errorf("error find files for cur node set: %w", err)
	}

	nodes := make([]*cache.Node, 0, len(filePaths))
	for _, path := range filePaths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error loading file to compute hash: %w", err)
		}

		node := &cache.Node{
			Path: path,
			Hash: md5.Sum(data),
		}

		nodes = append(nodes, node)
	}

	b.CurNodeSet = cache.CreateNodeSet(nodes)
	for _, node := range nodes {
		dependencies, err := b.getNodeFileBuilder(node).FindDependencies(node)
		if err != nil {
			return fmt.Errorf("error getting dependencies: %w", err)
		}

		for _, dependency := range dependencies {
			dependencyNode := b.CurNodeSet.PathToNode[dependency]
			node.Dependencies = append(node.Dependencies, dependencyNode)
		}
	}

	return nil
}
