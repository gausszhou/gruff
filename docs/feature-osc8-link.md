# OSC 8 超链接支持

## glamour 的实现方式

glamour 处理链接分两层：

### 1. 链接渲染 (`link.go`)

glamour 渲染链接时不直接写 ANSI 输出，而是生成一个 `token` 字符串，再用其风格系统包裹：

```go
// link.go
token := e.hyperlink + b.String() + e.resetHyperlink
el := &BaseElement{Token: token, Style: ctx.options.Styles.LinkText}
el.Render(w, ctx)
```

其中 `e.hyperlink` / `e.resetHyperlink` 通过 `makeHyperlink()` 生成：

```go
// 用 FNV32a 哈希生成稳定 id= 参数，辅助终端复用已解析的链接
hyperlink = ansi.SetHyperlink(link, urlID)   // → \x1b]8;id=xxx;url\x1b\\
resetHyperlink = ansi.ResetHyperlink()       // → \x1b]8;;\x1b\\
```

### 2. 换行的换行处理 (`wrap.go`)

glamour 的处理是通过 lipgloss 的 `WrapWriter`：

```go
type WrapWriter struct {
    w     io.Writer
    p     *ansi.Parser
    style uv.Style   // 当前活跃的 CSI 样式
    link  uv.Link    // 当前活跃的 OSC 8 链接
}
```

关键机制：
- **CSI 回调**：当解析器遇到 `\x1b[...m` 时，解析参数更新 `w.style`
- **OSC 回调**：当解析器遇到 `\x1b]8;...\x1b\\` 时，解析数据更新 `w.link`

换行处理的核心逻辑（`Write` 方法）：

```
遇到 \n 时:
  1. 关闭: 如果 style 非零，输出 \x1b[0m (全量 reset)
           如果 link  非零，输出 \x1b]8;;\x1b\\ (关闭链接)
  2. 输出 \n
  3. 重开: 如果 link  非零，输出 ansi.SetHyperlink(link.URL, link.Params)
           如果 style 非零，输出 style.String()
```

### 3. 为什么 glamour 能正确处理

glamour 不需要在 `wrapText` 中拆分单词——它在 Paragraph 的 `Finish()` 阶段把整个段落的已渲染内容交给 `lipgloss.Wrap()` 一次性处理。`WrapWriter` 逐字节扫描，利用 `ansi.Parser` 持续追踪当前活跃的样式和链接状态。每遇到 `\n` 就 close → write \n → reopen。

### 4. 为什么 gruff 不同

gruff 的架构是「先渲染出完整 ANSI 字符串，再用自研的 `wrapText` 做换行」。没有 lipgloss / chrono 的 ANSI parser 基础设施。`wrapText` 是一个 rune 级别的状态机，必须自己跟踪样式和 OSC 8 的活跃状态。

## gruff 的实现思路

### 核心挑战

`wrapText` 按 rune 逐字处理，以空格分隔单词。旧实现把 CSI/OSC 序列当 word 的一部分累积——当 flushWord 清空 word 后，后续单词丢失了前面的样式和链接状态。

例如渲染 `[a really long label text](https://example.com)` 得到：

```
\x1b]8;id=https://example.com;https://example.com\x1b\\\x1b[1m\x1b[38;2;92;156;245ma really long label text\x1b[22m\x1b[39m\x1b]8;;\x1b\\
```

旧 wrapText 处理：
- word "a" 累积了 OSC+CSI，flush 后丢失
- word "really" 无样式，无链接
- 换行后无样式回放

### 解决方案：activeStyle 追踪

引入 `activeStyle []byte`，在 `wrapText` 的处理循环中**全局维护当前活跃的转义序列**：

#### 状态机变化

```
旧: esc → 累积到 word → flush 时写入（activeStyle 未追踪，换行后样式丢失）
新: esc 仍累积在 word 内（写入时直接输出），同时更新 activeStyle 作为换行时样式重放的依据
```

#### activeStyle 更新规则

| 遇到的序列 | 对 activeStyle 的操作 |
|-----------|---------------------|
| `\x1b[1m` (bold on) | 追加 |
| `\x1b[22m` (bold off) | 移除 `\x1b[1m` |
| `\x1b[3m` / `\x1b[23m` | 同上 italic |
| `\x1b[4m` / `\x1b[24m` | 同上 underline |
| `\x1b[38;2;R;G;Bm` (fg) | 移除旧 fg → 追加新 fg |
| `\x1b[39m` (fg 默认) | 移除 fg |
| `\x1b[48;...m` (bg) | 同上 bg |
| `\x1b]8;id=url;url\x1b\\` (OSC 8) | 移除旧 → 追加新 |
| `\x1b]8;;\x1b\\` (OSC end) | 移除 OSC 8 |

#### 换行处理（模仿 glamour 的 close-write-reopen）

```
newLine():
  1. 如果 activeStyle 非空:
     输出 \x1b[0m       // 全量关闭所有 CSI 样式
     输出 osc8End        // 关闭 OSC 8 链接
  2. 填充当前行至 width
  3. 输出 \n
  4. 输出 bgCode + padding 空格
  5. 输出 activeStyle (重放活跃样式+链接)
  6. 重置 lineLen = padding
```

#### 换行时样式重放

`newLine()` 通过 `activeStyle` 实现：

```
newLine():
  1. 如果 activeStyle 非空:
      输出 \x1b[0m       // 全量关闭所有 CSI 样式
      输出 osc8End        // 关闭 OSC 8 链接
  2. 填充当前行至 width
  3. 输出 \n
  4. 输出 bgCode + padding 空格
  5. 输出 activeStyle (重放活跃样式+链接)
  6. 重置 lineLen = padding
```

#### flushWord 的换行判断

当 word 在当前行放不下时，调用 `newLine()` 将 word 移至下一行。`flushWord` 中的第一次换行（word 整体不拟合当前行）使用 `wordStartStyle`（该 word 开始前的活跃样式），确保新行上 OSC8 + SGR 样式正确重发。

#### 长词断行 (breakChunk)

单词过长时（如 URL `https://very-long-url...`），`breakChunk` 通过 `style *[]byte` 参数在切分的同时累积追踪 ANSI 转义码状态。断点处的样式信息用于 `newLine()`，保证续行样式/链接正确。切分循环结束后，累积的 `splitStyle` 覆盖全局 `activeStyle` 和 `wordStartStyle`，确保后续文本的样式上下文正确。

### 文件变更清单

| 文件 | 内容 |
|------|------|
| `gruff_ansi.go` | `osc8Link`/`osc8End` 常量、`updateActiveStyle`、`removeExact`/`removeCSIPrefix`/`removeOSC` 辅助函数、`breakChunk` 增加 style 参数追踪 |
| `gruff.go` | `wrapText` 新增 activeStyle 全局追踪 + wordStartStyle + newLine close-write-reopen |
| `gruff_renderer.go` | `*ast.Link` 外包 OSC 8、`*ast.AutoLink` 加 OSC 8 |
| `gruff_test.go` | Link 测试加入 OSC 8、新增 TestRender_LongURL_Wrap |
