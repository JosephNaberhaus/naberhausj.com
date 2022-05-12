package main

import (
	"flag"
	"github.com/JosephNaberhaus/naberhausj.com/builder/orchestrator"
	"log"
)

var source = flag.String("source", "", "the directory that contains the source files")
var output = flag.String("output", "", "the directory to output the built files to")

func main() {
	panic("test")
	builder, err := orchestrator.CreateBuilder(*source, *output)
	if err != nil {
		log.Fatal(err)
	}

	err = builder.Build()
	if err != nil {
		log.Fatal(err)
	}
}

//package main
//
//import (
//	"flag"
//	"fmt"
//	"github.com/JosephNaberhaus/naberhausj.com/builder/component"
//	"github.com/JosephNaberhaus/naberhausj.com/builder/css"
//	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
//	"github.com/JosephNaberhaus/naberhausj.com/builder/html"
//	"github.com/JosephNaberhaus/naberhausj.com/builder/resource"
//	"log"
//	"os"
//	"path/filepath"
//)
//
//var fast = flag.Bool("fast", false, "perform a fast developer build")
//var root = flag.String("root", "", "the directory that contains the src and out directories")
//
//func main() {
//	flag.Parse()
//
//	wd, err := os.Getwd()
//	println(wd)
//	if err != nil {
//		log.Fatal(fmt.Errorf("failed to get working directory: %w", err))
//	}
//
//	log.Printf("Looking for definitions within %s", wd)
//	sourceDirectory := filepath.Join(wd, *root, "src")
//	definitions, err := component.FindDefinitions(sourceDirectory)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("%d component(s) found", len(definitions))
//
//	log.Println("Loading all html files")
//	htmlFiles, err := html.LoadHtmlFiles(sourceDirectory)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("%d html file(s) found", len(htmlFiles))
//
//	log.Println("Building document collection")
//	nameToPathMap, err := definitions.CreateNameToPathMap()
//	if err != nil {
//		log.Fatal(err)
//	}
//	documents, err := html.BuildDocumentCollection(htmlFiles, nameToPathMap)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	log.Println("Identifying non-component html files")
//	nonComponentHtmlFiles := make([]*html.File, 0)
//	for _, file := range htmlFiles {
//		if !definitions.ContainsPath(file.Path()) {
//			nonComponentHtmlFiles = append(nonComponentHtmlFiles, file)
//		}
//	}
//
//	outputDirectory := filepath.Join(wd, *root, "out")
//	outDirectoryExists, err := fileExists(outputDirectory)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if outDirectoryExists {
//		log.Println("Clearing output directory")
//		err = file.RemoveContents(outputDirectory)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//
//	log.Println("Substituting and writing non-component html files")
//	for _, file := range nonComponentHtmlFiles {
//		wd, err := filepath.Rel(sourceDirectory, filepath.Dir(file.Path()))
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		result, err := documents.Substitute(file.Path())
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		if !*fast {
//			result, err = html.SubstituteDate(file.Path(), result)
//			if err != nil {
//				log.Fatal(err)
//			}
//		}
//
//		result, err = html.SubstituteImages(result, sourceDirectory, wd, outputDirectory, *fast)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		err = html.Output(sourceDirectory, outputDirectory, file.Path(), result)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//
//	log.Println("Combining all CSS files")
//	cssContent, err := css.ConcatFiles(sourceDirectory)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	log.Println("Producing main CSS file")
//	err = css.Output(outputDirectory, cssContent)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	log.Println("Linking resources")
//	err = resource.MoveResources(sourceDirectory, outputDirectory)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//func fileExists(path string) (bool, error) {
//	_, err := os.Stat(path)
//	if err == nil {
//		return true, nil
//	}
//	if os.IsNotExist(err) {
//		return false, nil
//	}
//	return false, err
//}
