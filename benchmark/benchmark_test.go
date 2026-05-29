package benchmark

import (
	"os"
	"testing"

	"charm.land/glamour/v2"
	"charm.land/glamour/v2/ansi"
	"charm.land/glamour/v2/styles"
	"github.com/gausszhou/gruff"
)

// gruffMinimalStyle returns a style config based on "dark" but stripped to
// match gruff's supported elements. Gruff handles:
//   Document, Paragraph, Heading (H1-H6), List, ListItem,
//   Text, String, Emphasis (bold/italic), CodeSpan, Link, Image,
//   FencedCodeBlock, CodeBlock, ThematicBreak, Table
//
// Everything else (BlockQuote, Strikethrough, TaskCheckBox, HTMLBlock,
// RawHTML, DefinitionList, etc.) is neutralized.
func gruffMinimalStyle() ansi.StyleConfig {
	cfg := styles.DarkStyleConfig

	// Neutralize elements gruff doesn't support
	cfg.CodeBlock.Chroma = nil
	cfg.CodeBlock.Theme = ""
	cfg.BlockQuote.IndentToken = nil
	cfg.BlockQuote.Indent = nil
	cfg.BlockQuote.Margin = nil
	cfg.Strikethrough = ansi.StylePrimitive{}
	cfg.Task.Ticked = ""
	cfg.Task.Unticked = ""
	cfg.DefinitionList = ansi.StyleBlock{}
	cfg.DefinitionTerm = ansi.StylePrimitive{}
	cfg.DefinitionDescription = ansi.StylePrimitive{}
	cfg.HTMLBlock = ansi.StyleBlock{}
	cfg.HTMLSpan = ansi.StyleBlock{}

	// Strip ANSI decorations to reduce allocations
	cfg.HorizontalRule.Format = "\n"
	cfg.Item.BlockPrefix = " "
	cfg.Enumeration.BlockPrefix = ""
	cfg.Code.Prefix = ""
	cfg.Code.Suffix = ""
	cfg.ImageText.Format = ""
	cfg.Image = ansi.StylePrimitive{}
	cfg.Link = ansi.StylePrimitive{}
	cfg.LinkText = ansi.StylePrimitive{}

	return cfg
}

func benchGruff(b *testing.B, file string) {
	source, err := os.ReadFile("../" + file)
	if err != nil {
		b.Fatal(err)
	}
	input := string(source)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gruff.Render(input)
	}
}

func BenchmarkGruff(b *testing.B) { benchGruff(b, "testdata/benchmark.md") }

func benchGlamour(b *testing.B, file string) {
	source, err := os.ReadFile("../" + file)
	if err != nil {
		b.Fatal(err)
	}
	input := string(source)

	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(gruffMinimalStyle()),
		glamour.WithWordWrap(0),
		glamour.WithTableWrap(false),
		glamour.WithInlineTableLinks(true),
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Render(input)
	}
}

func BenchmarkGlamour(b *testing.B) { benchGlamour(b, "testdata/benchmark.md") }
