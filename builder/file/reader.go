package file

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadNodes(root string) (NodeSet, error) {
	files, err := findFiles(root)
	if err != nil {
		return NodeSet{}, fmt.Errorf("error finding files for nodes: %w", err)
	}

	nodes := make([]*Node, 0, len(files))
	for _, file := range files {
		data, err := ioutil.ReadFile(filepath.Join(root, file))
		if err != nil {
			return NodeSet{}, fmt.Errorf("error reading data to compute hash: %w", err)
		}

		node := &Node{
			File: file,
			Hash: sha1.Sum(data),
		}
		nodes = append(nodes, node)
	}

	return CreateNodeSet(nodes), nil
}

func findFiles(root string) (paths []string, err error) {
	paths = make([]string, 0)

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking project files: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("error finding relative path for file: %w", err)
		}

		paths = append(paths, relativePath)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return paths, nil
}
