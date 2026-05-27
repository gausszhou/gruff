# Style 架构优化

## 背景

将 `Theme.Background` 和 `Theme.Padding` 归入 `Theme.Document` Style，统一了样式继承体系。

## 变更

### 旧结构

```go
type Theme struct {
    Background string       // 文档背景色
    Padding    int          // 文档内边距
    H1, H2, ... Style
}
```

### 新结构

```go
type Theme struct {
    Document Style          // 包含 Bg + Padding
    H1, H2, ... Style
}
```

## 背景继承规则

- 文档背景色由 `Theme.Document.Bg` 设置
- 每个节点类型的 `Style` 如果未设置 `Bg`，则继承 `Document.Bg`
- `Style.end(bg)` 中的 `bg` 参数接收 `Document.Bg`，当节点有自定义背景时用于恢复
- 避免使用 `\x1b[0m`（全重置），改用精确撤销码（`\x1b[22m`、`\x1b[39m`、`\x1b[49m`）

## 内边距继承

- 文档级内边距由 `Theme.Document.Padding` 设置
- 元素级内边距由对应 `Style.Padding` 设置（如代码块）
- `wrapText` 接收 `Theme.Document.Padding` 作为左右边距

## 示例

```go
var darkTheme = Theme{
    Document: Style{Bg: "#141414", Padding: 2},
    H1:       Style{Bold: true, Fg: "#FFFF87", Bg: "#5F5FFF"},
    Code:     Style{Fg: "#FF5F5F", Bg: "#303030", Padding: 1},
}
```
