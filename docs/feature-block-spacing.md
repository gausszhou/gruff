# 块级元素间距方案

## 问题

当前每个块级元素各自管理尾部换行，策略不统一：

| 元素 | 尾部 | 条件 |
|------|------|------|
| Paragraph | `\n\n` / `\n` / 无 | `inBlockquote`, `isInsideList`, `isInsideTable` |
| Heading | `\n\n`（无条件） | — |
| List | `\n`（顶层） / 无（嵌套） | `listParent(n) == nil` |
| ListItem | `\n`（无嵌套列表时） | `hadNestedList` |
| Blockquote | 无 | — |
| CodeBlock | `\n` | — |
| ThematicBreak | `\n\n`（无条件） | — |
| Table | `\n`（额外写） | — |

根因：**spacing 责任分散在各元素中**，全局 `TrimRight` + `wrapText` 叠加后产生不可预期的空行。

## 方案

**将所有块级元素的尾部统一为单个 `\n`，由父节点在相邻块级兄弟间插入额外 `\n` 形成空行。**

### 改动总览

```
UnifiedSpacing = element's own \n + (if parent adds separator) \n
```

效果：元素间 `\n\n` = 一个空行；元素末尾只有一个 `\n`，被 `TrimRight` 干净移除。

### 具体修改

#### 1. 新增 `isBlockLevel(n ast.Node) bool`

```go
func isBlockLevel(n ast.Node) bool {
    switch n.(type) {
    case *ast.Paragraph, *ast.Heading, *ast.List, *ast.Blockquote,
         *ast.FencedCodeBlock, *ast.CodeBlock, *ast.ThematicBreak,
         *extensionAst.Table:
        return true
    }
    return false
}
```

#### 2. Document 渲染时在块级兄弟间插入空行

```go
case *ast.Document:
    r.buf.WriteString(string(ansiBg(r.th.Bg)))
    for c := n.FirstChild(); c != nil; c = c.NextSibling() {
        if c != n.FirstChild() && isBlockLevel(c) {
            r.buf.WriteByte('\n') // 空行分隔
        }
        r.renderNode(c)
    }
```

这是整个方案的核心：**spacing 集中在父节点一处处理**。

#### 3. 各元素统一尾部为 `\n`

| 元素 | 当前 | 改为 |
|------|------|------|
| Paragraph (L38-44) | `\n\n`（普通）/ `\n`（blockquote）/ 无（list/table） | `\n`（无条件） |
| Heading (L51) | `\n\n` | `\n` |
| ThematicBreak (L147) | `\n\n` | `\n` |
| Table (L582) | `r.buf.WriteByte('\n')` | 删除 |
| List (L55-57) | `\n`（顶层） | 删除 |
| ListItem (L290) | 嵌套列表前写 `\n` | 删除 |
| ListItem (L299-301) | `!hadNestedList` 时写 `\n` | 无条件写 `\n` |

#### 4. 删除不再使用的辅助函数

Paragraph 的条件移除后，`isInsideList` 和 `isInsideTable` 不再被任何元素使用，可删除。

### 效果分析

#### 顶层元素间距

```
Document 子节点遍历:

  Paragraph1\n    ← 元素自身 \n
  \n              ← Document 插行
  Heading\n       ← 元素自身 \n
  \n              ← Document 插行
  CodeBlock\n     ← 元素自身 \n
  \n              ← Document 插行
  Table...\n      ← 元素自身 \n

→ Paragraph1\n\nHeading\n\nCodeBlock\n\nTable...\n
→ TrimRight → Paragraph1\n\nHeading\n\nCodeBlock\n\nTable...
```

所有顶层块元素之间均有且只有一个空行。

#### 列表内部

```
listItem1\n  ← ListItem 无条件 \n（内部 Paragraph 也写了一个 \n）
listItem2\n  ← ListItem 无条件 \n（内部 Paragraph 也写了一个 \n）
```

注意：Paragraph `\n` + ListItem `\n` = `\n\n`，列表项之间会多一个空行。这是"完全统一"的必然代价——如果要求列表项紧凑，需保留 Paragraph 的 `isInsideList` 条件。

#### Blockquote 内部

```
│ P1\n   ← Paragraph 的 \n
│ \n     ← Blockquote 的子元素间分隔 prefix+\n
│ P2\n   ← Paragraph 的 \n
```

与当前行为一致，Blockquote 自身控制子元素间距，不受 Document 插行影响。

#### 表格内部

Table 只移除末尾额外的 `r.buf.WriteByte('\n')`（L582）。行内 `\n`（来自 `renderTableRow` 和 `hline`）保持不变，为结构性换行。

### 验证场景

```
# Title

Paragraph text.

> Blockquote paragraph 1.
> Blockquote paragraph 2.

- List item 1
- List item 2

1. Ordered 1
2. Ordered 2

```code
code block
```

| A | B |
|---|---|
| 1 | 2 |

---

Thematic break above.
```

所有相邻块元素间应有且仅有一个空行。

### 风险与备选

| 风险 | 影响 | 备选 |
|------|------|------|
| 列表项间出现空行 | Paragraph `\n` + ListItem `\n` = `\n\n` | Option A: 保留 Paragraph 的 `isInsideList` 检查 |
| Blockquote 内段落间距变化 | 不变（Blockquote 自己控制） | 无需处理 |
| 表格后出现多余空格 | 移除了尾部 `\n`，行为干净 | 无需处理 |
