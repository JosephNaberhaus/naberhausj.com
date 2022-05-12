package build

import "github.com/JosephNaberhaus/naberhausj.com/builder/cache"

type ReaderWriter interface {
	Read(node *cache.Node) ([]byte, error)
	Write(node *cache.Node, file string, data []byte) error
}
