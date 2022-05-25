package builder

import (
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
)

type Handler interface {
	CanCache() bool
	DoesHandle(node *file.Node) bool
	Build(node *file.Node) (interface{}, error)
	Finalize() error
}
