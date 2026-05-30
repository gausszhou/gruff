# Unicode Table Rendering

Covers CJK, emoji, fullwidth, dingbats, math, currency, and accented characters.

## Emoji & Dingbats

| Char | Name | Codepoint |
|:----:|:-----|:----------|
| ✅ | Check Mark | U+2705 |
| ❌ | Cross Mark | U+274C |
| 🚀 | Rocket | U+1F680 |
| ⚠️ | Warning | U+26A0 |
| ✨ | Sparkles | U+2728 |
| ⭐ | Star | U+2B50 |
| 🎉 | Party Popper | U+1F389 |
| ✅ ✅ ✅ | Triple Check | — |
| 🚀🚀🚀 | Triple Rocket | — |

## CJK Characters

| Type | Characters | Description |
|:----:|:-----------|:------------|
| Chinese | 你好世界 | Simplified Chinese |
| Chinese | 畫龍點睛 | Traditional Chinese |
| Japanese | こんにちは | Hiragana + Kanji |
| Japanese | コンニチハ | Katakana |
| Korean | 안녕하세요 | Hangul |
| CJK mixed | 山﨑さん | CJK + Hiragana |
| CJK + emoji | 你好 ✅ | Mixed script |
| Fullwidth | ＡＢＣＸＹＺ | Fullwidth Latin |
| Fullwidth | １２３４５６ | Fullwidth digits |
| CJK punct | 【】「」『』 | CJK brackets |

## Latin & Accented

| Type | Characters | Notes |
|:----:|:-----------|:------|
| Basic | abcdefghijklmnopqrstuvwxyz | 26 letters |
| Uppercase | ABCDEFGHIJKLMNOPQRSTUVWXYZ | 26 letters |
| Accented | àáâãäåæçèéêëìíîï | Latin-1 |
| Accented | ñòóôõöøùúûüýþÿ | Extended |
| Accented | ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏ | Uppercase |
| Phonetic | ðŋœʃʒʔθðŋ | IPA symbols |
| German | ÄÖÜßäöü | Umlauts + eszett |
| French | çéèêëîïôûùœæ | French accent mix |
| Spanish | ¿ñóüíéá¡ | Spanish punctuation |
| Polish | ąćęłńóśźż | Polish diacritics |

## Symbols & Math

| Char | Name | Category |
|:----:|:-----|:---------|
| ∑ | Summation | Math |
| ∫ | Integral | Math |
| ∞ | Infinity | Math |
| π | Pi | Greek |
| λ | Lambda | Greek |
| ✓ | Check | Symbol |
| ✗ | Cross | Symbol |
| ★ | Star filled | Symbol |
| ♥ | Heart | Suit |
| ♦ | Diamond | Suit |
| ♠ | Spade | Suit |
| ♣ | Club | Suit |
| → | Arrow right | Arrow |
| ← | Arrow left | Arrow |
| ⇒ | Double arrow | Arrow |
| § | Section | Punctuation |
| ¶ | Pilcrow | Punctuation |
| © | Copyright | Legal |
| ® | Registered | Legal |
| ™ | Trademark | Legal |
| ± | Plus-minus | Math operator |
| ≠ | Not equal | Math operator |
| ≈ | Approximately | Math operator |

## Currency

| Symbol | Name | Code |
|:------:|:-----|:-----|
| ¥ | Yen | JPY/CNY |
| € | Euro | EUR |
| £ | Pound | GBP |
| $ | Dollar | USD |
| ¢ | Cent | Cent |
| ₩ | Won | KRW |
| ₹ | Rupee | INR |
| ₽ | Ruble | RUB |
| ₿ | Bitcoin | BTC |
| ฿ | Baht | THB |

## Mixed Scripts & Edge Cases

| Description | Value |
|:------------|:------|
| CJK + emoji + Latin | Hello 你好 🚀 世界！ |
| Math + Greek + symbol | E = mc² ≈ ∑π² ± ∞ |
| Accented + currency | Café €5,99 — ¥3000 |
| Fullwidth + halfwidth | ＡＢ half ＣＤ |
| Punctuation mix | ¿Qué pasa? → ¡Bien! |
| Subscript/Superscript | H₂O and E=mc² |
| Dashes & quotes | — En dash, ≠ — Em dash |
| Bracket nesting | 【「『Test』」】 | Fullwidth brackets |
| Arrows + symbols | ← ↑ → ↓ ↔ ⇒ ⇔ |
| All in one cell | ✅ Hello 你好 ★ ★ ∫ ¥ |
| **Bold CJK** | **你好世界** | Bold formatting |
| `Code with CJK` | `func(你好)` | Code formatting |
| ***All mixed*** | ***✅ 你好 ★ π → €*** | All styles |
| Multiple emoji | 😀😁😂🤣😃😄😅😆😇😈 | Many faces |

## Long Text Wrapping

| Column A | Column B | Column C |
|:---------|:---------|:---------|
| Short | Fits in one line. | Also short |
| Long plain | This is a very long plain text sentence that absolutely exceeds the maximum column width of forty characters and must wrap to multiple lines inside the cell automatically. | Short |
| Short | Short text | A very long text that wraps across multiple lines and demonstrates the wrapping behavior for the third column as well with more content here. |
| Long CJK | 这是一个非常长的中文句子，它完全超过了四十个字符的最大列宽度，必须自动换行到单元格内的多行才能完整显示。 | 短文本 |
| Long mixed | Hello 你好 🚀 This is a long 🔥 mixed 🌟 sentence with ✅ emoji, CJK 中文, and long text that wraps across multiple lines correctly in the table cell. | End |
| Long **bold** | This is a **very long sentence with bold formatting** that should still wrap correctly at word boundaries while preserving the **bold ANSI styling** across multiple lines of output in the rendered table. | **Bold end** |
| Long `code` | This has `inline code snippets` mixed `into a very long` sentence that `wraps across` multiple lines `while preserving` the code background color. | `code end` |
| Long ★ ★ ★ | ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ ★ | ★ ★ ★ |
| Emoji chain text | ✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅ | ✅✅✅ |
| All three wrap | Every single column in this row has text that is long enough to trigger word wrapping at the maximum column width of the rendered table in the terminal output. | Short again |

## Header Wrapping

| This is a very long header that definitely exceeds the maximum column width | Short header |
|:----------------------------------------------------------------------------|:-------------|
| Short cell content here | This data cell also has text that wraps across multiple lines to test combined header and body wrapping. |
| Another short row with text | ✅🚀⭐🌟 你好世界 Hello mixed content here for wrapping test |
| **Bold and `code` mixed** | Regular text that wraps under the short header column |

## Alignment Test

| Left | Center | Right |
|:-----|:------:|------:|
| 你好 | ✅ | ¥3000 |
| Hello | ∑  π  ∫ | ∞
| ★ ★ ★ | ＡＢＣ | １２３ |
| éèê | ♥♦♠♣ | →⇒⇔ |
