# gruff

一个轻量级、高性能的 [Go](https://go.dev) 库，用于将 Markdown 渲染为带 ANSI 颜色的终端输出。

## 特性

- **标题**（H1–H6），每个级别有独立的样式
- **粗体**、*斜体* 和 ***粗斜体***
- `行内代码` 和围栏/缩进**代码块** —— 绿色文本，语言标签为灰色
- [链接](https://github.com/gausszhou/gruff) —— 加粗链接文本 + 带下划线的 URL，支持 OSC 8 终端超链接
- 无序（`-`、`*`）和有序（`1.`）列表
- GFM 表格，带有 **UTF-8 框线边框**、对齐（左/中/右）、单元格内联格式，以及**最大列宽下自动换行**
- 删除线文本（`~~删除~~`）
- 任务列表（`- [x] 已完成`、`- [ ] 待办`）
- 水平分割线（`---`）
- 深色和浅色主题
- ANSI 感知的自动换行（`WithWordWrap`）

## 安装

```bash
go get github.com/gausszhou/gruff/gruff
```

### CLI

```bash
go install github.com/gausszhou/gruff@latest
gruff README.md              # 渲染一个文件
gruff render README.md       # 显式子命令
gruff render -w 40 README.md # 自定义换行宽度
gruff render -l README.md    # 浅色主题
```

## 用法

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

### 选项

| 函数 | 描述 |
|----------|-------------|
| `WithDark()` | 深色背景主题（默认） |
| `WithLight()` | 浅色背景主题 |
| `WithWordWrap(n)` | 在 `n` 列处换行 |

```go
out, err := gruff.Render(source, gruff.WithLight(), gruff.WithWordWrap(80))
```

## 性能

针对 `testdata/benchmark.md`（约 5.6 KB）重复 100 次（约 560 KB 输入）进行基准测试，
硬件为 Intel Core Ultra 7 255H。

| 指标         | gruff ¹     | glamour (minimal) ² | glamour (standard) ³ | 提升（对比 minimal / 对比 standard）|
|----------------|-------------|---------------------|----------------------|------------------------------------------|
| Time/op        | **~98 ms**  | ~435 ms             | ~1.39 s              | **~4.4× / ~14×**                         |
| Memory/op      | **~69 MB**  | ~137 MB             | ~441 MB              | **~2.0× / ~6.4×**                        |
| Allocations/op | **~459,000**| ~10,100,000         | ~39,300,000          | **~22× / ~86×**                          |

¹ gruff：`WithDark()`（无背景），`WithWordWrap(120)`。
² glamour minimal：`Chroma = nil`、`CleanInput`，关闭自动换行、表格换行，行内表格链接开启。
³ glamour standard：`WithStandardStyle("dark")`，120 列自动换行。

参见 [`docs/why-gruff-faster.md`](docs/why-gruff-faster.md) 了解性能差距的详细分析。

本地运行基准测试：

```bash
go test -bench=. -benchmem ./benchmark/
```

## 示例

可以直接运行的示例位于 [`examples/`](examples/) 目录：

| 示例 | 描述 |
|---------|-------------|
| [`basic`](examples/basic/) | 使用 CLI 标志（`--light`、`--wrap`）渲染 markdown |
| [`table`](examples/table/) | 表格专项演示，展示对齐和自动换行 |
| [`codeblock`](examples/codeblock/) | 带语言标签的代码块渲染 |
| [`custom-theme`](examples/custom-theme/) | 自定义 ANSI 颜色和样式定制 |
| [`api`](examples/api/) | `Render`、`RenderBytes` 和 `WithWordWrap` 用法 |
| [`compare-benchmark`](examples/compare-benchmark/) | 用 gruff vs glamour 渲染并排的基准 markdown |
| [`compare-glamour`](examples/compare-glamour/) | 并排的 glamour standard vs minimal |
| [`compare-theme`](examples/compare-theme/) | 并排的 gruff 深色 vs 浅色主题 |
| [`compare-simple`](examples/compare-simple/) | 无 viewport/bubbletea 的 glamour minimal vs standard |
| [`viewport-gruff`](examples/viewport-gruff/) | gruff 输出置于带滚动条的 bubbletea viewport 中 |
| [`viewport-glamour`](examples/viewport-glamour/) | glamour 输出置于带滚动条的 bubbletea viewport 中 |

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

## 主题定制

使用导出的 `Theme`、`Style` 和 `Color` 类型自定义颜色和样式：

```go
import "github.com/gausszhou/gruff"

customTheme := func() gruff.Option {
    return func(o *gruff.Options) {
        o.Theme.H1 = gruff.Style{Fg: gruff.Color(196), Bold: true}       // 红色
        o.Theme.Strong = gruff.Style{Bold: true, Fg: gruff.Color(51)}    // 青色
        o.Theme.Code = gruff.Style{Fg: gruff.Color("#50865a")}            // 绿色
        o.Theme.Link = gruff.Style{Bold: true, Fg: gruff.Color("#5c9cf5")}
        o.Theme.LinkURL = gruff.Style{Underline: true, Fg: gruff.Color("#5c9cf5")}
    }
}

out, _ := gruff.Render(md, customTheme())
```

## 工作原理

解析由 [`goldmark`](https://github.com/yuin/goldmark) 完成。`renderer.go` 中的递归类型 switch 遍历 AST，直接输出 SGR ANSI 码 —— 没有中间 DOM、没有 CSS、没有 HTML。每个内联样式使用精确定向的撤销码（例如 `\x1b[22m` 表示取消加粗、`\x1b[39m` 表示默认前景色），而不是 `\x1b[0m`，从而正确处理嵌套格式。

表格渲染采用两遍算法：先收集所有单元格内容并计算列宽，然后输出 UTF-8 框线边框和填充后的单元格内容。列宽设有最大值上限，超出部分自动换行。

代码块逐行渲染，语言标签用灰色显示、内容用绿色显示，每行填充到文档宽度，以获得干净的全宽效果。

## 依赖

**运行时：** `github.com/yuin/goldmark`、`github.com/mattn/go-runewidth`

**仅示例：** `charm.land/bubbles/v2`、`charm.land/bubbletea/v2`、`charm.land/lipgloss/v2`、`charm.land/glamour/v2`（不包含在库构建中）

## 许可证

MIT