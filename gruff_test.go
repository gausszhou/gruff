package gruff

import (
	"strings"
	"testing"
)

func TestRender_Heading(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "h1",
			input: "# Heading 1\n",
			want:  "\x1b[1m\x1b[4m\x1b[38;5;15mHeading 1\x1b[22m\x1b[24m\x1b[39m\n\n",
		},
		{
			name:  "h2",
			input: "## Heading 2\n",
			want:  "\x1b[1m\x1b[38;5;11mHeading 2\x1b[22m\x1b[39m\n\n",
		},
		{
			name:  "h6",
			input: "###### Heading 6\n",
			want:  "\x1b[38;5;8mHeading 6\x1b[39m\n\n",
		},
		{
			name:  "heading with inline",
			input: "# **Bold** heading\n",
			want:  "\x1b[1m\x1b[4m\x1b[38;5;15m\x1b[1mBold\x1b[22m heading\x1b[22m\x1b[24m\x1b[39m\n\n",
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

func TestRender_BoldItalic(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "bold",
			input: "**bold**\n",
			want:  "\x1b[1mbold\x1b[22m\n\n",
		},
		{
			name:  "italic",
			input: "*italic*\n",
			want:  "\x1b[3mitalic\x1b[23m\n\n",
		},
		{
			name:  "bold italic",
			input: "***both***\n",
			want:  "\x1b[3m\x1b[1mboth\x1b[22m\x1b[23m\n\n",
		},
		{
			name:  "nested bold in italic",
			input: "*italic and **bold** inside*\n",
			want:  "\x1b[3mitalic and \x1b[1mbold\x1b[22m inside\x1b[23m\n\n",
		},
		{
			name:  "mixed inline paragraph",
			input: "plain **bold** and *italic*.\n",
			want:  "plain \x1b[1mbold\x1b[22m and \x1b[3mitalic\x1b[23m.\n\n",
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

func TestRender_InlineCode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "inline code",
			input: "Use `code` here\n",
			want:  "Use \x1b[38;5;15m\x1b[48;5;236mcode\x1b[39m\x1b[49m here\n\n",
		},
		{
			name:  "code with bold",
			input: "**bold and `code`**\n",
			want:  "\x1b[1mbold and \x1b[38;5;15m\x1b[48;5;236mcode\x1b[39m\x1b[49m\x1b[22m\n\n",
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

func TestRender_List(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "unordered",
			input: "- item 1\n- item 2\n",
			want:  "  \x1b[38;5;11m• \x1b[39mitem 1\n  \x1b[38;5;11m• \x1b[39mitem 2\n\n",
		},
		{
			name:  "ordered",
			input: "1. first\n2. second\n",
			want:  "  \x1b[38;5;11m1. \x1b[39mfirst\n  \x1b[38;5;11m2. \x1b[39msecond\n\n",
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

func TestRender_Mixed(t *testing.T) {
	input := "# Title\n\nThis is **bold** and *italic* and `code`.\n\n- list with **bold**\n- list with `code`\n\n1. first\n2. second\n"

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
		{"contains code ANSI", func(s string) bool { return strings.Contains(s, "\x1b[48;5;236m") }},
		{"contains bullet", func(s string) bool { return strings.Contains(s, "•") }},
		{"contains ordered num", func(s string) bool { return strings.Contains(s, "1.") }},
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
	if got != "" {
		t.Errorf("Render() = %q, want empty string", got)
	}
}

func TestRender_PlainText(t *testing.T) {
	got, err := Render("Hello, world!\n")
	if err != nil {
		t.Fatal(err)
	}
	want := "Hello, world!\n\n"
	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
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

func BenchmarkRender(b *testing.B) {
	input := `# Gruff Markdown Renderer

A lightweight, high-performance **markdown** renderer for the terminal.

## Features

- *Italic* and **bold** text
- ` + "`" + `Inline code` + "`" + ` support
- Unordered lists:
  - Item 1
  - Item 2
- Ordered lists:
  1. First
  2. Second

## Code Example

` + "```" + `
func hello() {
    fmt.Println("Hello, World!")
}
` + "```" + `

## Mixed Content

This paragraph has **bold**, *italic*, ` + "`" + `code` + "`" + `, and ***bold italic*** all in one line.

And another paragraph to make it interesting.

---

The end.
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Render(input)
	}
}
