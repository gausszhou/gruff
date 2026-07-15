# Gruff Markdown Renderer

A lightweight, high-performance **markdown** renderer for the terminal.
Supports CJK (你好世界), emoji (✅ 🚀), math (∑π∞), and more.

## Text Formatting

- *Italic* and **bold** text
- `Inline code` with background
- ***Bold italic*** combined
- ~~Strikethrough~~ text
- **Bold with *italic* nested**
- *Italic with **bold** nested*

## Lists

- Item with **bold**
- Item with *italic*
- Item with `code`
- Item with ***both***
- Item with ✅ emoji
- Item with 中文

1. First ordered item with **bold text**
2. Second ordered item with *italic text*
3. Third ordered item with `inline code`
4. Fourth with 中文
5. Fifth with 🚀

## Links

- Plain link: [Gruff](https://github.com/gausszhou/gruff)
- Bold link: [**bold text**](https://example.com)
- Code link: [`code`](https://example.com)

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

## Table with Alignment

| Left | Center | Right |
|:-----|:------:|------:|
| 你好 | ✅ | ¥3000 |
| Hello | ∑ π ∫ | ∞ |
| ★ ★ ★ | ＡＢＣ | １２３ |
| éèê | ♥♦♠♣ | →⇒⇔ |
| **Bold** | *Italic* | `Code` |
| ✅ ✅ | 🚀 🚀 | ⭐ ⭐ |

## Thematic Break

---

Above is a horizontal rule.

## Combined Edge Cases

- **CJK with bold**: **你好世界**
- *Accented italic*: *àáâãäå*
- `Code with math`: `E = mc² ≈ ∑π²`
- ***Mixed styles***: ***✅ 你好 ★ π → €***
- Link with unicode: [你好世界](https://example.com)

## Code Blocks

Fenced code block:

```
func hello() {
    fmt.Println("Hello, World!")
}
```

Indented code block:

    line1
    line2

Visit [Gruff on GitHub](https://github.com/gausszhou/gruff) for more information.


## Unicode Showcase

| Category | Characters | Width | Reference |
|:---------|:-----------|:-----|:----------|
| Emoji | ✅ ❌ 🚀 ⭐ ⚠️ ✨ 🎉 😀 🔥 | Mixed | [Emoji List](https://unicode.org/emoji/charts/full-emoji-list.html) |
| CJK | 你好世界 畫龍點睛 | 2 ea | [CJK Unified](https://en.wikipedia.org/wiki/CJK_Unified_Ideographs) |
| Dingbats | ✓ ✗ ✘ ✝ ✞ ✟ | 1 ea | [Dingbat](https://en.wikipedia.org/wiki/Dingbat) |
| Suits | ♠ ♥ ♦ ♣ | 1 ea | <https://en.wikipedia.org/wiki/Playing_cards_in_Unicode> |
| Arrows | ← ↑ → ↓ ↔ ⇒ ⇔ | 1 ea | <https://en.wikipedia.org/wiki/Arrows_(Unicode_block)> |
| Operators | ± ≠ ≈ ≤ ≥ ≡ | 1 ea | https://en.wikipedia.org/wiki/Mathematical_operators_and_symbols_in_Unicode |
| Brackets | 【「『』」】 | 2 ea | https://en.wikipedia.org/wiki/Bracket#East_Asian_brackets |