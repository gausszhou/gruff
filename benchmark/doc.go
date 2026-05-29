// Package benchmark contains performance comparisons between gruff and glamour.
package benchmark

import (
	"charm.land/glamour/v2/ansi"
	"charm.land/glamour/v2/styles"
)

func strPtr(s string) *string { return &s }

// GruffMinimalStyle returns a style config based on "dark" with chroma disabled
// (basic colors instead of full syntax highlighting) and elements gruff doesn't
// handle neutralized, while preserving visual quality for supported features.
//
// Gruff handles: Document, Paragraph, Heading (H1-H6), List, ListItem,
// Text, String, Emphasis (bold/italic), CodeSpan, Link, Image,
// FencedCodeBlock, CodeBlock, ThematicBreak, Table.
func GruffMinimalStyle() ansi.StyleConfig {
	cfg := styles.DarkStyleConfig
	cfg.CodeBlock.Theme = ""
	cfg.CodeBlock.Color = strPtr("#50fa7b")
	cfg.CodeBlock.BackgroundColor = strPtr("#1e1e1e")
	cfg.Strikethrough = ansi.StylePrimitive{}
	cfg.DefinitionList = ansi.StyleBlock{}
	cfg.DefinitionTerm = ansi.StylePrimitive{}
	cfg.DefinitionDescription = ansi.StylePrimitive{}
	cfg.HTMLBlock = ansi.StyleBlock{}
	cfg.HTMLSpan = ansi.StyleBlock{}

	return cfg
}
