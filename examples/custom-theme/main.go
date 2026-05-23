package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gausszhou/gruff"
	"golang.org/x/term"
)

func customTheme() gruff.Option {
	return func(o *gruff.Options) {
		o.Theme.H1 = gruff.Style{Fg: "196", Bold: true}
		o.Theme.H2 = gruff.Style{Fg: "208", Bold: true}
		o.Theme.H3 = gruff.Style{Fg: "220", Bold: true}
		o.Theme.Strong = gruff.Style{Bold: true, Fg: "51"}
		o.Theme.Em = gruff.Style{Italic: true, Fg: "213"}
		o.Theme.Code = gruff.Style{Bg: "235", Fg: "120"}
		o.Theme.Link = gruff.Style{Underline: true, Fg: "39"}
		o.Theme.Bullet = gruff.Style{Fg: "202"}
		o.Theme.Numbered = gruff.Style{Fg: "202"}
	}
}

func main() {
	md := "# Custom Theme Demo\n\n" +
		"This example shows a **bold statement** with *emphasis* and `inline code`.\n\n" +
		"## Lists\n\n" +
		"- Item with **highlighted bold**\n" +
		"- Item with *italic emphasis*\n" +
		"- Item with `code style`\n\n" +
		"## Links\n\n" +
		"Visit [gruff](https://github.com/gausszhou/gruff) for more information.\n\n" +
		"## Table\n\n" +
		"| Feature | Status |\n" +
		"|---------|--------|\n" +
		"| Custom colors | ✅ |\n" +
		"| No external deps | ✅ |\n\n" +
		"## Code Block\n\n" +
		"```go\n" +
		"package main\n\n" +
		"import \"fmt\"\n\n" +
		"func main() {\n" +
		"    fmt.Println(\"Custom theme code block\")\n" +
		"}\n" +
		"```\n"

	var opts []gruff.Option
	opts = append(opts, customTheme())
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		opts = append(opts, gruff.WithWordWrap(w))
	}

	out, err := gruff.Render(md, opts...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(out)
}
