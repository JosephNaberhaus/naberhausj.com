package component

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Collection []*Definition

type Definition struct {
	Name         string `json:"name"`
	TemplatePath string `json:"templatePath"`
}

func (c Collection) CreateNameToPathMap() (nameToPathMap map[string]string, err error) {
	nameToPathMap = make(map[string]string)

	for _, definition := range c {
		if _, exists := nameToPathMap[definition.Name]; exists {
			return nil, fmt.Errorf("duplicate component name: %s", definition.Name)
		}

		nameToPathMap[definition.Name] = definition.TemplatePath
	}

	return nameToPathMap, nil
}

func (c Collection) ContainsPath(path string) bool {
	for _, definition := range c {
		if definition.TemplatePath == path {
			return true
		}
	}

	return false
}

// Walks the entire file system from the working directory down looking for component definitions files
func FindDefinitions(root string) (definitions Collection, err error) {
	definitionPaths, err := file.FindFilesWithSuffix(root, ".component.html")
	definitions = make(Collection, 0, len(definitionPaths))

	for _, path := range definitionPaths {
		definition := new(Definition)

		definition.Name = strings.TrimSuffix(filepath.Base(path), ".component.html")
		definition.TemplatePath = path

		definitions = append(definitions, definition)
	}

	return definitions, nil
}

func readFile(path string) (result []byte, err error) {
	toRead, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(toRead)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
