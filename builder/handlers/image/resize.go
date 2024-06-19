package image

import (
	"bytes"
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/command"
	"io/ioutil"
	"os"
	"os/exec"
)

func (i Image) resize(width, height int) (Image, error) {
	const tempFileName = "temp.webp"

	newSize := fmt.Sprintf("%dx%d!", width, height)
	cmd := exec.Command("convert", "-", "-resize", newSize, "-quality", "70", tempFileName)
	cmd.Stdin = bytes.NewReader(i.data)

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
