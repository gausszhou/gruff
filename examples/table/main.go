package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gausszhou/gruff"
	"golang.org/x/term"
)

func main() {
	wrap := flag.Int("wrap", 0, "word wrap width (0 = auto-detect)")
	flag.Parse()

	b, err := os.ReadFile("examples/table/table.md")
	if err != nil {
		log.Fatal(err)
	}
	md := string(b)

	var opts []gruff.Option
	if *wrap > 0 {
		opts = append(opts, gruff.WithWordWrap(*wrap))
	} else if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		opts = append(opts, gruff.WithWordWrap(w))
	}

	out, err := gruff.Render(md, opts...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(out)
}
