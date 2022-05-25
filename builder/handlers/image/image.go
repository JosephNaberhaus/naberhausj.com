package image

import (
	"bytes"
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/command"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

type Image struct {
	extension string
	data      []byte
}

func LoadImage(path string) (Image, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Image{}, err
	}

	return Image{
		extension: filepath.Ext(path),
		data:      data,
	}, nil
}

func (i Image) Dimensions() (int, int, error) {
	cmd := exec.Command("convert", "-", "-format", "%[w]x%[h]", "info:")
	cmd.Stdin = bytes.NewReader(i.data)

	output, err := command.RunWithErrorChecking(cmd)
	if err != nil {
		return 0, 0, err
	}

	var width, height int
	_, err = fmt.Sscanf(string(output), "%dx%d", &width, &height)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing dimensions: %w", err)
	}

	return width, height, nil
}
