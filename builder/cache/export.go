package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type exportNode struct {
	path         string
	hash         [16]byte
	dependencies []string

	writtenFiles []string
}

func Export(outputFile string, set NodeSet) error {
	exportedNodes := make([]exportNode, 0, len(set.Nodes))
	for _, node := range set.Nodes {
		dependencies := make([]string, 0, len(node.Dependencies))
		for _, dependency := range node.Dependencies {
			dependencies = append(dependencies, dependency.Path)
		}

		exportedNodes = append(exportedNodes, exportNode{
			path:         node.Path,
			hash:         node.Hash,
			dependencies: dependencies,
			writtenFiles: node.WrittenFiles,
		})
	}

	export, err := json.Marshal(exportedNodes)
	if err != nil {
		return fmt.Errorf("error marshalling export: %w", err)
	}

	err = ioutil.WriteFile(outputFile, export, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing export: %w", err)
	}

	return nil
}

func Import(file string) (NodeSet, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return NodeSet{}, fmt.Errorf("error importing cache file: %w", err)
	}

	var exportedNodes []exportNode
	err = json.Unmarshal(data, &exportedNodes)
	if err != nil {
		return NodeSet{}, fmt.Errorf("error unmarshaling cache: %w", err)
	}

	nodes := make([]*Node, 0, len(exportedNodes))
	pathToNode := map[string]*Node{}
	for _, exportedNode := range exportedNodes {
		node := &Node{
			Path:         exportedNode.path,
			Hash:         exportedNode.hash,
			Dependencies: make([]*Node, 0, len(exportedNode.dependencies)),
			WrittenFiles: exportedNode.writtenFiles,
		}

		nodes = append(nodes, node)
		pathToNode[node.Path] = node
	}

	for _, exportedNode := range exportedNodes {
		for _, dependency := range exportedNode.dependencies {
			pathToNode[exportedNode.path] = pathToNode[dependency]
		}
	}

	return CreateNodeSet(nodes), nil
}
