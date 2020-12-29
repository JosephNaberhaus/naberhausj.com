package html

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"io/ioutil"
	"regexp"
	"strings"
)

// Matches all image tags with an SVG source
// Capture group #1 = All attributes of the node
// Capture group #2 = The path of the SVG file
var svgImageMatcher = regexp.MustCompile("<img (.*src=['\"](.*.svg)['\"].*)>")

// Matches the name and equals portion of an XML tag including the whitespace
// Capture group #1 = The name of the attribute
var attributeNameMatcher = regexp.MustCompile("(?s) *([^\\s\\\\]*) *= *")

// Matches SVG tags
// Capture group #1 = All attributes of the SVG tag
var svgMatcher = regexp.MustCompile("(?s)<svg (.*?)>.*</svg>")

func embedSVGs(content, htmlPath, root string) (newContent string, err error) {
	if len(content) == 0 {
		return "", nil
	}

	splitContent := svgImageMatcher.Split(content, -1)
	images := svgImageMatcher.FindAllStringSubmatch(content, -1)

	newContentBuilder := strings.Builder{}
	newContentBuilder.WriteString(splitContent[0])

	for i, image := range images {
		svgPath := file.ResolveHTMLPath(image[2], htmlPath, root)

		svgFile, err := ioutil.ReadFile(svgPath)
		if err != nil {
			return "", fmt.Errorf("couldn't load SVG file at path %s: %w", svgPath, err)
		}

		imageAttributes, err := parseAttributes(image[1])
		if err != nil {
			return "", fmt.Errorf("couldn't parse image attributes at path %s: %w", svgPath, err)
		}

		delete(imageAttributes, "src")

		svgMatch := svgMatcher.FindSubmatch(svgFile)
		if len(svgMatch) == 0 {
			return "", fmt.Errorf("no SVG tag found for file with path: %s", svgPath)
		}

		svgAttributeString := string(svgMatch[1])

		svgAttributes, err := parseAttributes(svgAttributeString)
		if err != nil {
			return "", fmt.Errorf("malformed SVG attributes for file with path %s: %w", svgPath, err)
		}

		svgAttributes.substituteAttributes(imageAttributes)

		substitutedSVG := strings.Replace(string(svgMatch[0]), string(svgMatch[1]), svgAttributes.String(), 1)

		newContentBuilder.WriteString(substitutedSVG)
		newContentBuilder.WriteString(splitContent[i+1])
	}

	return newContentBuilder.String(), nil
}

type attributes map[string]string

func parseAttributes(attributesString string) (attributes attributes, err error) {
	attributesString = strings.TrimSpace(attributesString)

	attributeNames := attributeNameMatcher.FindAllStringSubmatch(attributesString, -1)
	attributeValues := attributeNameMatcher.Split(attributesString, -1)

	if len(attributeValues) >= 1 && attributeValues[0] == "" {
		attributeValues = attributeValues[1:]
	}

	if len(attributeNames) != len(attributeValues) {
		return nil, fmt.Errorf("malformed XML attributes \"%s\"", attributesString)
	}

	attributes = make(map[string]string)

	for i := range attributeNames {
		attributes[attributeNames[i][1]] = attributeValues[i]
	}

	return attributes, nil
}

func (attr attributes) substituteAttributes(newAttributes attributes) {
	for name, value := range newAttributes {
		attr[name] = value
	}
}

func (attr attributes) String() string {
	sb := strings.Builder{}

	for name, value := range attr {
		sb.WriteString(name)
		sb.WriteString("=")
		sb.WriteString(value)
		sb.WriteRune(' ')
	}

	return sb.String()
}
