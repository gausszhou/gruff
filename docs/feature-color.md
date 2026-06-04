# 颜色系统

## 类型

```go
type Color string
```

`Color` 可以保存：
- **24 位十六进制** — `"#RRGGBB"` 或 `"#RGB"`（简写），例如 `"#141414"`
- **8 位调色板索引** — `"0"`-`"255"`，例如 `"236"`（深灰色）

## 十六进制解析

十六进制颜色由内建 `hexRGB()` 直接解析，返回 R/G/B 值：

```go
func hexRGB(c string) (r, g, b uint8) {
    if len(c) < 7 || c[0] != '#' {
        return 0, 0, 0
    }
    return hex(c[1], c[2]), hex(c[3], c[4]), hex(c[5], c[6])
}
```

## ANSI 代码生成

| 输入 | 输出 | 格式 |
|---|---|---|
| `"#RRGGBB"`（十六进制） | `\x1b[38;2;R;G;Bm`（前景色）/ `\x1b[48;2;R;G;Bm`（背景色） | 24 位真彩色 |
| `"0"`-`"7"`（4 位） | `\x1b[3Nm`（前景色）/ `\x1b[4Nm`（背景色） | ANSI 3/4 位 |
| `"8"`-`"255"`（8 位） | `\x1b[38;5;Nm`（前景色）/ `\x1b[48;5;Nm`（背景色） | ANSI 8 位 |

## 内置调色板常量

```go
cBlack="0", cMaroon="1", cGreen="2",  cOlive="3",
cNavy="4",  cPurple="5", cTeal="6",   cSilver="7",
cGrey="8",  cRed="9",    cLime="10",  cYellow="11",
cBlue="12", cFuchsia="13", cCyan="14", cWhite="15",
cDarkBG="236"
```

## 文档背景色

`Theme.Background` 设置整个文档的背景色（深色主题默认为 `"#141414"`，浅色主题默认为 `""`）。

- 在文档开始时输出一次：`\x1b[48;2;20;20;20m`
- 在设置了背景色的样式（例如 `Code` 的 `Bg: "236"`）之后恢复：`Style.end(bg)` 输出 `\x1b[48;2;20;20;20m` 而不是 `\x1b[49m`
- 在表格单元格中，背景色在内容之后、填充之前恢复，防止内联代码背景色扩散

## 样式

```go
type Style struct {
    Fg        Color  // 前景色
    Bg        Color  // 背景色
    Bold      bool
    Italic    bool
    Underline bool
}

func (s Style) start() ansiCode   // 输出 ANSI 代码进入样式
func (s Style) end(bg Color) ansiCode  // 输出 ANSI 代码退出样式，恢复文档背景色
```

## 主题

```go
type Theme struct {
    Background              Color
    H1, H2, H3, H4, H5, H6 Style
    Strong                  Style
    Em                      Style
    Code, Link, LinkURL     Style
    Bullet, Numbered        Style
}
```

### 深色主题（默认）

| 元素 | 样式 |
|---|---|
| 背景色 | `"#141414"` |
| H1 | 粗体 + 下划线 + 白色 |
| H2 | 粗体 + 黄色 |
| H3 | 粗体 + 绿色 |
| H4 | 粗体 + 青色 |
| H5 | 粗体 + 灰色 |
| H6 | 灰色 |
| Code | 前景色：白色 + 背景色：236（深灰色） |
| Link | 下划线 + 青色 |

### 浅色主题

| 元素 | 样式 |
|---|---|
| 背景色 | (无) |
| H1 | 粗体 + 下划线 + 黑色 |
| H2 | 粗体 + 深蓝色 |
| Code | 前景色：黑色 + 背景色：银色 |
| Link | 下划线 + 深蓝色 |
