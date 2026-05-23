# Architecture

## Overview

Gruff is a markdown-to-ANSI terminal renderer with zero runtime dependencies beyond goldmark (parser) and go-runewidth / lipgloss v2 (Unicode width, color parsing).

```
input (markdown string)
    │
    ▼
goldmark parser ──► AST
    │
    ▼
nodeRenderer ──► ANSI string (raw)
    │
    ▼
wrapText (document-level word wrap)
    │
    ▼
output (ANSI string)
```

## Core Files

| File | Responsibility |
|---|---|
| `gruff.go` | Public API (`Render`, `RenderBytes`, `Options`, options funcs) |
| `renderer.go` | goldmark AST walker → ANSI output (block/inline rendering) |
| `ansi.go` | SGR escape codes, `Style`/`Color`/`Theme` types, helpers |

## Key Design Decisions

1. **No `\x1b[0m` (full reset) in inline styles** — use specific undo codes (`\x1b[22m` noBold, `\x1b[39m` default fg, `\x1b[49m`/doc bg) to preserve outer style state during nesting.

2. **Two-pass table rendering** — Pass 1 collects all cell content and calculates column widths; Pass 2 renders with UTF-8 box-drawing separators. See `feature-table.md`.

3. **ANSI-aware word wrap** — `wrapText` (document) and `wrapCellLines` (cells) count only visible characters using `displayWidth` (go-runewidth), preserving leading spaces.

4. **Color via lipgloss v2** — Colors are strings (hex `"#RRGGBB"` or palette `"0"-"255"`). Hex parsing delegates to `lipgloss.Color(s)`. 24-bit true color output for hex, 8-bit for palette indices. See `feature-color.md`.

5. **Document background** — `Theme.Background` emits `\x1b[48;2;R;G;Bm` at document start; `Style.end(bg)` restores doc background instead of `\x1b[49m` when a style had a background color.

## Public API

```go
func Render(source string, opts ...Option) (string, error)
func RenderBytes(source []byte, opts ...Option) ([]byte, error)

type Options struct {
    Theme    Theme
    WordWrap int
}

func WithDark() Option          // use dark theme
func WithLight() Option         // use light theme
func WithWordWrap(n int) Option // set maximum line width (default 120)
```
