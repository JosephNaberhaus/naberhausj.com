package orchestrator

import (
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
)

func (b *Builder) ShouldRebuild(node *file.Node) bool {
	cacheNode, ok := b.cache.PathToNode[node.File]
	if !ok {
		return true
	}

	if cacheNode.Hash != node.Hash {
		return true
	}

	for _, dependency := range cacheNode.Dependencies {
		dependencyNode, ok := b.current.PathToNode[dependency]
		if !ok {
			return true
		}

		if b.ShouldRebuild(dependencyNode) {
			return true
		}
	}

	return false
}
