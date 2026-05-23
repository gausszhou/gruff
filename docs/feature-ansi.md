# ANSI Code Generation

## SGR Constants

| Constant | Code | Effect |
|---|---|---|
| `ansiReset` | `\x1b[0m` | Full reset (used only as fallback in `Style.end()`) |
| `ansiBold` | `\x1b[1m` | Bold / bright |
| `ansiItalic` | `\x1b[3m` | Italic |
| `ansiUnderline` | `\x1b[4m` | Underline |
| `ansiNoBold` | `\x1b[22m` | Remove bold |
| `ansiNoItalic` | `\x1b[23m` | Remove italic |
| (inline) | `\x1b[24m` | Remove underline |
| (inline) | `\x1b[39m` | Default foreground |
| (inline) | `\x1b[49m` | Default background (or doc background) |

## Style Protocol

Every styled element calls `Style.start()` before its content and `Style.end(bg)` after.

### `start()`

1. Bold → `\x1b[1m`
2. Italic → `\x1b[3m`
3. Underline → `\x1b[4m`
4. Foreground color → `\x1b[38;...m`
5. Background color → `\x1b[48;...m`

### `end(bg Color)`

1. Italic → `\x1b[23m`
2. Bold → `\x1b[22m`
3. Underline → `\x1b[24m`
4. Foreground → `\x1b[39m` (always reset to default fg)
5. Background → `\x1b[48;2;R;G;Bm` (restore doc bg) or `\x1b[49m` (terminal default)

**Why specific undo codes instead of `\x1b[0m`?** Consider nested styles:

```
\x1b[1mbold \x1b[3mitalic and bold\x1b[23m just bold\x1b[22m plain
```

Using `\x1b[0m` at the inner close would also reset the outer bold state. By using specific undo codes, each style layer is independently managed.

## Nesting Examples

```
**bold**          → \x1b[1mbold\x1b[22m
*italic*          → \x1b[3mitalic\x1b[23m
***both***        → \x1b[3m\x1b[1mboth\x1b[22m\x1b[23m
**bold `code`**   → \x1b[1mbold \x1b[48;5;236;38;5;15mcode\x1b[39m\x1b[48;2;20;20;20m\x1b[22m
```

## Table Cell Background Isolation

After rendering each cell's content, the renderer emits:

```
\x1b[39m           ← reset foreground
\x1b[48;2;R;G;Bm  ← restore document background
```

This prevents inline code backgrounds from bleeding into the cell's padding area or the separator `│`.
