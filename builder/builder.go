package main

import (
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/component"
	"github.com/JosephNaberhaus/naberhausj.com/builder/css"
	"github.com/JosephNaberhaus/naberhausj.com/builder/html"
	"log"
	"os"
	"path/filepath"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get working directory: %w", err))
	}

	log.Printf("Looking for definitions within %s", wd)
	sourceDirectory := filepath.Join(wd, "src")
	definitions, err := component.FindDefinitions(sourceDirectory)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d component(s) found", len(definitions))

	log.Println("Loading all html files")
	htmlFiles, err := html.LoadHtmlFiles(sourceDirectory)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d html file(s) found", len(htmlFiles))

	log.Println("Building dependency graph")
	nameToPathMap, err := definitions.CreateNameToPathMap()
	if err != nil {
		log.Fatal(err)
	}
	graph, err := html.BuildDependencyGraph(htmlFiles, nameToPathMap)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Identifying non-component html files")
	nonComponentHtmlFiles := make([]*html.File, 0)
	for _, file := range htmlFiles {
		if !definitions.ContainsPath(file.Path()) {
			nonComponentHtmlFiles = append(nonComponentHtmlFiles, file)
		}
	}

	log.Println("Substituting non-component html files")
	for _, file := range nonComponentHtmlFiles {
		err = graph.SubstitutePath(file.Path())
		if err != nil {
			log.Fatal(err)
		}
	}

	outputDirectory := filepath.Join(wd, "out")
	outDirectoryExists, err := fileExists(outputDirectory)
	if err != nil {
		log.Fatal(err)
	}

	if outDirectoryExists {
		log.Println("Clearing output directory")
		err = os.RemoveAll(outputDirectory)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Outputting result")
	err = html.Output(sourceDirectory, outputDirectory, nonComponentHtmlFiles)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Combining all CSS files")
	cssContent, err := css.ConcatFiles(sourceDirectory)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Producing main CSS file")
	err = css.Output(outputDirectory, cssContent)
	if err != nil {
		log.Fatal(err)
	}
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
