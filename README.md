# gruff

A lightweight, high-performance [Go](https://go.dev) library for rendering Markdown to ANSI-colored terminal output.

## Features

- **Headings** (H1–H6) with distinct styles per level
- **Bold**, *italic*, and ***bold italic***
- `Inline code` and fenced/indented **code blocks** — green text, language tag in gray
- [Links](https://github.com/gausszhou/gruff) — bold link text + underlined URL with OSC 8 terminal hyperlink support
- Unordered (`-`, `*`) and ordered (`1.`) lists
- GFM tables with **UTF-8 box‑drawing borders**, alignment (left/center/right), inline formatting inside cells, and **automatic text wrapping** at a max column width
- Strikethrough text (`~~strikethrough~~`)
- Task lists (`- [x] done`, `- [ ] todo`)
- Thematic break (`---`)
- Dark and light themes
- ANSI-aware word wrap (`WithWordWrap`)

## Installation

```bash
go get github.com/gausszhou/gruff/gruff
```

### CLI

```bash
go install github.com/gausszhou/gruff@latest
gruff README.md              # render a file
gruff render README.md       # explicit subcommand
gruff render -w 40 README.md # custom wrap width
gruff render -l README.md    # light theme
```

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/gausszhou/gruff/gruff"
)

func main() {
    md := `# Hello, gruff!

This is **bold**, *italic*, and \`inline code\`.

| Feature | Status |
|---------|--------|
| Tables  | ✅     |
| Speed   | 🚀     |
`

    out, err := gruff.Render(md)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Print(out)
}
```

### Options

| Function | Description |
|----------|-------------|
| `WithDark()` | Dark background theme (default) |
| `WithLight()` | Light background theme |
| `WithWordWrap(n)` | Wrap output at `n` columns |

```go
out, err := gruff.Render(source, gruff.WithLight(), gruff.WithWordWrap(80))
```

## Performance

Benchmarked against `testdata/benchmark.md` (~5.6 KB) repeated 100× (~560 KB input),
on Intel Core Ultra 7 255H.

| Metric         | gruff ¹     | glamour (minimal) ² | glamour (standard) ³ | Improvement (vs minimal / vs standard) |
|----------------|-------------|---------------------|----------------------|----------------------------------------|
| Time/op        | **~98 ms**  | ~435 ms             | ~1.39 s              | **~4.4× / ~14×**                       |
| Memory/op      | **~69 MB**  | ~137 MB             | ~441 MB              | **~2.0× / ~6.4×**                      |
| Allocations/op | **~459,000**| ~10,100,000         | ~39,300,000          | **~22× / ~86×**                        |

¹ gruff: `WithDark()` (no background), `WithWordWrap(120)`.
² glamour minimal: `Chroma = nil`, `CleanInput`, word wrap off, table wrap off, inline table links on.
³ glamour standard: `WithStandardStyle("dark")`, word wrap at 120 cols.

See [`docs/why-gruff-faster.md`](docs/why-gruff-faster.md) for a detailed analysis of the
performance gap.

Run benchmarks locally:

```bash
go test -bench=. -benchmem ./benchmark/
```

## Examples

Ready-to-run examples are in the [`examples/`](examples/) directory:

| Example | Description |
|---------|-------------|
| [`basic`](examples/basic/) | Render markdown with CLI flags (`--light`, `--wrap`) |
| [`table`](examples/table/) | Table-specific demo showing alignment and word wrap |
| [`codeblock`](examples/codeblock/) | Code block rendering with language tags |
| [`custom-theme`](examples/custom-theme/) | Custom ANSI color and style customization |
| [`api`](examples/api/) | `Render`, `RenderBytes`, and `WithWordWrap` usage |
| [`compare-benchmark`](examples/compare-benchmark/) | Side-by-side benchmark markdown rendered with gruff vs glamour |
| [`compare-glamour`](examples/compare-glamour/) | Side-by-side glamour standard vs minimal |
| [`compare-theme`](examples/compare-theme/) | Side-by-side gruff dark vs light theme |
| [`compare-simple`](examples/compare-simple/) | Glamour minimal vs standard without viewport/bubbletea |
| [`viewport-gruff`](examples/viewport-gruff/) | Gruff output in a bubbletea viewport with scrollbar |
| [`viewport-glamour`](examples/viewport-glamour/) | Glamour output in a bubbletea viewport with scrollbar |

```bash
go run examples/basic/main.go
go run examples/table/main.go
go run examples/codeblock/main.go
go run examples/custom-theme/main.go
go run examples/api/main.go
go run examples/compare-benchmark/main.go
go run examples/compare-glamour/main.go
go run examples/compare-theme/main.go
go run examples/compare-simple/main.go
go run examples/viewport-gruff/main.go
go run examples/viewport-glamour/main.go
```

## Theme Customization

Use the exported `Theme`, `Style`, and `Color` types to customize colors and styles:

```go
import "github.com/gausszhou/gruff"

customTheme := func() gruff.Option {
    return func(o *gruff.Options) {
        o.Theme.H1 = gruff.Style{Fg: gruff.Color(196), Bold: true}       // red
        o.Theme.Strong = gruff.Style{Bold: true, Fg: gruff.Color(51)}    // cyan
        o.Theme.Code = gruff.Style{Fg: gruff.Color("#50865a")}            // green
        o.Theme.Link = gruff.Style{Bold: true, Fg: gruff.Color("#5c9cf5")}
        o.Theme.LinkURL = gruff.Style{Underline: true, Fg: gruff.Color("#5c9cf5")}
    }
}

out, _ := gruff.Render(md, customTheme())
```

## How It Works

Parsing is handled by [`goldmark`](https://github.com/yuin/goldmark). The AST is walked by a recursive type switch in `renderer.go` that emits SGR ANSI codes directly — no intermediate DOM, no CSS, no HTML. Each inline style uses specific undo codes (e.g. `\x1b[22m` for no-bold, `\x1b[39m` for default foreground) instead of `\x1b[0m`, so nested formatting is preserved correctly.

Table rendering uses a two-pass approach: collect all cell content and calculate column widths, then emit UTF-8 box‑drawing borders and padded cell content. Columns cap at a maximum width with automatic word wrapping.

Code blocks are rendered line-by-line with language tags shown in gray, content in green, and each line padded to the document width for a clean full-width appearance.

## Dependencies

**Runtime:** `github.com/yuin/goldmark`, `github.com/mattn/go-runewidth`

**Examples only:** `charm.land/bubbles/v2`, `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`, `charm.land/glamour/v2` (not included in library builds)

## License

MIT
