package html

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type DocumentCollection struct {
	documents     map[string]document
	nameToPathMap map[string]string
}

func (d *DocumentCollection) Substitute(toSubstitutePath string) (result string, err error) {
	return d.substituteParameters(toSubstitutePath, map[string]string{})
}

func (d *DocumentCollection) substituteParameters(toSubstitutePath string, parameters map[string]string) (result string, err error) {
	toSubstitute, exists := d.documents[toSubstitutePath]
	if !exists {
		return "", fmt.Errorf("couldn't find document at path: %s", toSubstitutePath)
	}

	resultBuilder := strings.Builder{}

	for _, node := range toSubstitute {
		substitutedNode := node.substituteParameters(parameters)

		switch v := substitutedNode.(type) {
		case textNode:
			resultBuilder.WriteString(string(v))
		case includeNode:
			path, exists := d.nameToPathMap[v.componentName]
			if !exists {
				return "", fmt.Errorf("couldn't find component with name: %s", v.componentName)
			}

			substitutionResult, err := d.substituteParameters(path, v.parameters)
			if err != nil {
				return "", err
			}

			resultBuilder.WriteString(substitutionResult)
		default:
			panic(fmt.Errorf("unkown type: %t", v))
		}
	}

	return resultBuilder.String(), nil
}

func SubstituteDate(toSubstitutePath string, content string) (string, error) {
	regex := regexp.MustCompile("<!--#created-on-->")
	if !regex.MatchString(content) {
		return content, nil
	}

	cmd := exec.Command("git", "log", "-1", "--diff-filter=A", "--format=%ad", "--date=iso-strict", toSubstitutePath)
	output := new(bytes.Buffer)
	cmd.Stdout = output
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	createdAt, err := time.Parse(time.RFC3339, strings.TrimSpace(output.String()))
	if err != nil {
		return "", err
	}

	createdAtStr := createdAt.Format("January 2006")
	result := regex.ReplaceAllString(content, createdAtStr)

	return result, nil
}

func BuildDocumentCollection(htmlFiles []*File, nameToPathMap map[string]string) (d *DocumentCollection, err error) {
	d = new(DocumentCollection)
	d.documents = make(map[string]document)
	d.nameToPathMap = nameToPathMap

	for _, file := range htmlFiles {
		document, err := SplitIncludes(file.content)
		if err != nil {
			return nil, err
		}

		d.documents[file.path] = document
	}

	return d, nil
}
