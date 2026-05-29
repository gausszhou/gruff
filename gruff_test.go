package gruff

import (
	"strings"
	"testing"
)

func TestRender_Heading(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check []string
	}{
		{
			name:  "h1",
			input: "# Heading 1\n",
			check: []string{"\x1b[1m\x1b[38;2;255;255;135m", "Heading 1\x1b[K", "\x1b[22m\x1b[39m"},
		},
		{
			name:  "h2",
			input: "## Heading 2\n",
			check: []string{"\x1b[1m\x1b[38;2;0;175;255mHeading 2", "\x1b[22m\x1b[39m"},
		},
		{
			name:  "h6",
			input: "###### Heading 6\n",
			check: []string{"\x1b[38;2;0;175;95mHeading 6", "\x1b[39m"},
		},
		{
			name:  "heading with inline",
			input: "# **Bold** heading\n",
			check: []string{"\x1b[1m\x1b[38;2;255;255;135m", "Bold\x1b[22m heading\x1b[K", "\x1b[22m\x1b[39m"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			for _, c := range tt.check {
				if !strings.Contains(got, c) {
					t.Errorf("output missing %q\n got: %q", c, got)
				}
			}
		})
	}
}

func TestRender_BoldItalic(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check []string
	}{
		{
			name:  "bold",
			input: "**bold**\n",
			check: []string{"\x1b[1mbold\x1b[22m"},
		},
		{
			name:  "italic",
			input: "*italic*\n",
			check: []string{"\x1b[3mitalic\x1b[23m"},
		},
		{
			name:  "bold italic",
			input: "***both***\n",
			check: []string{"\x1b[3m\x1b[1mboth\x1b[22m\x1b[23m"},
		},
		{
			name:  "nested bold in italic",
			input: "*italic and **bold** inside*\n",
			check: []string{"\x1b[3mitalic and \x1b[1mbold\x1b[22m inside\x1b[23m"},
		},
		{
			name:  "mixed inline paragraph",
			input: "plain **bold** and *italic*.\n",
			check: []string{"plain \x1b[1mbold\x1b[22m and \x1b[3mitalic\x1b[23m."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			for _, c := range tt.check {
				if !strings.Contains(got, c) {
					t.Errorf("output missing %q\n got: %q", c, got)
				}
			}
		})
	}
}

func TestRender_InlineCode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check []string
	}{
		{
			name:  "inline code",
			input: "Use `code` here\n",
			check: []string{"\x1b[38;2;80;250;123mcode\x1b[39m"},
		},
		{
			name:  "code with bold",
			input: "**bold and `code`**\n",
			check: []string{"\x1b[1mbold and \x1b[38;2;80;250;123mcode\x1b[39m"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			for _, c := range tt.check {
				if !strings.Contains(got, c) {
					t.Errorf("output missing %q\n got: %q", c, got)
				}
			}
		})
	}
}

func TestRender_Link(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check []string
	}{
		{
			name:  "basic link",
			input: "[Gruff](https://example.com)\n",
			check: []string{"\x1b[4m\x1b[38;2;92;156;245mGruff\x1b[24m\x1b[39m \x1b[38;2;128;128;128m(https://example.com)\x1b[39m"},
		},
		{
			name:  "link with bold text",
			input: "[**bold**](https://example.com)\n",
			check: []string{"\x1b[4m\x1b[38;2;92;156;245m\x1b[1mbold\x1b[22m\x1b[24m\x1b[39m \x1b[38;2;128;128;128m(https://example.com)\x1b[39m"},
		},
		{
			name:  "link in paragraph",
			input: "click [here](https://example.com) now\n",
			check: []string{"\x1b[4m\x1b[38;2;92;156;245mhere\x1b[24m\x1b[39m \x1b[38;2;128;128;128m(https://example.com)\x1b[39m now"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			for _, c := range tt.check {
				if !strings.Contains(got, c) {
					t.Errorf("output missing %q\n got: %q", c, got)
				}
			}
		})
	}
}

func TestRender_CodeBlock(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check []string
	}{
		{
			name:  "fenced code block",
			input: "```\ncode\n```\n",
			check: []string{"\x1b[38;2;80;250;123m  code\x1b[K", "\x1b[39m"},
		},
		{
			name:  "fenced code with language",
			input: "```go\nvar x = 1\n```\n",
			check: []string{"\x1b[38;2;128;128;128m  go\x1b[K", "\x1b[38;2;80;250;123m  var x = 1\x1b[K"},
		},
		{
			name:  "indented code block",
			input: "    indented\n",
			check: []string{"\x1b[38;2;80;250;123m  indented\x1b[K"},
		},
		{
			name:  "multi-line fenced code",
			input: "```\nline1\nline2\n```\n",
			check: []string{"\x1b[38;2;80;250;123m  line1\x1b[K", "\x1b[38;2;80;250;123m  line2\x1b[K"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			for _, c := range tt.check {
				if !strings.Contains(got, c) {
					t.Errorf("output missing %q\n got: %q", c, got)
				}
			}
		})
	}
}

func TestRender_List(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "unordered",
			input: "- item 1\n- item 2\n",
			want:  "    \x1b[K\n    • item 1  \x1b[K\n    • item 2  \x1b[K\n    \x1b[K\n  ",
		},
		{
			name:  "ordered",
			input: "1. first\n2. second\n",
			want:  "    \x1b[K\n    1. first  \x1b[K\n    2. second  \x1b[K\n    \x1b[K\n  ",
		},
		{
			name:  "nested unordered",
			input: "- Alpha\n- Beta\n  - Delta\n  - Epsilon\n",
			want:  "    \x1b[K\n    • Alpha  \x1b[K\n    • Beta  \x1b[K\n      • Delta  \x1b[K\n      • Epsilon  \x1b[K\n    \x1b[K\n  ",
		},
		{
			name:  "nested ordered",
			input: "1. first\n2. second\n   1. deep\n   2. deeper\n",
			want:  "    \x1b[K\n    1. first  \x1b[K\n    2. second  \x1b[K\n      1. deep  \x1b[K\n      2. deeper  \x1b[K\n    \x1b[K\n  ",
		},
		{
			name:  "task list",
			input: "- [x] done\n- [ ] todo\n",
			want:  "    \x1b[K\n  \x1b[38;2;80;250;123m[✓]\x1b[39m done  \x1b[K\n  \x1b[38;2;128;128;128m[ ]\x1b[39m todo  \x1b[K\n    \x1b[K\n  ",
		},
		{
			name:  "nested task list",
			input: "- [x] checked\n- [ ] unchecked\n  - [x] nested checked\n  - [ ] nested unchecked\n",
			want:  "    \x1b[K\n  \x1b[38;2;80;250;123m[✓]\x1b[39m checked  \x1b[K\n  \x1b[38;2;128;128;128m[ ]\x1b[39m unchecked  \x1b[K\n  \x1b[38;2;80;250;123m[✓]\x1b[39m nested checked  \x1b[K\n  \x1b[38;2;128;128;128m[ ]\x1b[39m nested unchecked  \x1b[K\n    \x1b[K\n  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("Render() =\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}

func TestRender_Blockquote(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple",
			input: "> A quote\n",
			want:  "    \x1b[K\n  \x1b[38;2;128;128;128m│ \x1b[39mA quote  \x1b[K\n  ",
		},
		{
			name:  "multi paragraph",
			input: "> First\n>\n> Second\n",
			want:  "    \x1b[K\n  \x1b[38;2;128;128;128m│ \x1b[39mFirst  \x1b[K\n  \x1b[38;2;128;128;128m│ \x1b[39m  \x1b[K\n  \x1b[38;2;128;128;128m│ \x1b[39mSecond  \x1b[K\n  ",
		},
		{
			name:  "with inline",
			input: "> **bold** and *italic*\n",
			want:  "    \x1b[K\n  \x1b[38;2;128;128;128m│ \x1b[39m\x1b[1mbold\x1b[22m and \x1b[3mitalic\x1b[23m  \x1b[K\n  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("Render() =\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}

func TestRender_Table(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "simple table",
			input: "| H1 | H2 |\n| --- | --- |\n| A | B |\n",
			want:  []string{"\u2502", "H1", "H2", "A", "B"},
		},
		{
			name:  "table with alignment",
			input: "| Left | Center | Right |\n|:-----|:------:|------:|\n| a    | b      | c     |\n",
			want:  []string{"Left", "Center", "Right", "a", "b", "c"},
		},
		{
			name:  "table with inline",
			input: "| Col1 | Col2 |\n|------|------|\n| `code` | **bold** |\n",
			want:  []string{"\x1b[38;2;80;250;123m", "\x1b[1m", "code", "bold"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("output missing %q", w)
				}
			}
		})
	}
}

func TestRender_Mixed(t *testing.T) {
	input := "# Title\n\nThis is **bold** and *italic* and `code`.\n\nA [link](https://example.com) here.\n\n- list with **bold**\n- list with `code`\n\n1. first\n2. second\n\n> A quote\n\n| A | B |\n|---|---|\n| 1 | 2 |\n"

	got, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}

	checks := []struct {
		name string
		fn   func(string) bool
	}{
		{"contains Title", func(s string) bool { return strings.Contains(s, "Title") }},
		{"contains bold ANSI", func(s string) bool { return strings.Contains(s, "\x1b[1m") }},
		{"contains italic ANSI", func(s string) bool { return strings.Contains(s, "\x1b[3m") }},
		{"contains code ANSI", func(s string) bool { return strings.Contains(s, "\x1b[38;2;80;250;123m") }},
		{"contains link underline", func(s string) bool { return strings.Contains(s, "\x1b[4m") }},
		{"contains link URL", func(s string) bool { return strings.Contains(s, "example.com") }},
		{"contains bullet", func(s string) bool { return strings.Contains(s, "•") }},
		{"contains ordered num", func(s string) bool { return strings.Contains(s, "1.") }},
		{"contains table separator", func(s string) bool { return strings.Contains(s, "\u2502") }},
		{"contains blockquote pipe", func(s string) bool { return strings.Contains(s, "│ ") }},
	}

	for _, c := range checks {
		if !c.fn(got) {
			t.Errorf("output should %s", c.name)
		}
	}
}

func TestRender_Empty(t *testing.T) {
	got, err := Render("")
	if err != nil {
		t.Fatal(err)
	}
	if got != "    \x1b[K\n  " {
		t.Errorf("Render() = %q, want bg + padding", got)
	}
}

func TestRender_PlainText(t *testing.T) {
	got, err := Render("Hello, world!\n")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "Hello, world!") {
		t.Errorf("Render() missing 'Hello, world!', got %q", got)
	}
}

func TestOptions(t *testing.T) {
	t.Run("light theme", func(t *testing.T) {
		got, err := Render("# Hello", WithLight())
		if err != nil {
			t.Fatal(err)
		}
		if got == "" {
			t.Error("expected non-empty output")
		}
	})

	t.Run("word wrap", func(t *testing.T) {
		input := "This is a long paragraph that should be wrapped at the specified width."
		got, err := Render(input, WithWordWrap(20))
		if err != nil {
			t.Fatal(err)
		}
		got = stripANSI(got)
		lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
		for _, line := range lines {
			if len(line) > 22 {
				t.Errorf("line length %d exceeds 22: %q", len(line), line)
			}
		}
	})
}

func TestRenderBytes(t *testing.T) {
	out, err := RenderBytes([]byte("# Hello"))
	if err != nil {
		t.Fatal(err)
	}
	if len(out) == 0 {
		t.Error("expected non-empty output")
	}
}

