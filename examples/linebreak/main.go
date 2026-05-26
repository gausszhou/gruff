package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gausszhou/gruff"
	"golang.org/x/term"
)

var sampleMD = "" +
	"# Line Break Examples\n\n" +

	"## Soft Line Break (two trailing spaces)\n\n" +
	"First line with two spaces at end  \n" +
	"Second line\n\n" +

	"## Hard Line Break (backslash)\n\n" +
	"First line with backslash\\\n" +
	"Second line\n\n" +

	"## Normal Paragraph (no break)\n\n" +
	"A single paragraph with text that " +
	"continues on the same line because " +
	"there is no trailing space or backslash.\n\n" +

	"## Line Break in List\n\n" +
	"- Item with soft break  \n  continuation\n" +
	"- Another item\n\n" +

	"## Paragraph with multiple breaks\n\n" +
	"Line one  \n" +
	"Line two  \n" +
	"Line three\n\n" +

	"## Very Long Single Line (word wrap test)\n\n" +
	"Lorem ipsum dolor sit amet consectetur adipiscing elit sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident sunt in culpa qui officia deserunt mollit anim id est laborum. Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium totam rem aperiam eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.\n\n" +

	"The end.\n"

func main() {
	light := flag.Bool("light", false, "use light theme")
	wrap := flag.Int("wrap", 80, "word wrap width")
	flag.Parse()

	var opts []gruff.Option
	if *light {
		opts = append(opts, gruff.WithLight())
	}
	if *wrap > 0 {
		opts = append(opts, gruff.WithWordWrap(*wrap))
	} else if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		opts = append(opts, gruff.WithWordWrap(w))
	}

	out, err := gruff.Render(sampleMD, opts...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(out)
}
