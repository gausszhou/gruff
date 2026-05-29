package benchmark

import (
	"strings"
	"testing"

	"github.com/gausszhou/gruff"
)

func generateMarkdown(n int) string {
	var sb strings.Builder
	for sb.Len() < n {
		sb.WriteString("# Section Title\n\n")
		sb.WriteString("This is a paragraph with **bold** and *italic* text.\n\n")
		sb.WriteString("- list item one\n")
		sb.WriteString("- list item **two**\n")
		sb.WriteString("- list item *three*\n\n")
		sb.WriteString("```\n")
		sb.WriteString("func hello() {\n")
		sb.WriteString("\tfmt.Println(\"world\")\n")
		sb.WriteString("}\n")
		sb.WriteString("```\n\n")
	}
	return sb.String()
}

func generateLargeMarkdown(n int) string {
	var sb strings.Builder
	sb.WriteString("# Million Char Markdown\n\n")
	for sb.Len() < n {
		sb.WriteString("This is a **bold** and *italic* paragraph with `inline code`.\n\n")
		sb.WriteString("- list item with **formatting**\n")
		sb.WriteString("- another *list* item\n\n")
		sb.WriteString("```go\n")
		sb.WriteString("func foo() {\n")
		sb.WriteString("\tfmt.Println(\"hello\")\n")
		sb.WriteString("}\n")
		sb.WriteString("```\n\n")
		sb.WriteString("> A **blockquote** with *style*\n\n")
	}
	return sb.String()
}

func TestRenderMarkdown5MChars(t *testing.T) {
	content := generateLargeMarkdown(5_000_000)

	rendered, err := gruff.Render(content, gruff.WithWordWrap(80))
	if err != nil {
		t.Fatal(err)
	}
	if len(rendered) == 0 {
		t.Fatal("Render returned empty string")
	}
	if len(rendered) < len(content)/2 {
		t.Errorf("rendered length = %d, expected at least %d", len(rendered), len(content)/2)
	}
	if !strings.Contains(rendered, "Million Char Markdown") {
		t.Error("rendered output should contain the heading text")
	}
	if !strings.Contains(rendered, "func foo") {
		t.Error("rendered output should contain code block content")
	}
}

func BenchmarkRenderLargeMarkdown(b *testing.B) {
	content := generateLargeMarkdown(1_000_000)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gruff.Render(content, gruff.WithWordWrap(80))
	}
}

func BenchmarkRenderMediumMarkdown(b *testing.B) {
	content := generateMarkdown(100_000)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gruff.Render(content, gruff.WithWordWrap(80))
	}
}

func BenchmarkRender100x10kMarkdown(b *testing.B) {
	var contents []string
	for i := 0; i < 100; i++ {
		contents = append(contents, generateLargeMarkdown(10_000))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, c := range contents {
			gruff.Render(c, gruff.WithWordWrap(80))
		}
	}
}
