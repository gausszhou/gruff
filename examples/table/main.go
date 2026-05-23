package main

import (
	"fmt"
	"log"

	"github.com/gausszhou/gruff"
)

func main() {
	md := "" +
		"# Unicode Table Rendering\n\n" +
		"Covers CJK, emoji, fullwidth, dingbats, math, currency, and accented characters.\n\n" +

		"## Emoji & Dingbats\n\n" +
		"| Char | Name | Codepoint |\n" +
		"|:----:|:-----|:----------|\n" +
		"| ✅ | Check Mark | U+2705 |\n" +
		"| ❌ | Cross Mark | U+274C |\n" +
		"| 🚀 | Rocket | U+1F680 |\n" +
		"| ⚠️ | Warning | U+26A0 |\n" +
		"| ✨ | Sparkles | U+2728 |\n" +
		"| ⭐ | Star | U+2B50 |\n" +
		"| 🎉 | Party Popper | U+1F389 |\n" +
		"| ✅ ✅ ✅ | Triple Check | — |\n" +
		"| 🚀🚀🚀 | Triple Rocket | — |\n\n" +

		"## CJK Characters\n\n" +
		"| Type | Characters | Description |\n" +
		"|:----:|:-----------|:------------|\n" +
		"| Chinese | 你好世界 | Simplified Chinese |\n" +
		"| Chinese | 畫龍點睛 | Traditional Chinese |\n" +
		"| Japanese | こんにちは | Hiragana + Kanji |\n" +
		"| Japanese | コンニチハ | Katakana |\n" +
		"| Korean | 안녕하세요 | Hangul |\n" +
		"| CJK mixed | 山﨑さん | CJK + Hiragana |\n" +
		"| CJK + emoji | 你好 ✅ | Mixed script |\n" +
		"| Fullwidth | ＡＢＣＸＹＺ | Fullwidth Latin |\n" +
		"| Fullwidth | １２３４５６ | Fullwidth digits |\n" +
		"| CJK punct | 【】「」『』 | CJK brackets |\n\n" +

		"## Latin & Accented\n\n" +
		"| Type | Characters | Notes |\n" +
		"|:----:|:-----------|:------|\n" +
		"| Basic | abcdefghijklmnopqrstuvwxyz | 26 letters |\n" +
		"| Uppercase | ABCDEFGHIJKLMNOPQRSTUVWXYZ | 26 letters |\n" +
		"| Accented | àáâãäåæçèéêëìíîï | Latin-1 |\n" +
		"| Accented | ñòóôõöøùúûüýþÿ | Extended |\n" +
		"| Accented | ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏ | Uppercase |\n" +
		"| Phonetic | ðŋœʃʒʔθðŋ | IPA symbols |\n" +
		"| German | ÄÖÜßäöü | Umlauts + eszett |\n" +
		"| French | çéèêëîïôûùœæ | French accent mix |\n| Spanish | ¿ñóüíéá¡ | Spanish punctuation |\n" +
		"| Polish | ąćęłńóśźż | Polish diacritics |\n\n" +

		"## Symbols & Math\n\n" +
		"| Char | Name | Category |\n" +
		"|:----:|:-----|:---------|\n" +
		"| ∑ | Summation | Math |\n" +
		"| ∫ | Integral | Math |\n" +
		"| ∞ | Infinity | Math |\n" +
		"| π | Pi | Greek |\n" +
		"| λ | Lambda | Greek |\n" +
		"| ✓ | Check | Symbol |\n" +
		"| ✗ | Cross | Symbol |\n" +
		"| ★ | Star filled | Symbol |\n" +
		"| ♥ | Heart | Suit |\n" +
		"| ♦ | Diamond | Suit |\n" +
		"| ♠ | Spade | Suit |\n" +
		"| ♣ | Club | Suit |\n" +
		"| → | Arrow right | Arrow |\n" +
		"| ← | Arrow left | Arrow |\n" +
		"| ⇒ | Double arrow | Arrow |\n" +
		"| § | Section | Punctuation |\n" +
		"| ¶ | Pilcrow | Punctuation |\n" +
		"| © | Copyright | Legal |\n" +
		"| ® | Registered | Legal |\n" +
		"| ™ | Trademark | Legal |\n" +
		"| ± | Plus-minus | Math operator |\n" +
		"| ≠ | Not equal | Math operator |\n" +
		"| ≈ | Approximately | Math operator |\n\n" +

		"## Currency\n\n" +
		"| Symbol | Name | Code |\n" +
		"|:------:|:-----|:-----|\n" +
		"| ¥ | Yen | JPY/CNY |\n" +
		"| € | Euro | EUR |\n" +
		"| £ | Pound | GBP |\n" +
		"| $ | Dollar | USD |\n" +
		"| ¢ | Cent | Cent |\n" +
		"| ₩ | Won | KRW |\n" +
		"| ₹ | Rupee | INR |\n" +
		"| ₽ | Ruble | RUB |\n" +
		"| ₿ | Bitcoin | BTC |\n" +
		"| ฿ | Baht | THB |\n\n" +

		"## Mixed Scripts & Edge Cases\n\n" +
		"| Description | Value |\n" +
		"|:------------|:------|\n" +
		"| CJK + emoji + Latin | Hello 你好 🚀 世界！ |\n" +
		"| Math + Greek + symbol | E = mc² ≈ ∑π² ± ∞ |\n" +
		"| Accented + currency | Café €5,99 — ¥3000 |\n" +
		"| Fullwidth + halfwidth | ＡＢ half ＣＤ |\n" +
		"| Punctuation mix | ¿Qué pasa? → ¡Bien! |\n" +
		"| Subscript/Superscript | H₂O and E=mc² |\n" +
		"| Dashes & quotes | — En dash, ≠ — Em dash |\n" +
		"| Bracket nesting | 【「『Test』」】 | Fullwidth brackets |\n" +
		"| Arrows + symbols | ← ↑ → ↓ ↔ ⇒ ⇔ |\n" +
		"| All in one cell | ✅ Hello 你好 ★ ★ ∫ ¥ |\n" +
		"| **Bold CJK** | **你好世界** | Bold formatting |\n" +
		"| `Code with CJK` | `func(你好)` | Code formatting |\n" +
		"| ***All mixed*** | ***✅ 你好 ★ π → €*** | All styles |\n" +
		"| Multiple emoji | 😀😁😂🤣😃😄😅😆😇😈 | Many faces |\n\n" +

		"## Long Text Wrapping\n\n" +
		"| Column A | Column B | Column C |\n" +
		"|:---------|:---------|:---------|\n" +
		"| Short | Fits in one line. | Also short |\n" +
		"| Long plain | This is a very long plain text sentence that absolutely exceeds the maximum column width of forty characters and must wrap to multiple lines inside the cell automatically. | Short |\n" +
		"| Short | Short text | A very long text that wraps across multiple lines and demonstrates the wrapping behavior for the third column as well with more content here. |\n" +
		"| Long CJK | 这是一个非常长的中文句子，它完全超过了四十个字符的最大列宽度，必须自动换行到单元格内的多行才能完整显示。 | 短文本 |\n" +
		"| Long mixed | Hello 你好 🚀 This is a long 🔥 mixed 🌟 sentence with ✅ emoji, CJK 中文, and long text that wraps across multiple lines correctly in the table cell. | End |\n" +
		"| Long **bold** | This is a **very long sentence with bold formatting** that should still wrap correctly at word boundaries while preserving the **bold ANSI styling** across multiple lines of output in the rendered table. | **Bold end** |\n" +
		"| Long `code` | This has `inline code snippets` mixed `into a very long` sentence that `wraps across` multiple lines `while preserving` the code background color. | `code end` |\n" +
		"| Long ★ ★ ★ | ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ | ★ ★ ★ |\n" +
		"| Emoji chain text | ✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅ | ✅✅✅ |\n" +
		"| All three wrap | Every single column in this row has text that is long enough to trigger word wrapping at the maximum column width of the rendered table in the terminal output. | Short again |\n\n" +

		"## Header Wrapping\n\n" +
		"| This is a very long header that definitely exceeds the maximum column width | Short header |\n" +
		"|:----------------------------------------------------------------------------|:-------------|\n" +
		"| Short cell content here | This data cell also has text that wraps across multiple lines to test combined header and body wrapping. |\n" +
		"| Another short row with text | ✅🚀⭐🌟 你好世界 Hello mixed content here for wrapping test |\n" +
		"| **Bold and `code` mixed** | Regular text that wraps under the short header column |\n\n" +

		"## Alignment Test\n\n" +
		"| Left | Center | Right |\n" +
		"|:-----|:------:|------:|\n" +
		"| 你好 | ✅ | ¥3000 |\n" +
		"| Hello | ∑  π  ∫ | ∞\n" +
		"| ★ ★ ★ | ＡＢＣ | １２３ |\n" +
		"| éèê | ♥♦♠♣ | →⇒⇔ |\n"

	out, err := gruff.Render(md)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(out)
}
