# Gruff vs Glamour Benchmark

This document contains a balanced mix of markdown elements for performance testing.

## Links

Another lengthy link: [goldmark extension](https://github.com/yuin/goldmark/tree/main/extension#table-extension) for table support.

Stack Overflow: [ANSI escape codes](https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences/33206814#33206814) reference.

Charm: [bubbletea mouse](https://github.com/charmbracelet/bubbletea/blob/master/tutorials/mouse.md#full-mouse-mode-without-scrolling) docs.

Autolink: Can i use <https://caniuse.com/?search=css-container-queries>

Autolink: W3C  <https://www.w3.org/TR/css-color-4/#lab-to-lch-conversion-drafts>

Stack Overflow: Plain Text https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences/33206814#33206814 reference.

Charm: Plain Text https://github.com/charmbracelet/bubbletea/blob/master/tutorials/mouse.md#full-mouse-mode-without-scrolling docs.


## 2-Level Title

### 3-Level Title

#### 4-Level Title

##### 5-Level Title

###### 6-Level Title

## Unicode Showcase

| Category | Characters | Width |
|:---------|:-----------|:-----|
| Emoji | ✅ ❌ 🚀 ⭐ ⚠️ ✨ 🎉 😀 🔥 | Mixed |
| CJK | 你好世界 畫龍點睛 | 2 ea |
| Japanese | こんにちは コンニチハ | 2 ea |
| Korean | 안녕하세요 | 2 ea |
| Fullwidth | ＡＢＣＸＹＺ １２３ | 2 ea |
| Accented | àáâãäåæçèéêë éèêë çæœ | 1 ea |
| Math | ∑ ∫ ∞ π λ α β γ | 1 ea |
| Greek | αβγδεζηθικλμνξοπρςστυφχψω | 1 ea |
| Currency | ¥ € £ $ ¢ ₩ ₹ ₽ ₿ | Mixed |
| Symbols | ★ ♥ ♦ ♠ ♣ → ← ⇒ § ¶ © ® | Mixed |
| Dingbats | ✓ ✗ ✘ ✝ ✞ ✟ | 1 ea |
| Suits | ♠ ♥ ♦ ♣ | 1 ea |
| Arrows | ← ↑ → ↓ ↔ ⇒ ⇔ | 1 ea |
| Operators | ± ≠ ≈ ≤ ≥ ≡ | 1 ea |
| Brackets | 【「『』」】 | 2 ea |

## Text Formatting

Markdown supports **bold** and *italic* text, as well as ***bold italic***.
You can also use `inline code` for short snippets, or ~~strikethrough~~ for crossed-out text.
Standard paragraphs are the most common element in any document.

## Mixed Content

This paragraph has **bold**, *italic*, `code`, and a [link](https://example.com) all in one sentence.
Here is another one with ***all three*** styles combined and some `inline code` sprinkled throughout.

## Lists

Here are some unordered items:

- Alpha
- Beta
- Gamma
  - Delta
  - Epsilon

And ordered lists:

1. First
2. Second
3. Third

### Nested Lists

1. Item one
   - Sub-item A
   - Sub-item B
     1. Deep item 1
     2. Deep item 2
2. Item two
   - Sub-item C


## Tables

| Name  | Type   | Default | Description       |
| ----- | ------ | ------- | ----------------- |
| Theme | string | "dark"  | Color theme       |
| Width | int    | 120     | Word wrap width   |
| Debug | bool   | false   | Enable debug mode |

## Blockquotes

> This is a blockquote.
> It can span multiple lines.
>
> And even contain nested elements.

## Thematic Break

---

## Task List

- [x] Learned markdown syntax
- [x] Wrote benchmark document
- [ ] Run performance tests
- [ ] Analyze results

## Code Blocks

Here is a Go example:

```go
package main

import "fmt"

func main() {
  fmt.Println("Hello, World!")
}
```

And a JavaScript example:

```javascript
function greet(name) {
  return `Hello, ${name}!`;
}
console.log(greet("World"));
```

```python
def fibonacci(n):
    a, b = 0, 1
    for _ in range(n):
        yield a
        a, b = b, a + b
```

```rust
fn main() {
    let msg = "Hello, Rust!";
    println!("{}", msg);
}
```

### Code Without Language

```
GET    /users              # 获取用户列表
GET    /users/123          # 获取单个用户
POST   /users              # 创建用户
PUT    /users/123          # 更新用户（完整）
PATCH  /users/123          # 更新用户（部分）
DELETE /users/123          # 删除用户
```

## Summary

This benchmark file exercises all major markdown features including headings, paragraphs, text formatting, lists, links, code blocks, tables, blockquotes, and task lists. The mixed structure provides a realistic workload for performance comparison between Gruff and Glamour.
