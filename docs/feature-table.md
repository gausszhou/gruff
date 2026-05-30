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
```

此时 `colWidths[i]` 为每列内容的"自然宽度"（Naturals）。

### displayWidth 与 Emoji 宽度

列宽通过 `displayWidth(stripANSI(content))` 计算，该函数处理 Emoji Variation Selector-16（U+FE0F）：

```
input:  "✅ ❌ 🚀 ⭐ ⚠️ ✨ 🎉 😀 🔥"
         ↑   ↑   ↑   ↑ ↑↑  ↑   ↑   ↑
         │   │   │   │ ││  │   │   └─ 🔥 U+1F525        → width 2
         │   │   │   │ ││  │   └───── 😀 U+1F600        → width 2
         │   │   │   │ ││  └──────── 🎉 U+1F389         → width 2
         │   │   │   │ │└──────────── ✨ U+2728          → width 2
         │   │   │   │ └───────────── ⚠️ U+26A0 + U+FE0F  → width 2 (VS16 强制 emoji)
         │   │   │   └─────────────── ⭐ U+2B50          → width 2
         │   │   └─────────────────── 🚀 U+1F680         → width 2
         │   └─────────────────────── ❌ U+274C          → width 2
         └─────────────────────────── ✅ U+2705          → width 2
                                 空格 × 8                → width 1 each
                                                         → 总计 26
```

- `go-runewidth` 对 ambiguous width 字符（如 U+26A0 ⚠）返回 1，不因 U+FE0F 自动提升
- 自定义 `displayWidth` 弥补此不足：凡是后跟 U+FE0F 的码位强制为宽度 2
- U+FE0F 本身计 0，RI pairs（Regional Indicator，如国旗）仍沿用 `runewidth` 的宽度

## 第二遍 — 列宽计算

```go
// 固定开销
overhead = 3 * (numCols - 1) + 2
// 3 = " │ " 每个列分隔符占 3 个显示宽度
// (numCols - 1) = 分隔符数量
// +2 = 首尾各 1 个空格填充（无外边框时仍保留内部间距）

naturalTotal = sum(colWidths) + overhead
maxWidth = wordWrap - padding  // 可用总宽度

if naturalTotal <= maxWidth:
    // 内容可自然放下，无需压缩
    // 仅将 <3 的列设为 3（最小列宽，保证 " │ " 可读）
    colWidths[i] = max(colWidths[i], 3)
else:
    // 内容超宽，所有列均分剩余空间
    equal = (maxWidth - overhead) / numCols
    if equal < 3:
        equal = 3  // 硬下限
    for i:
        colWidths[i] = equal  // 统一覆盖
```

关键区别：
- **不按比例分配**：当总宽超出时，所有列被强制设为相同宽度 `equal`，而非按自然宽度比例缩小。
- **最小列宽 3**：下限为 3 个显示宽度，即使内容或可用空间更小。
- **`+2` 开销**：左右边缘各 1 空格（renderTableRow 中左右各补一个空格，见下方）。
- **`padding` 影响**：`wordWrap - padding` 才是实际可用宽度（`padding` 由 `WithWordWrap` 设置，默认 0）。

## 第三遍 — 换行与渲染

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
