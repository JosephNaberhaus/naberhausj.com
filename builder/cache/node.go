package cache

type Node struct {
	File         string
	Hash         [20]byte
	Dependencies []string
	WrittenFiles []string
	Artifact     Artifact
}
