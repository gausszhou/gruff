# 表格渲染

表格通过两遍算法在 goldmark GFM 表格 AST 上进行渲染。

## AST 结构

```
Table
├── TableHeader
│   ├── TableCell (对齐方式: cell.Alignment)
│   ├── TableCell
│   └── ...
├── TableRow
│   ├── TableCell
│   ├── TableCell
│   └── └── ...
```

**注意：** `TableHeader.Alignments` 在 goldmark 运行时未设置。对齐方式通过 `TableCell.Alignment` 按单元格设置。

## 第一遍 — 测量

```go
// 伪代码
for each row:
  for each cell:
    将单元格内容渲染为 ANSI 字符串（通过 renderSubtree）
    存储内容 + 对齐方式
    跟踪每列的最大宽度（剥离 ANSI 后的 displayWidth）

// 限制列宽
overhead = 3 * (numCols - 1)  // " │ " 分隔符
maxPerCol = (wordWrap - overhead) / numCols
if maxPerCol < 20 { maxPerCol = 20 }
```

## 第二遍 — 渲染

```go
// 将每个单元格的内容换行到其列宽
for each cell:
    cell.lines = wrapCellLines(content, colWidth)

// 渲染行
render header
hline()  // ────────┼──────────┼────────

for each body row:
    render row
    if not last:
        hline()
```

## 视觉结构（无外边框）

```
 H1       │ H2       │ H3
──────────┼──────────┼──────────
 A        │ B        │ C
──────────┼──────────┼──────────
 D        │ E        │ F
```

- **无外边框** — 没有 `┌─┬─┐` 或 `└─┴─┘`，没有开头/结尾的 `│`
- 分隔行使用 `───┼───┼───`，无 T 形连接符
- 更简洁的外观，更易于嵌套在 lipgloss 框架内

## 单元格渲染细节

- 首尾填充：每个单元格开头和结尾的 `" "`
- 列分隔符：`\x1b[38;5;8m│\x1b[39m`（灰色 `│`）
- 单元格内容后，先输出 `\x1b[39m`（重置前景色）+ 文档背景色（`\x1b[48;2;R;G;Bm`），然后再填充，防止内联样式背景色扩散
- 对齐方式：`v` 左对齐（默认）、`^` 居中、`>` 右对齐

## 单元格中的自动换行

参见 `feature-word-wrap.md` 了解完整的自动换行算法，该算法处理表格单元格内的 CJK/表情符号字符级换行。
