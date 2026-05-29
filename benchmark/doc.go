// Package benchmark contains performance comparisons between gruff and glamour.
package benchmark

import (
	"charm.land/glamour/v2/ansi"
	"charm.land/glamour/v2/styles"
)

// GruffMinimalStyle returns a style config based on "dark" but stripped to
// match gruff's supported elements. Gruff handles:
//
//	Document, Paragraph, Heading (H1-H6), List, ListItem,
//	Text, String, Emphasis (bold/italic), CodeSpan, Link, Image,
//	FencedCodeBlock, CodeBlock, ThematicBreak, Table
//
// Everything else (BlockQuote, Strikethrough, TaskCheckBox, HTMLBlock,
// RawHTML, DefinitionList, etc.) is neutralized.
func GruffMinimalStyle() ansi.StyleConfig {
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
