package image

import (
	"bytes"
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/command"
	"io/ioutil"
	"os"
	"os/exec"
)

type resizeFunc func(img Image, width, height int) (Image, error)

func resizeToPng(img Image, width, height int) (Image, error) {
	newSize := fmt.Sprintf("%dx%d!", width, height)
	cmd := exec.Command("convert", "-", "-resize", newSize, "png:-")
	cmd.Stdin = bytes.NewReader(img.data)

	output, err := command.RunWithErrorChecking(cmd)
	if err != nil {
		return Image{}, fmt.Errorf("error resizing to png: %w", err)
	}

	return Image{
		extension: ".png",
		data:      output,
	}, nil
}

func resizeToJpg(img Image, width, height int) (Image, error) {
	newSize := fmt.Sprintf("%dx%d!", width, height)
	cmd := exec.Command("convert", "-", "-resize", newSize, "-quality", "70", "jpg:-")
	cmd.Stdin = bytes.NewReader(img.data)

	output, err := command.RunWithErrorChecking(cmd)
	if err != nil {
		return Image{}, fmt.Errorf("error resizing to jpg: %w", err)
	}

	return Image{
		extension: ".jpg",
		data:      output,
	}, nil
}

func resizeToWebp(img Image, width, height int) (Image, error) {
	const tempFileName = "temp.webp"

	newSize := fmt.Sprintf("%dx%d!", width, height)
	cmd := exec.Command("convert", "-", "-resize", newSize, "-quality", "70", tempFileName)
	cmd.Stdin = bytes.NewReader(img.data)

	_, err := command.RunWithErrorChecking(cmd)
	if err != nil {
		return Image{}, fmt.Errorf("error resizing to webp: %w", err)
	}

	output, err := ioutil.ReadFile(tempFileName)
	if err != nil {
		return Image{}, fmt.Errorf("error reading temporary webp file: %w", err)
	}

	err = os.Remove(tempFileName)
	if err != nil {
		return Image{}, fmt.Errorf("error removing temporary webp file: %w", err)
	}

	return Image{
		extension: ".webp",
		data:      output,
	}, nil
}

var resizeFunctions = []resizeFunc{
	resizeToPng,
	resizeToJpg,
	resizeToWebp,
}

func (i Image) resize(width, height int) (Image, error) {
	best := i

	for _, resizeFunc := range resizeFunctions {
		result, err := resizeFunc(i, width, height)
		if err != nil {
			return Image{}, fmt.Errorf("error resizing image: %w", err)
		}

		if len(result.data) < len(best.data) {
			best = result
		}
	}

	return best, nil
}
