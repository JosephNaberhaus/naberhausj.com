package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func RemoveContents(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		err = os.RemoveAll(filepath.Join(dir, fileInfo.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
