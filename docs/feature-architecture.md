# 架构

## 概述

Gruff 是一个 Markdown 到 ANSI 终端渲染器，运行时除了 goldmark（解析器）、clipperhouse/displaywidth（Unicode 宽度）、mattn/go-runewidth（退路宽度）之外，没有其他依赖。

```
输入 (markdown 字符串)
    │
    ▼
goldmark 解析器 ──► AST
    │
    ▼
nodeRenderer ──► ANSI 字符串 (原始)
    │
    ▼
wrapText (文档级自动换行)
    │
    ▼
输出 (ANSI 字符串)
```

## 核心文件

| 文件 | 职责 |
|---|---|
| `gruff.go` | 公开 API（`Render`、`RenderBytes`、`Options`、选项函数） |
| `renderer.go` | goldmark AST 遍历器 → ANSI 输出（块级/内联渲染） |
| `ansi.go` | SGR 转义码、`Style`/`Color`/`Theme` 类型、辅助函数 |

## 关键设计决策

1. **内联样式不使用 `\x1b[0m`（完全重置）** — 使用特定的撤销代码（`\x1b[22m` 取消粗体、`\x1b[39m` 默认前景色、`\x1b[49m`/文档背景色）来保留嵌套时的外层样式状态。

2. **两遍表格渲染** — 第一遍收集所有单元格内容并计算列宽；第二遍使用 UTF-8 制表符分隔符渲染。参见 `feature-table.md`。

3. **ANSI 感知自动换行** — `wrapText`（文档级）和 `wrapCellLines`（单元格级）使用 `displayWidth`（clipperhouse/displaywidth）仅计算可见字符，保留前导空格。

4. **通过 lipgloss v2 处理颜色** — 颜色为字符串（十六进制 `"#RRGGBB"` 或调色板 `"0"-"255"`）。十六进制解析委托给 `lipgloss.Color(s)`。十六进制输出 24 位真彩色，调色板索引输出 8 位。参见 `feature-color.md`。

5. **文档背景色** — `Theme.Background` 在文档开始时输出 `\x1b[48;2;R;G;Bm`；当样式设置了背景色时，`Style.end(bg)` 恢复文档背景色而不是 `\x1b[49m`。

## 公开 API

```go
func Render(source string, opts ...Option) (string, error)
func RenderBytes(source []byte, opts ...Option) ([]byte, error)

type Options struct {
    Theme    Theme
    WordWrap int
}

func WithDark() Option          // 使用深色主题
func WithLight() Option         // 使用浅色主题
func WithWordWrap(n int) Option // 设置最大行宽（默认 120）
```
