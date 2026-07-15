# Gruff vs Glamour Benchmark

This document contains a balanced mix of markdown elements for performance testing.

## Links

Visit [GitHub](https://github.com) for more information.

Check out [Gruff](https://github.com/gausszhou/gruff) the markdown renderer.

You can also reference [Go](https://go.dev) documentation.

Example of a long URL: [golang.org/x/term](https://pkg.go.dev/golang.org/x/term#Terminal.ReadPassword) wraps nicely.

Another lengthy link: [goldmark extension](https://github.com/yuin/goldmark/tree/main/extension#table-extension) for table support.

Autolink: <https://github.com/charmbracelet/glamour/tree/master/styles/gallery>

Wikipedia: [Unicode](https://en.wikipedia.org/wiki/Unicode#Standardization_and_development_process) consortium page.

Stack Overflow: [ANSI escape codes](https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences/33206814#33206814) reference.

RFC: [URI Generic Syntax](https://datatracker.ietf.org/doc/html/rfc3986#section-3.2.2).

Go docs: [text/template](https://pkg.go.dev/text/template#pkg-functions) functions reference.

NPM: [chalk](https://www.npmjs.com/package/chalk/v/5.3.0#256-and-truecolor-color-support) color library.

Docker Hub: [golang image](https://hub.docker.com/_/golang/tags?page=1&name=1.23-bookworm).

Go issue: [proposal: spec: generic type aliases](https://github.com/golang/go/issues/46477#issuecomment-1990934563).

Charm: [bubbletea mouse](https://github.com/charmbracelet/bubbletea/blob/master/tutorials/mouse.md#full-mouse-mode-without-scrolling) docs.

Rust lang: [Option](https://doc.rust-lang.org/std/option/enum.Option.html#method.map_or_else) map_or_else method.

Kubernetes: [pod-lifecycle](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-states) container states.

CPAN: [Moose::Manual::Attributes](https://metacpan.org/pod/Moose::Manual::Attributes#Delegation-with-type-coercion-AND-currying) docs.

LWN: [kernel development](https://lwn.net/Articles/1012580/) latest articles.

ArXiv: [Attention is All You Need](https://arxiv.org/abs/1706.03762) paper.

Wikipedia: [SQL injection](https://en.wikipedia.org/wiki/SQL_injection#Technical_implementations) technical details.

Autolink: <https://caniuse.com/?search=css-container-queries>

Autolink: <https://www.w3.org/TR/css-color-4/#lab-to-lch-conversion-drafts>

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
