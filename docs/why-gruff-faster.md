# 为什么 Gruff 性能优于 Glamour

## 基准测试数据

测试环境：AMD Ryzen 7 4800H，输入为 CommonMark 规范文档（~200KB，~9700 行）。

| 指标 | Gruff | Glamour (minimal) | 差距 |
|---|---|---|---|
| 耗时 | 19.4 ms | 263 ms | **13.5x** |
| 内存分配 | 4.9 MB | 22.3 MB | **4.6x** |
| 分配次数 | 17,751 | 2,489,723 | **140x** |

> Glamour 已配置为 minimal 模式：`CodeBlock.Chroma = nil`、`WordWrap = 0`、`TableWrap = false`、`InlineTableLinks = true`，所有 Gruff 不支持的元素（BlockQuote、Strikethrough、TaskCheckBox、HTMLBlock、DefinitionList 等）样式已归零。

## 逐点分析

### 1. 渲染架构：直通 vs 缓冲链

**Gruff** 遍历 AST 时直接调用 `io.Writer.Write`，不走中间缓冲。

```
AST → switch → Write(ANSI)     ← 直通
```

**Glamour** 对每个 block 级元素（Document、Paragraph、Heading、ListItem 等）创建一个 `bytes.Buffer` 推入 block stack，子节点内容先写入当前 buffer，结束时再 flush 到父级 buffer。层级越深，内容副本越多。

```
AST → Element{Buffer} → push stack → 子节点写入 Buffer → Finish → flush → pop
```

**对 alloc 的影响**：每 paragraph 分配一个 `BlockElement{Block: &bytes.Buffer{}}`，每 heading 同样分配。对于 9700 行的文档，数千次 buffer 分配 + 逐级复制，是 140x alloc 差距的主要来源。

---

### 2. 元素创建方式：switch 值 vs 对象分配

**Gruff** 在 `renderNode()` 中用一个 `switch n.Kind()` 直接处理，不创建额外对象。

```go
switch n.Kind() {
case ast.KindParagraph:
    // 直接输出
case ast.KindHeading:
    // 直接输出
}
```

**Glamour** 在 `NewElement()` 中为每个 AST 节点创建一个 `Element` 结构体 + 具体 `Renderer` 对象（`ParagraphElement`、`HeadingElement`、`LinkElement` 等）。每个都是 heap alloc。

**对 alloc 的影响**：每个 AST 节点多 1-2 次 heap alloc。9700 行文档的 AST 节点数以万计。

---

### 3. 链接处理：直接输出 vs URL 解析 + FNV 哈希 + OSC-8

**Gruff** 对链接直接输出文本和 URL，不做任何额外处理。

```go
case *ast.Link:
    r.renderChildren(n)        // 链接文本
    w.Write(r.th.LinkURL, url) // URL
```

**Glamour** 对每个 Link 和 AutoLink 调用 `makeHyperlink()`，其中执行：
- `url.Parse(link)` — URL 解析
- `fnv.New32a()` + `io.WriteString(h, link)` — FNV-32a 哈希
- `ansi.SetHyperlink(link, urlID)` / `ansi.ResetHyperlink()` — OSC-8 超链接 ANSI 序列生成

此外 `LinkElement.renderTextPart()` 为每个子节点创建中间 `bytes.Buffer`。

**对 alloc 的影响**：链接越多开销越大。CommonMark 规范文档包含大量 URL 引用。

---

### 4. 文本处理：原始输出 vs HTML 反义 + 转义替换

**Gruff** 对 `ast.Text` 直接写入 `n.Segment.Value`，不做任何文本变换。

**Glamour** 在 `BaseElement.doRender()` 中执行：
- `html.UnescapeString(s)` — 每段文本做 HTML 反义
- `escapeReplacer.Replace(s)` — 17 组 Markdown 转义字符替换（`\\\\` → `\\`、`\\*` → `*` 等），无论文本中是否有转义符

```go
var escapeReplacer = strings.NewReplacer(
    "\\\\", "\\", "\\`", "`", "\\*", "*",
    "\\_", "_", "\\{", "{", "\\}", "}",
    "\\[", "[", "\\]", "]", "\\<", "<",
    "\\>", ">", "\\(", ")", "\\)", ")",
    "\\#", "#", "\\+", "+", "\\-", "-",
    "\\.", ".", "\\!", "!", "\\|", "|",
)
```

`strings.NewReplacer` 每次 `doRender` 都创建新 string（即使没有转义字符）。

---

### 5. 样式级联：直接查表 vs 逐级合并

**Gruff** 用 `headingStyle()` 直接根据 heading level 查 Theme 字段（O(1)）。

```go
func headingStyle(lv int, th Theme) Style {
    switch lv {
    case 1: return th.H1
    case 2: return th.H2
    // ...
    }
}
```

**Glamour** 对每个元素调用 `cascadeStyle()` / `cascadeStylePrimitives()`，从 block stack 当前样式逐字段合并父级样式和子级样式：

```go
func cascadeStylePrimitive(parent, child StylePrimitive, ...) StylePrimitive {
    s := child
    s.Color = parent.Color
    s.BackgroundColor = parent.BackgroundColor
    // ... 复制十几个字段
    if child.Color != nil { s.Color = child.Color }
    // ...
}
```

虽为值复制（stack alloc），但叠加大量节点后复制负担可观。

---

### 6. 自动换行：文档级单次 vs 段落级多次

**Gruff** 的 `wrapText` 在整个文档渲染完成后对**最终输出**做一次换行，所有段落共用一次处理。

**Glamour** 在每个 `ParagraphElement.Finish()` 和 `HeadingElement.Finish()` 中分别调用 `lipgloss.Wrap()`，每个段落/标题独立执行换行计算。

虽然通过 `WithWordWrap(0)` 可跳过 lipgloss，但 `BlockElement` 的 buffer 分配仍然存在。

---

### 7. 表格链接处理：无

**Gruff** 不做表格内链接的特殊处理。

**Glamour** 默认（`InlineTableLinks = false`）在表格渲染时：
- 用 `collectLinksAndImages()` 遍历整个表格 AST 子树
- 收集所有 AutoLink、Link、Image 节点
- 去重后替换为 `[N]` 引用标记
- 在表格下方渲染脚注列表

通过 `WithInlineTableLinks(true)` 可跳过此处理。

---

### 8. 额外扩展渲染器注册（即使 style 为空也执行）

Glamour 的 `ANSIRenderer.RegisterFuncs()` 为以下元素注册独立渲染器：

| 元素 | 即使 style 归零仍执行的操作 |
|---|---|
| `BlockQuote` | 创建 `BlockElement` + 推 block stack |
| `Strikethrough` | 创建 `BaseElement` + ANIS 样式渲染 |
| `TaskCheckBox` | 创建 `TaskElement`（ListItem 内） |
| `DefinitionList` | 创建 `BlockElement` + 推 block stack |
| `HTMLBlock` | `ctx.SanitizeHTML()` HTML 清理 |
| `RawHTML` | `ctx.SanitizeHTML()` HTML 清理 |
| `Footnote` / `FootnoteList` | block stack 操作 |
| `Emoji` | Emoji 渲染 |

Gruff 对这些元素不注册任何特殊处理，子节点直接 fallthrough 到 `default: r.renderChildren(n)`。

---

### 9. ANSI 分层样式系统

**Gruff** 每个 Style 直接输出 SGR 序列（`start()` / `end()`），不经过中间层。

**Glamour** 通过 `lipgloss.Style` / `ansi.Style` 链式 API 构建 ANIS 字符串，每次 `renderText()` 都：
```go
style := ansi.Style{}
style = style.ForegroundColor(lipgloss.Color(...))
style = style.BackgroundColor(lipgloss.Color(...))
style = style.Bold()
// ...
style.Styled(s)
```

`lipgloss.Color(s)` 每次解析颜色字符串（hex 或 0-255），构建样式对象，再调用 `Styled()` 生成 ANSI 码。

---

### 10. 模板引擎用于 ImageText 格式化

Glamour 的 `ImageElement` 渲染时调用 `formatToken()`：

```go
func formatToken(format string, token string) (string, error) {
    var b bytes.Buffer
    v := map[string]interface{}{"text": token}
    tmpl, _ := template.New(format).Parse(format)
    tmpl.Execute(&b, v)
    return b.String(), nil
}
```

使用 `text/template` 解析和执行模板，即使格式为空字符串，`template.New(format).Parse(format)` 也会执行一次。

**Gruff** 对 Image 只输出文本子节点，无模板处理。

## 总结

| 对比维度 | Gruff | Glamour | 差距来源 |
|---|---|---|---|
| 渲染路径 | 直通 io.Writer | block stack + bytes.Buffer 链 | 架构设计 |
| 元素分发 | switch 值 | Element + Renderer 对象分配 | 对象模型 |
| 链接处理 | 无 | URL 解析 + FNV 哈希 + OSC-8 | 功能冗余 |
| 文本处理 | 原始输出 | HTML 反义 + 转义替换 | 功能冗余 |
| 样式合并 | O(1) 查表 | cascadeStyle 逐级复制 | 架构设计 |
| 自动换行 | 文档级 1 次 | 段落级 N 次 | 架构设计 |
| 代码高亮 | 无 | chroma（可关） | 功能冗余 |
| 表格链接 | 无 | AST 遍历 + 收集 + 去重 | 功能冗余 |
| 额外渲染器 | 无（fallthrough） | ~10 个专用渲染器 | 功能冗余 |
| 颜色系统 | 裸 SGR | lipgloss / ansi.Style 链式 API | 分层设计 |

**核心结论**：Gruff 的每像素 3-5 条指令哲学使渲染路径保持极简。Glamour 的功能丰富性（OSC-8 超链接、HTML 清理、chroma 语法高亮、表格链接脚注、自定义模板等）带来了架构性开销，即使通过配置关闭所有额外功能，其渲染框架本身（block stack、Element 对象分配、style cascading、html unescaping）仍然导致 13x 耗时和 140x 分配次数的差距。
