package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gausszhou/gruff"
	"golang.org/x/term"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: gruff <file.md>")
	}

	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	opts := []gruff.Option{gruff.WithWordWrap(w)}
	if err != nil || w <= 0 {
		opts = []gruff.Option{gruff.WithWordWrap(80)}
	}

	out, err := gruff.Render(string(b), opts...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(out)
}
