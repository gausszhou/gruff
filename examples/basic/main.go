package main

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/gausszhou/gruff"
)

//go:embed testdata/sample.md
var sampleMD string

func main() {
	out, err := gruff.Render(sampleMD)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(out)
}
