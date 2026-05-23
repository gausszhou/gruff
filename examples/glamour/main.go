package main

import (
	"flag"
	"fmt"
	"log"

	"charm.land/glamour/v2"
)

var sampleMD = "" +
	"# Gruff vs Glamour\n\n" +
	"A comparison of **markdown** rendering between `gruff` and *glamour*.\n\n" +
	"## Text Formatting\n\n" +
	"- **Bold text**\n" +
	"- *Italic text*\n" +
	"- ***Bold italic***\n" +
	"- `Inline code`\n" +
	"- ~~Strikethrough~~\n\n" +
	"## Code Blocks\n\n" +
	"```go\n" +
	"package main\n\n" +
	"import \"fmt\"\n\n" +
	"func main() {\n" +
	"    fmt.Println(\"Hello\")\n" +
	"}\n" +
	"```\n\n" +
	"```\n" +
	"plain code block\n" +
	"```\n\n" +
	"## Tables\n\n" +
	"| Left | Center | Right |\n" +
	"|:-----|:------:|------:|\n" +
	"| Hello | ✅ | ¥3000 |\n" +
	"| **Bold** | *Italic* | `Code` |\n" +
	"| 🚀 | ⭐ | ❌ |\n\n" +
	"## Links\n\n" +
	"- [Gruff](https://github.com/gausszhou/gruff)\n" +
	"- **Bold link**: [**gruff**](https://github.com/gausszhou/gruff)\n\n" +
	"## Lists\n\n" +
	"- unordered item 1\n" +
	"- unordered item 2\n" +
	"  1. nested ordered\n" +
	"  2. nested ordered\n\n" +
	"1. ordered one\n" +
	"2. ordered two\n\n" +
	"---\n\n" +
	"The end.\n"

func main() {
	style := flag.String("s", "dark", "glamour style: dark, light, notty")
	wrap := flag.Int("wrap", 60, "word wrap width")
	flag.Parse()

	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(*style),
		glamour.WithWordWrap(*wrap),
	)
	if err != nil {
		log.Fatal(err)
	}

	out, err := r.Render(sampleMD)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(out)
}
