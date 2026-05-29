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
	"# Code Block Examples\n\n" +
	"A dedicated example showcasing code block rendering.\n\n" +

	"## Plain Fenced Code\n\n" +
	"A fenced code block without language annotation:\n\n" +
	"```\n" +
	"line1\n" +
	"line2\n" +
	"line3\n" +
	"```\n\n" +

	"## Go Language\n\n" +
	"Fenced code block with Go language tag:\n\n" +
	"```go\n" +
	"package main\n\n" +
	"import \"fmt\"\n\n" +
	"func main() {\n" +
	"    fmt.Println(\"Hello, World!\")\n" +
	"}\n" +
	"```\n\n" +

	"## Python Language\n\n" +
	"A fenced code block with Python language tag:\n\n" +
	"```python\n" +
	"def hello():\n" +
	"    print(\"Hello, World!\")\n\n" +
	"hello()\n" +
	"```\n\n" +

	"## Indented Code Block\n\n" +
	"An indented code block (no language tag):\n\n" +
	"    line one\n" +
	"    line two\n" +
	"    line three\n\n" +

	"## Mixed Content\n\n" +
	"A paragraph with **bold**, *italic*, and `inline code` followed by more code:\n\n" +
	"```\n" +
	"# This is a comment\n" +
	"result = 42\n" +
	"print(f\"The answer is {result}\")\n" +
	"```\n\n" +

	"Visit [Gruff on GitHub](https://github.com/gausszhou/gruff) for more information.\n"

func main() {
	light := flag.Bool("light", false, "use light theme")
	wrap := flag.Int("wrap", 120, "word wrap width")
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
