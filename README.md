# gruff

A lightweight, high-performance [Go](https://go.dev) library for rendering Markdown to ANSI-colored terminal output.

## Features

- **Headings** (H1–H6) with distinct styles per level
- **Bold**, *italic*, and ***bold italic***
- `Inline code` and fenced/indented **code blocks** — green text, language tag in gray
- [Links](https://github.com/gausszhou/gruff) — underlined + blue text with gray URL suffix
- Unordered (`-`, `*`) and ordered (`1.`) lists
- GFM tables with **UTF-8 box‑drawing borders**, alignment (left/center/right), inline formatting inside cells, and **automatic text wrapping** at a max column width
- Strikethrough text (`~~strikethrough~~`)
- Task lists (`- [x] done`, `- [ ] todo`)
- Thematic break (`---`)
- Dark and light themes
- ANSI-aware word wrap (`WithWordWrap`)

## Installation

```bash
go get github.com/gausszhou/gruff
```

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/gausszhou/gruff"
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

Benchmarked against `testdata/benchmark.md` (2.4 KB) repeated 100× (~240 KB input).
Glamour is tested in both minimal mode (chroma disabled, word wrap off, table wrap off,
inline table links on) and standard mode.

| Metric         | gruff      | glamour (minimal) | glamour (standard) | Improvement (vs minimal) |
|----------------|------------|-------------------|--------------------|--------------------------|
| Time/op        | ~25.7 ms   | ~265 ms           | ~1.19 s            | **~10× / ~46×**          |
| Memory/op      | ~14.0 MB   | ~52.9 MB          | ~222 MB            | **~3.8× / ~16×**         |
| Allocations/op | ~209,556   | ~4,505,093        | ~20,348,392        | **~21× / ~97×**          |

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

```bash
go run examples/basic/main.go
go run examples/table/main.go
go run examples/codeblock/main.go
go run examples/custom-theme/main.go
go run examples/api/main.go
go run examples/compare-benchmark/main.go
go run examples/compare-glamour/main.go
go run examples/compare-theme/main.go
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
        o.Theme.Link = gruff.Style{Underline: true, Fg: gruff.Color("#5c9cf5")}
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

**Test/Benchmark only:** `github.com/charmbracelet/glamour` (not included in production builds)

## License

MIT
