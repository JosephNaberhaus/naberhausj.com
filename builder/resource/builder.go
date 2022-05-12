package resource

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/build"
	"github.com/JosephNaberhaus/naberhausj.com/builder/cache"
)

type Builder struct {
	ReaderWriter build.ReaderWriter
}

func (b *Builder) CanHandleFile(_ *cache.Node) bool {
	return true
}

func (b *Builder) FindDependencies(_ *cache.Node) ([]string, error) {
	return nil, nil
}

func (b *Builder) BuildNode(node *cache.Node) error {
	data, err := b.ReaderWriter.Read(node)
	if err != nil {
		return fmt.Errorf("error reading resource: %w", err)
	}

	err = b.ReaderWriter.Write(node, node.Path, data)
	if err != nil {
		return fmt.Errorf("error writing resource: %w", err)
	}

	return nil
}
