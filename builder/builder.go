package main

import (
	"flag"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"github.com/JosephNaberhaus/naberhausj.com/builder/handlers/css"
	"github.com/JosephNaberhaus/naberhausj.com/builder/handlers/html"
	"github.com/JosephNaberhaus/naberhausj.com/builder/handlers/image"
	"github.com/JosephNaberhaus/naberhausj.com/builder/handlers/resource"
	"github.com/JosephNaberhaus/naberhausj.com/builder/orchestrator"
	"os"
)

var nocache = flag.Bool("nocache", false, "whether to do a full rebuild")
var src = flag.String("src", "", "the source directory")
var out = flag.String("out", "", "the output directory")

func main() {
	flag.Parse()

	if *nocache && file.Exists(*out) {
		err := os.RemoveAll(*out)
		if err != nil {
			panic(err)
		}
	}

	builder, err := orchestrator.CreateBuilder(*src, *out)
	if err != nil {
		panic(err)
	}

	builder.AddHandler(html.CreateHandler(builder))
	builder.AddHandler(css.CreateHandler(builder))
	builder.AddHandler(image.CreateHandler(builder))
	builder.AddHandler(resource.CreateHandler(builder))

	err = builder.Build()
	if err != nil {
		panic(err)
	}
}
