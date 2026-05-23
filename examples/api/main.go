package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gausszhou/gruff"
	"golang.org/x/term"
)

func termWidth() int {
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		return w
	}
	return 0
}

func main() {
	source := "# API Demo\n\n" +
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
	bytes := []byte("# Bytes Demo\n\nRendering from `[]byte` input.\n")
	out3, err := gruff.RenderBytes(bytes, gruff.WithWordWrap(termWidth()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== RenderBytes ===")
	fmt.Print(string(out3))

	// Word wrap at 40 columns
	longText := "# Word Wrap\n\n" +
		strings.Repeat("This is a long sentence that wraps. ", 8) + "\n"
	out4, err := gruff.Render(longText, gruff.WithWordWrap(40))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== Word wrap at 40 ===")
	fmt.Print(out4)
}
