package file

type NodeSet struct {
	Nodes      []*Node
	PathToNode map[string]*Node
}

func CreateNodeSet(nodes []*Node) NodeSet {
	pathToNodes := map[string]*Node{}
	for _, node := range nodes {
		pathToNodes[node.File] = node
	}

	return NodeSet{
		Nodes:      nodes,
		PathToNode: pathToNodes,
	}
}
