package main

import (
	"flag"
	"fmt"
	"github.com/JosephNaberhaus/naberhausj.com/builder/file"
	"github.com/JosephNaberhaus/naberhausj.com/builder/handlers/css"
	"github.com/JosephNaberhaus/naberhausj.com/builder/handlers/html"
	"github.com/JosephNaberhaus/naberhausj.com/builder/handlers/image"
	"github.com/JosephNaberhaus/naberhausj.com/builder/handlers/resource"
	"github.com/JosephNaberhaus/naberhausj.com/builder/orchestrator"
	"os"
	"time"
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

	start := time.Now()
	err = builder.Build()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Build finished in %dms\n", time.Since(start).Milliseconds())
}
