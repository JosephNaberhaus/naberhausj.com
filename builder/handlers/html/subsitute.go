package html

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/builder"
	"github.com/JosephNaberhaus/naberhausj.com/builder/cache"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"path/filepath"
)

func substituteImports(
	node *file.Node,
	orchestrator builder.Orchestrator,
	content []ContentNode,
) ([]ContentNode, error) {
	var newContent []ContentNode
	for _, contentNode := range content {
		if contentNode.isImport() {
			importedContent, err := substituteImport(node, orchestrator, contentNode)
			if err != nil {
				return nil, fmt.Errorf("error subsituting import: %w", err)
			}
			newContent = append(newContent, importedContent...)
		} else {
			newContent = append(newContent, contentNode)
		}
	}

	return newContent, nil
}

func substituteVariable(
	contentNode ContentNode,
	parameters map[string]string,
) (ContentNode, error) {
	variableName := contentNode.variableName()
	value, ok := parameters[variableName]
	if ok {
		return ContentNode(value), nil
	}
	return "", fmt.Errorf("no value provided for parameter: %s", variableName)
}

func substituteImport(
	node *file.Node,
	orchestrator builder.Orchestrator,
	contentNode ContentNode,
) ([]ContentNode, error) {
	componentName := contentNode.importName()
	artifacts, err := orchestrator.FindDependencies(
		node,
		func(node *file.Node) bool {
			return componentName+componentExtension == filepath.Base(node.File)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error substituting import: %w", err)
	}
	if len(artifacts) != 1 {
		return nil, fmt.Errorf("found %d artifacts with the name \"%s\", expected 1", len(artifacts), componentName)
	}

	var artifact ComponentArtifact
	err = cache.UnmarshalArtifact(artifacts[0], &artifact)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling artifact: %w", err)
	}

	parameters, err := contentNode.importParameters()
	if err != nil {
		return nil, fmt.Errorf("error getting parameters for import: %w", err)
	}

	return substituteContent(artifact.Content, parameters)
}

func substituteContent(
	content []ContentNode,
	parameters map[string]string,
) ([]ContentNode, error) {
	var newContent []ContentNode
	for _, artifactContent := range content {
		if artifactContent.isVariable() {
			result, err := substituteVariable(artifactContent, parameters)
			if err != nil {
				return nil, fmt.Errorf("error subsituting variable: %w", err)
			}
			newContent = append(newContent, result)
		} else if artifactContent.isDirective() {
			directiveName := artifactContent.directiveName()
			directiveParameters, err := artifactContent.directiveParameters()
			if err != nil {
				return nil, err
			}

			substitutedParameters := map[string]string{}
			for key, value := range directiveParameters {
				valueContent := toContent(value)
				substitutedContent, err := substituteContent(valueContent, parameters)
				if err != nil {
					return nil, fmt.Errorf("error substituting value content in directive: %w", err)
				}

				substitutedParameters[key] = contentToString(substitutedContent)
			}

			fmt.Printf("%s %v\n", directiveName, substitutedParameters)

			substitutedDirective, err := newDirective(directiveName, substitutedParameters)
			if err != nil {
				return nil, fmt.Errorf("error creating a new directive: %w", err)
			}

			newContent = append(newContent, substitutedDirective)
		} else {
			newContent = append(newContent, artifactContent)
		}
	}

	return newContent, nil
}

func substituteDirectives(
	node *file.Node,
	orchestrator builder.Orchestrator,
	content []ContentNode,
) ([]ContentNode, error) {
	var newContent []ContentNode
	for _, contentNode := range content {
		if contentNode.isDirective() {
			directiveName := contentNode.directiveName()
			directiveHandler, ok := directiveHandlers[directiveName]
			if !ok {
				return nil, fmt.Errorf("unkown directive: %s", directiveName)
			}

			parameters, err := contentNode.directiveParameters()
			if err != nil {
				return nil, fmt.Errorf("error getting parameters for directive: %w", err)
			}

			result, err := directiveHandler(node, parameters, orchestrator)
			if err != nil {
				return nil, fmt.Errorf("error handling directive: %w", err)
			}

			newContent = append(newContent, result...)
		} else {
			newContent = append(newContent, contentNode)
		}
	}

	return newContent, nil
}
