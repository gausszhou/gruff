# Word Wrap

Gruff has two word-wrap subsystems:

1. **`wrapText`** — document-level line wrapping (applied to the entire output)
2. **`wrapCellLines`** — cell-level wrapping (per column, inside tables)

## Document Wrap (`wrapText` in `gruff.go`)

Applied after rendering, controlled by `WithWordWrap(n)` (default 120, disabled with `WithWordWrap(0)`).

- Accumulates words at space boundaries
- ANSI escape sequences are transparent (counted as zero display width)
- Wraps at the nearest word boundary when line exceeds width
- Preserves existing `\n` line breaks
- Preserves leading whitespace

## Cell Wrap (`wrapCellLines` in `renderer.go`)

Applied per column during the table rendering Pass 2.

### Latin/ASCII Wrapping

- Accumulates characters into "words" at space boundaries
- When a word doesn't fit on the current line, it starts a new line
- A space separator (ASCII 0x20) is inserted between words only when `wordVisLen > 0` (not for pure-ANSI "words")

### CJK & Emoji Wrapping

CJK characters (Chinese, Japanese, Korean) and emoji have `displayWidth > 1` (double-width). They **never have spaces between them** in natural text, so space-only word boundaries don't work.

For any character where `runewidth.RuneWidth(r) > 1`:

1. Flush any pending word (Latin text accumulated before the CJK char)
2. If the double-width char doesn't fit on the current line, start a new line
3. Write the char directly (no space separator)
4. Continue to next char

This ensures correct wrapping for:

```
这是一个非常长的中文句子，它完全超过了
四十个字符的最大列宽度，必须自动换行到
单元格内的多行才能完整显示。
```

And mixed scripts:

```
Hello你好🚀 This is a long🔥 mixed🌟  
sentence with✅ emoji, CJK中文, and   
long text that wraps...
```

### ANSI Handling

ANSI escape sequences can appear anywhere in the content (from inline styles like **bold** or `code`). The wrap function:

- Treats `\x1b[...m` sequences as zero-width tokens
- Accumulates ANSI codes into the current word
- Only strips ANSI for display-width calculation (`displayWidth(stripANSI(...))`)
- Preserves ANSI codes across wrapped lines (no broken escapes)
- Does not insert spaces before pure-ANSI "words" (wordVisLen == 0)

### Column Width

Max column width is calculated as:

```go
overhead := 3 * (numCols - 1)            // " │ " per gap
maxCol := (wordWrap - overhead) / numCols // capped per col
if maxCol < 20 { maxCol = 20 }           // min 20 chars
```
