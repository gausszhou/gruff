package main

import (
	"fmt"
	"log"

	"github.com/gausszhou/gruff"
)

func main() {
	md := `# Gruff Markdown Renderer

A lightweight, high-performance **markdown** renderer for the terminal.

## Features

- *Italic* and **bold** text
- ` + "`" + `Inline code` + "`" + ` support
- Unordered and ordered lists
- Headings (H1 through H6)

## Example

This paragraph has **bold**, *italic*, ` + "`" + `code` + "`" + `, and ***bold italic*** together.

1. First ordered item
2. Second ordered item
3. Third ordered item

---

The end.
`

	out, err := gruff.Render(md)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(out)
}
