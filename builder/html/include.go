package html

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type documentNode interface {
	substituteParameters(map[string]string) documentNode
}

type textNode struct {
	content string
}

type includeNode struct {
	componentName string
	parameters    map[string]string
}

type document []documentNode

func (t *textNode) substituteParameters(parameters map[string]string) documentNode {
	result := new(textNode)
	result.content = substituteParameters(t.content, parameters)

	return result
}

func (t *includeNode) substituteParameters(parameters map[string]string) documentNode {
	result := new(includeNode)
	result.componentName = t.componentName
	result.parameters = make(map[string]string)

	for oldKey, oldValue := range t.parameters {
		result.parameters[oldKey] = substituteParameters(oldValue, parameters)
	}

	return result
}

func (d document) getIncludeNodes() []*includeNode {
	includeNodes := make([]*includeNode, 0)

	for _, node := range d {
		if node, ok := node.(*includeNode); ok {
			includeNodes = append(includeNodes, node)
		}
	}

	return includeNodes
}

func substituteParameters(value string, parameters map[string]string) string {
	result := value

	for name, value := range parameters {
		result = strings.ReplaceAll(result, "#"+name, value)
	}

	return result
}

func SplitIncludes(fileContent string) (document document, err error) {
	nodes := make([]documentNode, 0)

	directiveMatcher := regexp.MustCompile("<!--@(.*)({.*})-->")

	lastIndex := 0
	for _, indices := range directiveMatcher.FindAllStringSubmatchIndex(fileContent, -1) {
		textNode := new(textNode)
		textNode.content = fileContent[lastIndex:indices[0]]
		nodes = append(nodes, textNode)

		includeNode := new(includeNode)
		includeNode.componentName = fileContent[indices[2]:indices[3]]

		includeNode.parameters = make(map[string]string)
		err = json.Unmarshal([]byte(fileContent[indices[4]:indices[5]]), &includeNode.parameters)
		if err != nil {
			return nil, fmt.Errorf("error loading component paramters: %w", err)
		}

		nodes = append(nodes, includeNode)

		lastIndex = indices[1]
	}

	lastTextNode := new(textNode)
	lastTextNode.content = fileContent[lastIndex:len(fileContent)]
	nodes = append(nodes, lastTextNode)

	return nodes, nil
}
