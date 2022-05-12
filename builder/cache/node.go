package cache

type Node struct {
	Path         string
	Hash         [16]byte
	Dependencies []*Node

	WrittenFiles []string
}

func (n *Node) Equals(other *Node) bool {
	if n.Path != other.Path {
		return false
	}

	if n.Hash != other.Hash {
		return false
	}

	if len(n.Dependencies) != len(other.Dependencies) {
		return false
	}

	for i := range n.Dependencies {
		if !n.Dependencies[i].Equals(other.Dependencies[i]) {
			return false
		}
	}

	return true
}
