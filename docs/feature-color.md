# Color System

## Type

```go
type Color string
```

A `Color` holds either:
- **24-bit hex** — `"#RRGGBB"` or `"#RGB"` (shorthand), e.g. `"#141414"`
- **8-bit palette index** — `"0"`-`"255"`, e.g. `"236"` (dark grey)

## Hex Parsing

Hex colors are parsed by `lipgloss.Color(s)` from `charm.land/lipgloss/v2`, which returns a `color.Color` interface. Gruff calls `.RGBA()` to extract R/G/B values:

```go
func hexRGB(c Color) (r, g, b uint8) {
    cc := lipgloss.Color(string(c))
    rr, gg, bb, _ := cc.RGBA()
    return uint8(rr >> 8), uint8(gg >> 8), uint8(bb >> 8)
}
```

## ANSI Code Generation

| Input | Output | Format |
|---|---|---|
| `"#RRGGBB"` (hex) | `\x1b[38;2;R;G;Bm` (fg) / `\x1b[48;2;R;G;Bm` (bg) | 24-bit true color |
| `"0"`-`"7"` (4-bit) | `\x1b[3Nm` (fg) / `\x1b[4Nm` (bg) | ANSI 3/4-bit |
| `"8"`-`"255"` (8-bit) | `\x1b[38;5;Nm` (fg) / `\x1b[48;5;Nm` (bg) | ANSI 8-bit |

## Built-in Palette Constants

```go
cBlack="0", cMaroon="1", cGreen="2",  cOlive="3",
cNavy="4",  cPurple="5", cTeal="6",   cSilver="7",
cGrey="8",  cRed="9",    cLime="10",  cYellow="11",
cBlue="12", cFuchsia="13", cCyan="14", cWhite="15",
cDarkBG="236"
```

## Document Background

`Theme.Background` sets the entire document's background color (default `"#141414"` for dark theme, `""` for light theme).

- Emitted once at document start: `\x1b[48;2;20;20;20m`
- Restored after any style that sets a background (e.g., `Code` with `Bg: "236"`): `Style.end(bg)` emits `\x1b[48;2;20;20;20m` instead of `\x1b[49m`
- In table cells, the background is restored after content, before padding, preventing inline-code background bleed

## Style

```go
type Style struct {
    Fg        Color  // foreground color
    Bg        Color  // background color
    Bold      bool
    Italic    bool
    Underline bool
}

func (s Style) start() ansiCode   // emit ANSI codes to enter style
func (s Style) end(bg Color) ansiCode  // emit ANSI codes to exit style, restore doc bg
```

## Theme

```go
type Theme struct {
    Background              Color
    H1, H2, H3, H4, H5, H6 Style
    Strong                  Style
    Em                      Style
    Code, Link, LinkURL     Style
    Bullet, Numbered        Style
}
```

### Dark Theme (default)

| Element | Style |
|---|---|
| Background | `"#141414"` |
| H1 | Bold + Underline + White |
| H2 | Bold + Yellow |
| H3 | Bold + Green |
| H4 | Bold + Cyan |
| H5 | Bold + Grey |
| H6 | Grey |
| Code | Fg: White + Bg: 236 (dark grey) |
| Link | Underline + Cyan |

### Light Theme

| Element | Style |
|---|---|
| Background | (none) |
| H1 | Bold + Underline + Black |
| H2 | Bold + Navy |
| Code | Fg: Black + Bg: Silver |
| Link | Underline + Navy |
