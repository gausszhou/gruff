# Gruff Markdown Renderer

A lightweight, high-performance **markdown** renderer for the terminal.

## Features

- *Italic* and **bold** text
- `Inline code` support
- Unordered and ordered lists
- Headings (H1 through H6)

## Usage Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/gausszhou/gruff"
)

func main() {
    out, err := gruff.Render("# Hello World")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Print(out)
}
```

## Formatting Showcase

This paragraph demonstrates **bold**, *italic*, `inline code`, and ***bold italic*** styles working together seamlessly.

1. First ordered item with **bold text**
2. Second ordered item with *italic text*
3. Third ordered item with `inline code`

### Nested Emphasis

Here we have **bold with *italic inside*** and *italic with **bold inside***.

## Code Elements

Inline code like `var x = 42` should stand out with a distinct background color.

## Lists with Mixed Content

- Item with **bold**
- Item with *italic*
- Item with `code`
- Item with ***both***

## Thematic Break

---

Above is a horizontal rule.

## Final Section

A plain paragraph to close the document with some **formatting** to make sure everything works correctly in the *final* output.
