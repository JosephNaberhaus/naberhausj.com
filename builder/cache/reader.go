package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ReadNodeSet(file string) (NodeSet, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return NodeSet{}, fmt.Errorf("error reading node set: %w", err)
	}

	var nodes []*Node
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		return NodeSet{}, fmt.Errorf("error unmarshalling node set: %w", err)
	}

	return CreateNodeSet(nodes), nil
}
