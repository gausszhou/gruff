package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gausszhou/gruff"
)

func main() {
	source, err := os.ReadFile("../../testdata/sample.md")
	if err != nil {
		log.Fatal(err)
	}

	out, err := gruff.Render(string(source))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(out)
}
