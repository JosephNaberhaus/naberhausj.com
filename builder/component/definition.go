package component

import (
	"encoding/json"
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Collection []*Definition

type Definition struct {
	Name         string   `json:"name"`
	TemplatePath string   `json:"templatePath"`
	Resources    []string `json:"resources"`
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

func (c Collection) GetAllResourcePaths() []string {
	resourcePaths := make([]string, 0)

	for _, definition := range c {
		if definition.Resources != nil {
			resourcePaths = append(resourcePaths, definition.Resources...)
		}
	}

	return resourcePaths
}

// Walks the entire file system from the working directory down looking for component definitions files
func FindDefinitions(root string) (definitions Collection, err error) {
	definitionPaths, err := file.FindFilesWithSuffix(root, ".component.json")
	definitions = make(Collection, 0, len(definitionPaths))

	for _, path := range definitionPaths {
		definitionBytes, err := readFile(path)
		if err != nil {
			return nil, fmt.Errorf("error loading component definition: %w", err)
		}

		definition := new(Definition)
		err = json.Unmarshal(definitionBytes, definition)
		if err != nil {
			return nil, fmt.Errorf("error reading definition json: %w", err)
		}

		rootPath := filepath.Dir(path)
		definition.TemplatePath = filepath.Join(rootPath, definition.TemplatePath)

		if definition.Resources != nil {
			for i, path := range definition.Resources {
				definition.Resources[i] = filepath.Join(rootPath, path)
			}
		}

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
