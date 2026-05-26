package main

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/gausszhou/gruff"
)

func main() {
	var md strings.Builder
	for r := rune(0x20); r <= 0x10FFFF; r++ {
		if r >= 0x7F && r <= 0x9F {
			continue
		}
		if r >= 0xD800 && r <= 0xDFFF {
			continue
		}
		if !unicode.IsPrint(r) {
			continue
		}
		md.WriteRune(r)
		md.WriteRune(' ')
	}
	out, err := gruff.Render(md.String(), gruff.WithWordWrap(80))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(out)
}
