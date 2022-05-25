package file

import "os"

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	panic("unexpected error when checking if file exists")
}
