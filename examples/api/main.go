package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gausszhou/gruff/gruff"
	"golang.org/x/term"
)

func termWidth() int {
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		return w
	}
	return 0
}

func main() {
	source := "# Gruff API Demo\n\n" +
		"A lightweight, high-performance **markdown** renderer for the terminal.\n" +
		"Supports CJK (你好世界), emoji (✅ 🚀), math (∑π∞), and more.\n\n" +
		"## Text Formatting\n\n" +
		"- *Italic* and **bold** text\n" +
		"- `Inline code` with background\n" +
		"- ***Bold italic*** combined\n" +
		"- ~~Strikethrough~~ text\n" +
		"- [Link to GitHub](https://github.com/gausszhou/gruff)\n\n" +
		"## Lists\n\n" +
		"- Unordered item with **bold**\n" +
		"- Item with `code`\n" +
		"- Item with ✅ emoji\n" +
		"- Item with 中文\n\n" +
		"1. Ordered item with *italic*\n" +
		"2. Second item\n\n" +
		"## Task List\n\n" +
		"- [x] Completed task\n" +
		"- [ ] Pending task\n\n" +
		"## Table\n\n" +
		"| Left | Center | Right |\n" +
		"|:-----|:------:|------:|\n" +
		"| 你好 | ✅ | ¥3000 |\n" +
		"| Hello | ∑ π | ∞ |\n" +
		"| **Bold** | *Italic* | `Code` |\n\n" +
		"## Code Block\n\n" +
		"```\n" +
		"func hello() {\n" +
		"    fmt.Println(\"Hello, World!\")\n" +
		"}\n" +
		"```\n\n" +
		"> A wise blockquote.\n" +
		"> Multiple lines.\n\n" +
		"---\n\n" +
		"Use `Render` for strings, `RenderBytes` for `[]byte`, and `WithWordWrap` to control line width.\n"

	// Render with default dark theme
	out1, err := gruff.Render(source, gruff.WithWordWrap(termWidth()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== Default (dark theme) ===")
	fmt.Print(out1)

	// Render with light theme
	out2, err := gruff.Render(source, gruff.WithLight(), gruff.WithWordWrap(termWidth()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== Light theme ===")
	fmt.Print(out2)

	// RenderBytes for []byte input
	bytes := []byte("# Bytes Demo\n\nRendering from `[]byte` input.\n\n" +
		"- **Bold** and *italic*\n" +
		"- `inline code`\n")
	out3, err := gruff.RenderBytes(bytes, gruff.WithWordWrap(termWidth()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== RenderBytes ===")
	fmt.Print(string(out3))

	// Word wrap at 40 columns
	wrapSource := "# Word Wrap\n\n" +
		strings.Repeat("This is a long sentence that wraps. ", 8) + "\n\n" +
		"| Col1 | Col2 |\n" +
		"|:-----|:-----|\n" +
		"| 你好世界 | ✅ 🚀 |\n" +
		"| ∑π∞ | ★♥♦ |\n"
	out4, err := gruff.Render(wrapSource, gruff.WithWordWrap(40))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== Word wrap at 40 ===")
	fmt.Print(out4)
}
