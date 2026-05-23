package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gausszhou/gruff"
)

var sampleMD = `# Gruff Markdown Renderer

A lightweight, high-performance **markdown** renderer for the terminal.

## Features

- *Italic* and **bold** text
- Inline code support
- Unordered and ordered lists
- Headings (H1 through H6)

## Formatting Showcase

This paragraph demonstrates **bold**, *italic*, inline code, and ***bold italic*** styles.

1. First ordered item with **bold text**
2. Second ordered item with *italic text*
3. Third ordered item with inline code

## Table

| Feature     | Status | Priority | Notes                      |
|:------------|:------:|:--------:|:---------------------------|
| Headings    | ✅     | High     | H1 through H6              |
| Bold        | ✅     | High     | Use **double asterisks**   |
| Italic      | ✅     | High     | Use *single asterisks*     |
| Inline Code | ✅     | High     | Use backticks              |
| Lists       | ✅     | High     | Ordered and unordered      |
| Links       | ✅     | Medium   | [text](url) with underline |

Visit [Gruff on GitHub](https://github.com/gausszhou/gruff) for more information.
`

func main() {
	light := flag.Bool("light", false, "use light theme")
	wrap := flag.Int("wrap", 0, "word wrap width (0 = no wrap)")
	flag.Parse()

	var opts []gruff.Option
	if *light {
		opts = append(opts, gruff.WithLight())
	}
	if *wrap > 0 {
		opts = append(opts, gruff.WithWordWrap(*wrap))
	}

	out, err := gruff.Render(sampleMD, opts...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(out)
}
