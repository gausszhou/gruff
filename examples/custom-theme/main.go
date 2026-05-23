package main

import (
	"fmt"
	"log"

	"github.com/gausszhou/gruff"
)

func customTheme() gruff.Option {
	return func(o *gruff.Options) {
		o.Theme.H1 = gruff.Style{Fg: gruff.Color(196), Bold: true}
		o.Theme.H2 = gruff.Style{Fg: gruff.Color(208), Bold: true}
		o.Theme.H3 = gruff.Style{Fg: gruff.Color(220), Bold: true}
		o.Theme.Strong = gruff.Style{Bold: true, Fg: gruff.Color(51)}
		o.Theme.Em = gruff.Style{Italic: true, Fg: gruff.Color(213)}
		o.Theme.Code = gruff.Style{Bg: gruff.Color(235), Fg: gruff.Color(120)}
		o.Theme.Link = gruff.Style{Underline: true, Fg: gruff.Color(39)}
		o.Theme.Bullet = gruff.Style{Fg: gruff.Color(202)}
		o.Theme.Numbered = gruff.Style{Fg: gruff.Color(202)}
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
		"| No external deps | ✅ |\n"

	out, err := gruff.Render(md, customTheme())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(out)
}
