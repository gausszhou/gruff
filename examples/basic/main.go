package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gausszhou/gruff"
)

var sampleMD = "# Gruff Markdown Renderer\n\n" +
	"A lightweight, high-performance **markdown** renderer for the terminal.\n" +
	"Supports CJK (你好世界), emoji (✅ 🚀), math (∑π∞), and more.\n\n" +
	"## Text Formatting\n\n" +
	"- *Italic* and **bold** text\n" +
	"- `Inline code` with background\n" +
	"- ***Bold italic*** combined\n" +
	"- ~~Strikethrough~~ text\n" +
	"- **Bold with *italic* nested**\n" +
	"- *Italic with **bold** nested*\n\n" +
	"## Lists\n\n" +
	"- Item with **bold**\n" +
	"- Item with *italic*\n" +
	"- Item with `code`\n" +
	"- Item with ***both***\n" +
	"- Item with ✅ emoji\n" +
	"- Item with 中文\n\n" +
	"1. First ordered item with **bold text**\n" +
	"2. Second ordered item with *italic text*\n" +
	"3. Third ordered item with `inline code`\n" +
	"4. Fourth with 中文\n" +
	"5. Fifth with 🚀\n\n" +
	"## Links\n\n" +
	"- Plain link: [Gruff](https://github.com/gausszhou/gruff)\n" +
	"- Bold link: [**bold text**](https://example.com)\n" +
	"- Code link: [`code`](https://example.com)\n\n" +
	"## Unicode Showcase\n\n" +
	"| Category | Characters | Width |\n" +
	"|:---------|:-----------|:-----|\n" +
	"| Emoji | ✅ ❌ 🚀 ⭐ ⚠️ ✨ 🎉 😀 🔥 | Mixed |\n" +
	"| CJK | 你好世界 畫龍點睛 | 2 ea |\n" +
	"| Japanese | こんにちは コンニチハ | 2 ea |\n" +
	"| Korean | 안녕하세요 | 2 ea |\n" +
	"| Fullwidth | ＡＢＣＸＹＺ １２３ | 2 ea |\n" +
	"| Accented | àáâãäåæçèéêë éèêë çæœ | 1 ea |\n" +
	"| Math | ∑ ∫ ∞ π λ α β γ | 1 ea |\n" +
	"| Greek | αβγδεζηθικλμνξοπρςστυφχψω | 1 ea |\n" +
	"| Currency | ¥ € £ $ ¢ ₩ ₹ ₽ ₿ | Mixed |\n" +
	"| Symbols | ★ ♥ ♦ ♠ ♣ → ← ⇒ § ¶ © ® | Mixed |\n" +
	"| Dingbats | ✓ ✗ ✘ ✝ ✞ ✟ | 1 ea |\n" +
	"| Suits | ♠ ♥ ♦ ♣ | 1 ea |\n" +
	"| Arrows | ← ↑ → ↓ ↔ ⇒ ⇔ | 1 ea |\n" +
	"| Operators | ± ≠ ≈ ≤ ≥ ≡ | 1 ea |\n" +
	"| Brackets | 【「『』」】 | 2 ea |\n\n" +
	"## Table with Alignment\n\n" +
	"| Left | Center | Right |\n" +
	"|:-----|:------:|------:|\n" +
	"| 你好 | ✅ | ¥3000 |\n" +
	"| Hello | ∑ π ∫ | ∞ |\n" +
	"| ★ ★ ★ | ＡＢＣ | １２３ |\n" +
	"| éèê | ♥♦♠♣ | →⇒⇔ |\n" +
	"| **Bold** | *Italic* | `Code` |\n" +
	"| ✅ ✅ | 🚀 🚀 | ⭐ ⭐ |\n\n" +
	"## Thematic Break\n\n" +
	"---\n\n" +
	"Above is a horizontal rule.\n\n" +
	"## Combined Edge Cases\n\n" +
	"- **CJK with bold**: **你好世界**\n" +
	"- *Accented italic*: *àáâãäå*\n" +
	"- `Code with math`: `E = mc² ≈ ∑π²`\n" +
	"- ***Mixed styles***: ***✅ 你好 ★ π → €***\n" +
	"- Link with unicode: [你好世界](https://example.com)\n\n" +
	"Visit [Gruff on GitHub](https://github.com/gausszhou/gruff) for more information.\n"

func main() {
	light := flag.Bool("light", false, "use light theme")
	wrap := flag.Int("wrap", 0, "word wrap width (0 = no wrap)")
	flag.Parse()

	var opts []gruff.Option
	if *light {
		opts = append(opts, gruff.WithLight())
	}
	if *wrap > 0 {
		opts = append(opts, gruff.WithWordWrap(*wrap))
	}

	out, err := gruff.Render(sampleMD, opts...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(out)
}
