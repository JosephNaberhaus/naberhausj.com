package orchestrator

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/cache"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"os"
	"path/filepath"
)

func cacheFile(root string) string {
	return filepath.Join(root, "build-cache.json")
}

func readCache(root string) (cache.NodeSet, error) {
	if !file.Exists(cacheFile(root)) {
		return cache.NodeSet{}, nil
	}

	c, err := cache.ReadNodeSet(cacheFile(root))
	if err != nil {
		return cache.NodeSet{}, fmt.Errorf("error loading cache: %w", err)
	}

	// Remove any nodes from the cache if one of their files has been deleted.
	for i, node := range c.Nodes {
		for _, writtenFile := range node.WrittenFiles {
			if !file.Exists(filepath.Join(root, writtenFile)) {
				c.Nodes = append(c.Nodes[:i], c.Nodes[i+1:]...)
				delete(c.PathToNode, node.File)
				break
			}
		}
	}

	return c, nil
}

func (b *Builder) writeNewCache() error {
	newCacheNodes := make([]*cache.Node, 0, len(b.current.Nodes))
	for _, node := range b.current.Nodes {
		cacheNode := &cache.Node{
			File:         node.File,
			Hash:         node.Hash,
			Dependencies: node.Dependencies,
			WrittenFiles: node.WrittenFiles,
			Artifact:     b.built[node],
		}
		newCacheNodes = append(newCacheNodes, cacheNode)
	}
	newCache := cache.CreateNodeSet(newCacheNodes)

	err := cache.WriteNodeSet(cacheFile(b.out), newCache)
	if err != nil {
		return fmt.Errorf("error writing the new cache: %w", err)
	}

	return nil
}

func (b *Builder) deleteWrittenFile(cacheNode *cache.Node) error {
	for _, writtenFile := range cacheNode.WrittenFiles {
		outputFile := filepath.Join(b.out, writtenFile)
		if file.Exists(outputFile) {
			err := os.Remove(outputFile)
			if err != nil {
				return fmt.Errorf("error deleting cached file: %w", err)
			}
		}
	}

	return nil
}
