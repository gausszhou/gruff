# Table Rendering

Tables are rendered via a two-pass algorithm over the goldmark GFM table AST.

## AST Structure

```
Table
├── TableHeader
│   ├── TableCell (alignment: cell.Alignment)
│   ├── TableCell
│   └── ...
├── TableRow
│   ├── TableCell
│   ├── TableCell
│   └── ...
└── ...
```

**Note:** `TableHeader.Alignments` is unset at runtime in goldmark. Alignment is per-cell via `TableCell.Alignment`.

## Pass 1 — Measure

```go
// pseudocode
for each row:
  for each cell:
    render cell content to ANSI string (via renderSubtree)
    store content + alignment
    track maxWidth per column (displayWidth of stripped content)

// Cap columns
overhead = 3 * (numCols - 1)  // " │ " separators
maxPerCol = (wordWrap - overhead) / numCols
if maxPerCol < 20 { maxPerCol = 20 }
```

## Pass 2 — Render

```go
// Wrap each cell's content to its column width
for each cell:
    cell.lines = wrapCellLines(content, colWidth)

// Render rows
render header
hline()  // ────────┼──────────┼────────

for each body row:
    render row
    if not last:
        hline()
```

## Visual Structure (no outer borders)

```
 H1       │ H2       │ H3
──────────┼──────────┼──────────
 A        │ B        │ C
──────────┼──────────┼──────────
 D        │ E        │ F
```

- **No outer borders** — no `┌─┬─┐` or `└─┴─┘`, no leading/trailing `│`
- Separator row uses `───┼───┼───` without tee connectors
- Cleaner look, easier to nest inside lipgloss frames

## Cell Rendering Details

- Leading/trailing padding: `" "` at start and end of each cell
- Column separator: `\x1b[38;5;8m│\x1b[39m` (gray `│`)
- After cell content, emit `\x1b[39m` (reset fg) + doc background (`\x1b[48;2;R;G;Bm`) before padding to prevent inline-style background bleed
- Alignment: `v` left (default), `^` center, `>` right

## Word Wrap in Cells

See `feature-word-wrap.md` for the full word wrap algorithm, which handles CJK/emoji character-level wrapping inside table cells.
