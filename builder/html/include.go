package html

import (
	"fmt"
	"regexp"
)

func FindIncludes(fileContent string) []string {
	includes := make([]string, 0)

	directiveMatcher := regexp.MustCompile("<!--@(.*)-->")

	for _, indices := range directiveMatcher.FindAllStringSubmatchIndex(fileContent, -1) {
		includes = append(includes, fileContent[indices[2]:indices[3]])
	}

	return includes
}

func createIncludeDirective(componentName string) string {
	return fmt.Sprintf("<!--@%s-->", componentName)
}
