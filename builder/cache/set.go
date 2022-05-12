package cache

type NodeSet struct {
	Nodes      []*Node
	PathToNode map[string]*Node
}

func CreateNodeSet(nodes []*Node) NodeSet {
	nodeSet := NodeSet{
		Nodes:      make([]*Node, 0, len(nodes)),
		PathToNode: map[string]*Node{},
	}

	for _, node := range nodes {
		nodeSet.Nodes = append(nodeSet.Nodes, node)
		nodeSet.PathToNode[node.Path] = node
	}

	return nodeSet
}
