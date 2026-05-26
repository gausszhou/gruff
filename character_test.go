package gruff

import (
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
)

// ============================================================
// displayWidth — character display width calculation
// ============================================================

func TestDisplayWidth_ASCII(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"", 0},
		{" ", 1},
		{"a", 1},
		{"Hello", 5},
		{"Hello, World!", 13},
	}
	for _, tt := range tests {
		got := displayWidth(tt.input)
		if got != tt.want {
			t.Errorf("displayWidth(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestDisplayWidth_CJK(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"中", 2},
		{"文", 2},
		{"中文", 4},
		{"你好世界", 8},
		{"Hello世界", 9},
		{"コンニチハ", 10},             // katakana
		{"한글", 4},                    // hangul
		{"中文English混合", 15},          // 中(2)+文(2)+E(1)+n(1)+g(1)+l(1)+i(1)+s(1)+h(1)+混(2)+合(2)
	}
	for _, tt := range tests {
		got := displayWidth(tt.input)
		if got != tt.want {
			t.Errorf("displayWidth(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestDisplayWidth_Fullwidth(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"Ａ", 2},  // fullwidth A
		{"，", 2},  // fullwidth comma
		{"（）", 4}, // fullwidth parens
	}
	for _, tt := range tests {
		got := displayWidth(tt.input)
		if got != tt.want {
			t.Errorf("displayWidth(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestDisplayWidth_Emoji(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"thumbs up", "👍", 2},
		{"party popper", "🎉", 2},
		{"bubbles", "🫧", 2},
		{"smile", "😀", 2},
		{"heart", "❤", 1}, // without VS16 — ambiguous, width 1
		{"three emoji", "👍🎉🫧", 6},
		{"emoji with text", "Hi👍there", 9}, // H(1)+i(1)+👍(2)+t(1)+h(1)+e(1)+r(1)+e(1)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := displayWidth(tt.input)
			if got != tt.want {
				t.Errorf("displayWidth(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// Known limitation: go-runewidth is code-point-level, not grapheme-cluster-level.
// VS16 (U+FE0F) does not increase width beyond the base codepoint width.
// Run with `RUNEWIDTH_EASTASIAN=1` to treat ambiguous chars as wide.
func TestDisplayWidth_VS16(t *testing.T) {
	// ❤ (U+2764) is East Asian Ambiguous → width 1 with default settings
	// Adding VS16 (U+FE0F, zero-width) should force emoji presentation (width 2)
	// but go-runewidth doesn't handle this at the grapheme cluster level.
	base := "❤"       // U+2764
	withVS16 := "❤️"   // U+2764 + U+FE0F
	want := 1          // known limitation: go-runewidth doesn't promote to width 2
	gotBase := displayWidth(base)
	gotVS16 := displayWidth(withVS16)
	if gotBase != want {
		t.Errorf("displayWidth(❤) = %d, want %d", gotBase, want)
	}
	// VS16 itself is zero-width, so total should be same as base
	if gotVS16 != gotBase {
		t.Errorf("displayWidth(❤️) = %d, want %d (VS16 should be zero-width in go-runewidth)", gotVS16, gotBase)
	}
}

// Known limitation: go-runewidth doesn't handle regional indicator pairs
// as single graphemes. go-runewidth v0.0.23 with default settings treats
// an RI pair as having combined width 1.
func TestDisplayWidth_FlagSequence(t *testing.T) {
	china := "🇨🇳" // U+1F1E8 + U+1F1F3, a regional indicator pair
	got := displayWidth(china)
	want := 1 // known limitation: go-runewidth reports 1 for RI pairs
	if got != want {
		t.Errorf("displayWidth(🇨🇳) = %d, want %d (known go-runewidth limitation)", got, want)
	}
}

func TestDisplayWidth_Combining(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"precomposed é", "é", 1},
		{"decomposed e+combining", "e\u0301", 1},
		{"a+combining grave", "a\u0300", 1},
		{"cafe with combining", "caf\u00e9", 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := displayWidth(tt.input)
			if got != tt.want {
				t.Errorf("displayWidth(%[1]q %[1]s) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestDisplayWidth_ZeroWidth(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"combining acute", "\u0301", 0},
		{"combining grave", "\u0300", 0},
		{"ZWJ", "\u200D", 0},
		{"ZWNJ", "\u200C", 0},
		{"VS16", "\uFE0F", 0},
		{"VS15", "\uFE0E", 0},
		{"soft hyphen", "\u00AD", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := displayWidth(tt.input)
			if got != tt.want {
				t.Errorf("displayWidth(U+%04X) = %d, want %d", []rune(tt.input)[0], got, tt.want)
			}
		})
	}
}

func TestDisplayWidth_Halfwidth(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"ｶ", 1}, // halfwidth katakana
		{"ﾟ", 1},
		{"Helloｶ", 6},
	}
	for _, tt := range tests {
		got := displayWidth(tt.input)
		if got != tt.want {
			t.Errorf("displayWidth(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

// ============================================================
// stripANSI — ANSI escape sequence removal with Unicode
// ============================================================

func TestStripANSI_Unicode(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"\x1b[1m中文\x1b[22m", "中文"},
		{"\x1b[38;2;80;134;90m世界\x1b[39m", "世界"},
		{"\x1b[1m👍\x1b[22m", "👍"},
		{"彩\x1b[31m虹\x1b[39m色", "彩虹色"},
	}
	for _, tt := range tests {
		got := stripANSI(tt.input)
		if got != tt.want {
			t.Errorf("stripANSI(%q) = %q, want %q", tt.input, got, tt.want)
		}
		// After stripping, displayWidth should reflect visual width
		w := displayWidth(got)
		if w == 0 && got != "" {
			t.Errorf("displayWidth(stripANSI(%q)) = 0, want > 0", tt.input)
		}
	}
}

// ============================================================
// wrapCellLines — table cell word wrapping with Unicode
// ============================================================

func TestWrapCellLines_ASCII(t *testing.T) {
	tests := []struct {
		content string
		width   int
		want    []string
	}{
		{"Hello", 10, []string{"Hello"}},
		{"Hello World", 10, []string{"Hello", "World"}},
		{"a b c d e", 3, []string{"a b", "c d", "e"}},
	}
	for _, tt := range tests {
		got := wrapCellLines(tt.content, tt.width)
		if len(got) != len(tt.want) {
			t.Errorf("wrapCellLines(%q, %d) = %q, want %q", tt.content, tt.width, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("wrapCellLines(%q, %d)[%d] = %q, want %q", tt.content, tt.width, i, got[i], tt.want[i])
			}
			if runewidth.StringWidth(got[i]) > tt.width {
				t.Errorf("wrapCellLines(%q, %d)[%d] visLen=%d > width=%d",
					tt.content, tt.width, i, runewidth.StringWidth(got[i]), tt.width)
			}
		}
	}
}

func TestWrapCellLines_CJK(t *testing.T) {
	tests := []struct {
		name    string
		content string
		width   int
	}{
		{"exact fit 8", "你好世界", 8},
		{"exact 6", "你好世", 6},
		{"overflow 6→3+3", "你好世界", 6},
		{"single char break", "你好世界", 4},
		{"mixed CJK ASCII", "Hello世界", 8},
		{"all CJK 5 chars", "天地玄黄宇", 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := wrapCellLines(tt.content, tt.width)
			if len(lines) == 0 {
				t.Fatal("expected at least one line")
			}
			for i, line := range lines {
				vlen := runewidth.StringWidth(line)
				if vlen > tt.width {
					t.Errorf("line %d visLen=%d > width=%d: %q", i, vlen, tt.width, line)
				}
			}
		})
	}
}

func TestWrapCellLines_ZeroWidth(t *testing.T) {
	// After the fix: wordVisLen += RuneWidth(r) instead of wordVisLen++
	// Zero-width chars should not contribute to word width.
	tests := []struct {
		name    string
		content string
		width   int
	}{
		{"combining acute", "a\u0301bc", 3},
		{"combining grave x2", "a\u0300bc", 3},
		{"precomposed cafe", "caf\u00e9", 4},
		{"ZWJ not counted", "a\u200Db", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := wrapCellLines(tt.content, tt.width)
			if len(lines) == 0 {
				t.Fatal("expected at least one line")
			}
			for i, line := range lines {
				vlen := runewidth.StringWidth(line)
				if vlen > tt.width {
					t.Errorf("line %d visLen=%d > width=%d: %q", i, vlen, tt.width, line)
				}
			}
		})
	}
}

func TestWrapCellLines_Emoji(t *testing.T) {
	tests := []struct {
		name    string
		content string
		width   int
	}{
		{"three emoji fit", "👍🎉🫧", 6},
		{"emoji overflow", "👍🎉🫧", 4},
		{"emoji mixed", "Hi👍there", 8},
		{"emoji + cjk", "中👍文", 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := wrapCellLines(tt.content, tt.width)
			if len(lines) == 0 {
				t.Fatal("expected at least one line")
			}
			for i, line := range lines {
				vlen := runewidth.StringWidth(line)
				if vlen > tt.width {
					t.Errorf("line %d visLen=%d > width=%d: %q", i, vlen, tt.width, line)
				}
			}
		})
	}
}

func TestWrapCellLines_NoOverflow(t *testing.T) {
	// Ensure no line exceeds the specified width for various content types.
	tests := []struct {
		name    string
		content string
		width   int
	}{
		{"long ascii", "a b c d e f g h i j k l m n o p", 8},
		{"cjk long", "天地玄黄宇宙洪荒日月盈昃辰宿列张", 12},
		{"mixed long", "Hello World 你好 世界 Test 123", 10},
		{"cjk single words", "中文 测试 宽度 计算 是否 正确", 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := wrapCellLines(tt.content, tt.width)
			for i, line := range lines {
				vlen := runewidth.StringWidth(line)
				if vlen > tt.width {
					t.Errorf("line %d visLen=%d > width=%d: %q", i, vlen, tt.width, line)
				}
			}
		})
	}
}

// ============================================================
// wrapText — document-level word wrapping with Unicode
// ============================================================

func TestWrapText_CJK(t *testing.T) {
	tests := []struct {
		name  string
		input string
		width int
	}{
		{"cjk with spaces", "你好 世界 这是一个 测试", 10},
		{"mixed cjk ascii", "Hello World 中文 测试", 12},
		{"emoji wrap", "👍 🎉 🫧 emoji test here", 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := wrapText(tt.input, tt.width)
			for _, line := range strings.Split(wrapped, "\n") {
				w := runewidth.StringWidth(line)
				if w > tt.width {
					t.Errorf("visLen=%d > width=%d: %q", w, tt.width, line)
				}
			}
		})
	}
}

// ============================================================
// Rendering integration — CJK / emoji in Markdown structures
// ============================================================

func TestRender_CJKHeading(t *testing.T) {
	input := "# 你好世界\n"
	got, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	clean := stripANSI(got)
	if !strings.Contains(clean, "你好世界") {
		t.Errorf("output missing CJK heading text, got %q", clean)
	}
}

func TestRender_CJKInline(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"bold CJK", "**中文加粗**\n"},
		{"italic CJK", "*中文斜体*\n"},
		{"bold italic CJK", "***中文粗斜体***\n"},
		{"inline code CJK", "`中文代码`\n"},
		{"link with CJK", "[中文链接](https://example.com)\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			clean := stripANSI(got)
			if !strings.Contains(clean, "中文") {
				t.Errorf("output missing CJK text in %s, got %q", tt.name, clean)
			}
		})
	}
}

func TestRender_CJKList(t *testing.T) {
	input := "- 项目一\n- 项目二\n- 项目三\n"
	got, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range []string{"项目一", "项目二", "项目三"} {
		if !strings.Contains(got, item) {
			t.Errorf("output missing list item %q", item)
		}
	}
}

func TestRender_CJKTable(t *testing.T) {
	input := "| 中文 | English |\n|------|---------|\n| 你好 | Hello |\n| 世界 | World |\n"
	got, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"中文", "English", "你好", "Hello", "世界", "World"} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q", want)
		}
	}
	// Verify table uses box-drawing separators
	if !strings.Contains(got, "│") {
		t.Error("CJK table missing column separators")
	}
}

func TestRender_EmojiTable(t *testing.T) {
	input := "| Emoji | Desc |\n|-------|------|\n| 👍 | Thumbs up |\n| 🎉 | Party |\n"
	got, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"👍", "🎉", "Thumbs up", "Party"} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestRender_CJKCodeBlock(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"fenced CJK code", "```\n中文代码\n```\n"},
		{"indented CJK code", "    中文代码\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(got, "中文代码") {
				t.Errorf("output missing CJK code, got %q", got)
			}
		})
	}
}

func TestRender_CJKParagraph(t *testing.T) {
	input := "这是一个中文段落。它包含**加粗**和*斜体*，以及`代码`。\n"
	got, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "中文段落") {
		t.Errorf("output missing CJK paragraph text")
	}
}

func TestRender_EmojiInline(t *testing.T) {
	input := "# Emoji Title\n\nParagraph with 👍 emoji.\n"
	got, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "👍") {
		t.Errorf("output missing emoji, got %q", got)
	}
}

// ============================================================
// Table column alignment with CJK — visual width verification
// ============================================================

// ============================================================
// Table column width expansion with CJK content
// ============================================================

func TestRender_CJKTableWidthExpansion(t *testing.T) {
	// "李四四" (width 6) should expand column from 4 to 6
	input := "| 姓名 | 年龄 |\n|------|------|\n| 李四四 | 35 |\n"
	got, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	sepLine := stripANSI(lines[1])
	dashCount := strings.Count(sepLine, "─")
	if dashCount != 14 { // seg(6)+seg(4) = 8+6 dashes
		t.Errorf("separator has %d dashes, want 14 (col widths [6,4])", dashCount)
	}
	if !strings.Contains(got, "李四四") {
		t.Error("CJK content missing from table")
	}
}

// ============================================================
// Table column alignment with CJK — visual width verification
// ============================================================

func TestRender_CJKTableContentPreserved(t *testing.T) {
	input := "| 姓名 | 年龄 |\n|------|------|\n| 张三 | 28 |\n| 李四四 | 35 |\n"
	got, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"姓名", "年龄", "张三", "28", "李四四", "35"} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

// ============================================================
// Mixed CJK + emoji + ASCII in one document
// ============================================================

func TestRender_MixedCJKEmoji(t *testing.T) {
	input := `# 综合测试

这是一个**混合**文档，包含：
- 中文列表项
- English items
- 👍 Emoji 项目

| 类型 | 示例 |
|------|------|
| 中文 | 你好 |
| Emoji | 🎉 |
| 混合 | Hello中文👍 |
`
	got, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"综合测试", "混合", "中文列表项", "English", "👍", "🎉", "你好", "Hello中文👍"} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q", want)
		}
	}
}
