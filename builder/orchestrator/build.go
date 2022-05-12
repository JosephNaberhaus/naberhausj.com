package orchestrator

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/cache"
	"os"
)

func (b *Builder) Build() error {
	// Clear the build of any nodes that don't exist anymore
	for _, cacheNode := range b.CacheNodeSet.Nodes {
		if _, ok := b.CurNodeSet.PathToNode[cacheNode.Path]; !ok {
			err := b.clearNodeBuild(cacheNode)
			if err != nil {
				return fmt.Errorf("error cleaning up old build: %w", err)
			}
		}
	}

	for _, node := range b.CurNodeSet.Nodes {
		err := b.BuildNode(node)
		if err != nil {
			return fmt.Errorf("error building: %w", err)
		}
	}

	return nil
}

func (b *Builder) BuildNode(node *cache.Node) error {
	if !b.didNodeChange(node) {
		return nil
	}

	err := b.clearNodeBuild(b.CacheNodeSet.PathToNode[node.Path])
	if err != nil {
		return fmt.Errorf("error building node: %w", node)
	}

	for _, dependency := range node.Dependencies {
		err = b.BuildNode(dependency)
		if err != nil {
			return fmt.Errorf("error building dependency: %w", err)
		}
	}

	err = b.getNodeFileBuilder(node).BuildNode(node)
	if err != nil {
		return fmt.Errorf("error building node: %w", err)
	}

	return nil
}

func (b *Builder) clearNodeBuild(node *cache.Node) error {
	for _, writtenFile := range node.WrittenFiles {
		err := os.Remove(writtenFile)
		if err != nil {
			return fmt.Errorf("error remove previous file: %w", err)
		}
	}

	return nil
}
