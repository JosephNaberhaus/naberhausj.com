package html

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var variableRegex = regexp.MustCompile("#([a-zA-Z0-9]+)")
var onlyVariableRegex = toFullMatch(variableRegex)
var importRegex = regexp.MustCompile("<!--@(.*)({.*})-->")
var onlyImportRegex = toFullMatch(importRegex)
var directiveRegex = regexp.MustCompile("<!--#(.*)({.*})-->")
var onlyDirectiveRegex = toFullMatch(directiveRegex)

var allTokensRegex = or(variableRegex, importRegex, directiveRegex)

func toFullMatch(regex *regexp.Regexp) *regexp.Regexp {
	return regexp.MustCompile("^" + regex.String() + "$")
}

func or(regexes ...*regexp.Regexp) *regexp.Regexp {
	combined := strings.Builder{}
	for i, regex := range regexes {
		combined.WriteRune('(')
		combined.WriteString(regex.String())
		combined.WriteRune(')')
		if i != len(regexes)-1 {
			combined.WriteRune('|')
		}
	}

	return regexp.MustCompile(combined.String())
}

type ContentNode string

func newDirective(name string, parameters map[string]string) (ContentNode, error) {
	jsonParameters, err := json.Marshal(parameters)
	if err != nil {
		return "", err
	}

	return ContentNode(fmt.Sprintf("<!--#%s%s-->", name, string(jsonParameters))), nil
}

func toContent(value string) []ContentNode {
	var content []ContentNode

	prevIndex := 0
	for _, match := range allTokensRegex.FindAllStringIndex(value, -1) {
		startIndex, endIndex := match[0], match[1]
		if startIndex != prevIndex {
			content = append(content, ContentNode(value[prevIndex:startIndex]))
		}

		content = append(content, ContentNode(value[startIndex:endIndex]))
		prevIndex = endIndex
	}
	if prevIndex != len(value)-1 {
		content = append(content, ContentNode(value[prevIndex:]))
	}

	return content
}

func (c ContentNode) isVariable() bool {
	return onlyVariableRegex.MatchString(string(c))
}

func (c ContentNode) variableName() string {
	return onlyVariableRegex.FindStringSubmatch(string(c))[1]
}

func (c ContentNode) isImport() bool {
	return onlyImportRegex.MatchString(string(c))
}

func (c ContentNode) importName() string {
	return onlyImportRegex.FindStringSubmatch(string(c))[1]
}

func (c ContentNode) importParameters() (map[string]string, error) {
	parametersStr := onlyImportRegex.FindStringSubmatch(string(c))[2]
	return parseParameters(parametersStr)
}

func (c ContentNode) isDirective() bool {
	return onlyDirectiveRegex.MatchString(string(c))
}

func (c ContentNode) directiveName() string {
	return onlyDirectiveRegex.FindStringSubmatch(string(c))[1]
}

func (c ContentNode) directiveParameters() (map[string]string, error) {
	parametersStr := onlyDirectiveRegex.FindStringSubmatch(string(c))[2]
	return parseParameters(parametersStr)
}

func parseParameters(parametersStr string) (map[string]string, error) {
	parameters := map[string]string{}
	err := json.Unmarshal([]byte(parametersStr), &parameters)
	if err != nil {
		return nil, fmt.Errorf("error parsing parameters: %w", err)
	}

	return parameters, nil
}

func contentToString(content []ContentNode) string {
	sb := strings.Builder{}
	for _, node := range content {
		sb.WriteString(string(node))
	}
	return sb.String()
}
