package html

import (
	"errors"
	"fmt"
	"strings"
)

type substitutionState int8

const (
	stateUnsubstituted = iota
	stateWaitingOnChildren
	stateSubstituted
)

type dependencyEdge struct {
	componentName string
	to            *dependencyNode
}

type dependencyNode struct {
	htmlFile         *File
	dependents       []*dependencyNode
	dependencies     []dependencyEdge
	substitutedState substitutionState
}

type DependencyGraph map[string]*dependencyNode

func (g DependencyGraph) Substitute(toSubstitute *dependencyNode) error {
	if toSubstitute.substitutedState == stateSubstituted {
		return nil
	}

	if toSubstitute.substitutedState == stateWaitingOnChildren {
		return errors.New("include loop between components")
	}

	toSubstitute.substitutedState = stateWaitingOnChildren

	if len(toSubstitute.dependencies) == 0 {
		return nil
	}

	for _, dependency := range toSubstitute.dependencies {
		err := g.Substitute(dependency.to)
		if err != nil {
			return err
		}

		toSubstitute.htmlFile.content = strings.ReplaceAll(
			toSubstitute.htmlFile.content,
			createIncludeDirective(dependency.componentName),
			dependency.to.htmlFile.content,
		)
	}

	toSubstitute.substitutedState = stateSubstituted

	return nil
}

func (g DependencyGraph) SubstitutePath(path string) error {
	node, exists := g[path]
	if !exists {
		return fmt.Errorf("invalid path: %s", path)
	}

	return g.Substitute(node)
}

func BuildDependencyGraph(files []*File, nameToPathMap map[string]string) (graph DependencyGraph, err error) {
	graph = make(DependencyGraph)

	for _, file := range files {
		node := new(dependencyNode)
		node.htmlFile = file
		node.dependencies = make([]dependencyEdge, 0)
		node.dependents = make([]*dependencyNode, 0)
		node.substitutedState = stateUnsubstituted

		graph[file.path] = node
	}

	for _, file := range files {
		fileNode := graph[file.path]

		for _, include := range FindIncludes(file.content) {
			includePath, exists := nameToPathMap[include]
			if !exists {
				return nil, fmt.Errorf("could not find component named: %s", include)
			}

			includeNode, exists := graph[includePath]
			if !exists {
				return nil, fmt.Errorf("could not find template file at: %s", includePath)
			}

			includeNode.dependents = append(includeNode.dependents, fileNode)
			fileNode.dependencies = append(fileNode.dependencies, dependencyEdge{
				componentName: include,
				to:            includeNode,
			})
		}
	}

	return graph, nil
}
