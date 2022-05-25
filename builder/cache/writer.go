package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteNodeSet(file string, set NodeSet) error {
	err := os.MkdirAll(filepath.Dir(file), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory for node set: %w", err)
	}

	data, err := json.Marshal(set.Nodes)
	if err != nil {
		return fmt.Errorf("error marshalling node set: %w", err)
	}

	err = ioutil.WriteFile(file, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing node set: %w", err)
	}

	return nil
}
