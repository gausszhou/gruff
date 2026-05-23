package main

import (
	"fmt"
	"log"

	"github.com/gausszhou/gruff"
)

func main() {
	md := "# Table Rendering Demo\n\n" +
		"## Basic Table\n\n" +
		"| Left | Center | Right |\n" +
		"|:-----|:------:|------:|\n" +
		"| a    | b      | c     |\n" +
		"| short | centered | right |\n\n" +
		"## Table with Inline Formatting\n\n" +
		"| Feature | Code | Notes |\n" +
		"|---------|------|-------|\n" +
		"| **Bold** | `**text**` | Double asterisks |\n" +
		"| *Italic* | `*text*` | Single asterisks |\n" +
		"| `Code` | `backticks` | Inline code |\n" +
		"| ***Both*** | `***text***` | Combine both |\n\n" +
		"## Table with Word Wrap\n\n" +
		"| Description | Details |\n" +
		"|-------------|---------|\n" +
		"| Long Content | This is a very long sentence that demonstrates how the table cell word wrapping works when the content exceeds the maximum column width of 40 characters. |\n" +
		"| Alignment Matters | Left aligned text wraps naturally. The column width is capped so long text wraps to multiple lines automatically. |\n\n" +
		"## Benchmarks\n\n" +
		"| Metric | gruff | glamour | Improvement |\n" +
		"|:-------|:-----:|:-------:|:-----------:|\n" +
		"| Time per op | ~429 µs | ~3,152 µs | **~7×** |\n" +
		"| Memory per op | ~426 KB | ~1,897 KB | **~4×** |\n" +
		"| Allocations | ~5,341 | ~23,357 | **~4×** |\n"

	out, err := gruff.Render(md)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(out)
}
