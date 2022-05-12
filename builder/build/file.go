package build

import "github.com/JosephNaberhaus/naberhausj.com/builder/cache"

type FileBuilder interface {
	CanHandleFile(node *cache.Node) bool
	FindDependencies(node *cache.Node) ([]string, error)
	BuildNode(node *cache.Node) error
}
