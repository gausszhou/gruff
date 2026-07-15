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

| Category | Characters | Width | Reference |
|:---------|:-----------|:-----|:----------|
| Emoji | ✅ ❌ 🚀 ⭐ ⚠️ ✨ 🎉 😀 🔥 | Mixed | [Emoji List](https://unicode.org/emoji/charts/full-emoji-list.html) |
| CJK | 你好世界 畫龍點睛 | 2 ea | [CJK Unified](https://en.wikipedia.org/wiki/CJK_Unified_Ideographs) |
| Japanese | こんにちは コンニチハ | 2 ea | [Hiragana](https://en.wikipedia.org/wiki/Hiragana) / [Katakana](https://en.wikipedia.org/wiki/Katakana) |
| Korean | 안녕하세요 | 2 ea | [Hangul](https://en.wikipedia.org/wiki/Hangul) |
| Fullwidth | ＡＢＣＸＹＺ １２３ | 2 ea | [Halfwidth and Fullwidth](https://en.wikipedia.org/wiki/Halfwidth_and_fullwidth_forms) |
| Accented | àáâãäåæçèéêë éèêë çæœ | 1 ea | [Latin Extended](https://en.wikipedia.org/wiki/Latin_Extended-A) |
| Math | ∑ ∫ ∞ π λ α β γ | 1 ea | [Mathematical Operators](https://en.wikipedia.org/wiki/Mathematical_operators_and_symbols_in_Unicode) |
| Greek | αβγδεζηθικλμνξοπρςστυφχψω | 1 ea | [Greek and Coptic](https://en.wikipedia.org/wiki/Greek_and_Coptic) |
| Currency | ¥ € £ $ ¢ ₩ ₹ ₽ ₿ | Mixed | [Currency Symbols](https://en.wikipedia.org/wiki/Currency_symbol_(typography)) |
| Symbols | ★ ♥ ♦ ♠ ♣ → ← ⇒ § ¶ © ® | Mixed | [Miscellaneous Symbols](https://en.wikipedia.org/wiki/Miscellaneous_Symbols) |
| Dingbats | ✓ ✗ ✘ ✝ ✞ ✟ | 1 ea | [Dingbat](https://en.wikipedia.org/wiki/Dingbat) |
| Suits | ♠ ♥ ♦ ♣ | 1 ea | [Playing Cards](https://en.wikipedia.org/wiki/Playing_cards_in_Unicode) |
| Arrows | ← ↑ → ↓ ↔ ⇒ ⇔ | 1 ea | [Arrows block](https://en.wikipedia.org/wiki/Arrows_(Unicode_block)) |
| Operators | ± ≠ ≈ ≤ ≥ ≡ | 1 ea | https://en.wikipedia.org/wiki/Mathematical_operators_and_symbols_in_Unicode |
| Brackets | 【「『』」】 | 2 ea | https://en.wikipedia.org/wiki/Bracket#East_Asian_brackets |

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

| Name  | Type   | Default | Description       | More Info |
| ----- | ------ | ------- | ----------------- | --------- |
| Theme | string | "dark"  | Color theme       | [ANSI colors](https://en.wikipedia.org/wiki/ANSI_escape_code#Colors) |
| Width | int    | 120     | Word wrap width   | [Terminal width](https://en.wikipedia.org/wiki/Page_width) |
| Debug | bool   | false   | Enable debug mode | [Debugging](https://en.wikipedia.org/wiki/Debugging) |

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
