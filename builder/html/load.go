package html

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"io/ioutil"
)

type File struct {
	path    string
	content string
}

func (f *File) Path() string {
	return f.path
}

func (f *File) Content() string {
	return f.content
}

func LoadHtmlFiles(root string) (htmlFiles []*File, err error) {
	htmlFiles = make([]*File, 0)

	paths, err := file.FindFilesWithSuffix(root, ".html")
	if err != nil {
		return nil, fmt.Errorf("error finding html files; %w", err)
	}

	for _, path := range paths {
		fileContent, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error loading html file: %w", err)
		}

		htmlFile := new(File)
		htmlFile.path = path
		htmlFile.content = string(fileContent)
		htmlFiles = append(htmlFiles, htmlFile)
	}

	return htmlFiles, nil
}
