package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gausszhou/gruff/gruff"
	"golang.org/x/term"
)

func main() {
	light := flag.Bool("light", false, "use light theme")
	wrap := flag.Int("wrap", 120, "word wrap width")
	mdFile := flag.String("md", "examples/codeblock/codeblock.md", "path to markdown file")
	flag.Parse()

	raw, err := os.ReadFile(*mdFile)
	if err != nil {
		log.Fatal(err)
	}

	var opts []gruff.Option
	if *light {
		opts = append(opts, gruff.WithLight())
	}
	if *wrap > 0 {
		opts = append(opts, gruff.WithWordWrap(*wrap))
	} else if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		opts = append(opts, gruff.WithWordWrap(w))
	}

	out, err := gruff.Render(string(raw), opts...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(out)
}
