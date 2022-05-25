package builder

import (
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
)

type Artifact = interface{}
type NodeSearchFunc = func(node *file.Node) bool

type Orchestrator interface {
	LoadDependency(node *file.Node, dependency string) (Artifact, error)
	FindDependencies(node *file.Node, search NodeSearchFunc) ([]Artifact, error)
	AbsPath(node *file.Node) string
	Write(node *file.Node, file string, data []byte) error
}
