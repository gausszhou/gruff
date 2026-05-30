# Gruff vs Glamour Benchmark

This document contains a balanced mix of markdown elements for performance testing.

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

## Links

Visit [GitHub](https://github.com) for more information.

Check out [Gruff](https://github.com/gausszhou/gruff) the markdown renderer.

You can also reference [Go](https://go.dev) documentation.

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

## Summary

This benchmark file exercises all major markdown features including headings, paragraphs, text formatting, lists, links, code blocks, tables, blockquotes, and task lists. The mixed structure provides a realistic workload for performance comparison between Gruff and Glamour.
