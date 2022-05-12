package image

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type HtmlImage struct {
	Class  string `json:"class"`
	Source string `json:"src"`
	Alt    string `json:"alt"`
	Width  int    `json:"width"`
}

func ParseHtmlImage(value string) (HtmlImage, error) {
	var h HtmlImage

	err := json.Unmarshal([]byte(value), &h)
	if err != nil {
		return HtmlImage{}, fmt.Errorf("error parsing image: %w", err)
	}

	return h, nil
}

func imageWidths() []int {
	breakpoints := make([]int, 0, 10)
	for i := 10; i > 0; i-- {
		breakpoints = append(breakpoints, i*200)
	}
	return breakpoints
}

const fallbackWidth = 800

func (h HtmlImage) ImageElement(srcDir, wd, outputDir string, fast bool) (string, error) {
	image, err := LoadImage(filepath.Join(srcDir, wd, h.Source))
	if err != nil {
		return "", fmt.Errorf("error loading image: %w", err)
	}

	inputWidth, inputHeight, err := image.Dimensions()
	if err != nil {
		return "", err
	}

	displayWidth := h.Width
	displayHeight := (displayWidth * inputHeight) / inputWidth

	createImageWidth := func(outputWidth int) (string, error) {
		outputHeight := (outputWidth * inputHeight) / inputWidth

		resized, err := image.Resize(outputWidth, outputHeight)
		if err != nil {
			return "", err
		}

		sourceWithExtension := strings.TrimSuffix(h.Source, filepath.Ext(h.Source))
		resizedSource := sourceWithExtension + strconv.Itoa(outputWidth) + resized.extension
		outputPath := filepath.Join(outputDir, wd, resizedSource)

		err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
		if err != nil {
			return "", err
		}

		err = ioutil.WriteFile(outputPath, resized.data, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("error writing file: %w", err)
		}

		return resizedSource, nil
	}

	srcSet := strings.Builder{}
	if !fast {
		for _, width := range imageWidths() {
			resizedSource, err := createImageWidth(width)
			if err != nil {
				return "", err
			}

			srcSet.WriteString(fmt.Sprintf("%s %dw,", resizedSource, width))
		}
	}

	var fallbackSource string
	if fast {
		fallbackSource = h.Source

		outputPath := filepath.Join(outputDir, wd, h.Source)
		err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
		if err != nil {
			return "", err
		}

		err = ioutil.WriteFile(outputPath, image.data, os.ModePerm)
		if err != nil {
			return "", err
		}
	} else {
		fallbackSource, err = createImageWidth(fallbackWidth)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf(
		"<img class=\"%s\" width=\"%d\" height=\"%d\" src=\"%s\" srcset=\"%s\" sizes=\"%dpx\" alt=\"%s\">",
		h.Class,
		displayWidth,
		displayHeight,
		fallbackSource,
		srcSet.String(),
		displayWidth,
		h.Alt,
	), nil
}
