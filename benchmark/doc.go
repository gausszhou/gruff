package benchmark

import (
	"strings"

	"charm.land/glamour/v2/ansi"
	"charm.land/glamour/v2/styles"
)

// CleanInput strips \r from input to prevent corruption in Chroma=nil fallback
// (BaseElement leaves \r intact, causing visual glitches on terminals).
func CleanInput(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// GlamourMinimalStyle returns a style config aligned with gruff's visual output.
// Chroma = nil skips the chroma pipeline for maximum speed. Call CleanInput
// on the markdown before rendering to prevent \r corruption in code blocks.
func GlamourMinimalStyle() ansi.StyleConfig {
	cfg := styles.DarkStyleConfig

	cfg.Strikethrough = ansi.StylePrimitive{}
	cfg.Image = ansi.StylePrimitive{}
	cfg.ImageText = ansi.StylePrimitive{}
	cfg.DefinitionList = ansi.StyleBlock{}
	cfg.DefinitionTerm = ansi.StylePrimitive{}
	cfg.DefinitionDescription = ansi.StylePrimitive{}
	cfg.HTMLBlock = ansi.StyleBlock{}
	cfg.HTMLSpan = ansi.StyleBlock{}

	cfg.CodeBlock.Chroma = nil

	return cfg
}

func GlamourStandardStyle() ansi.StyleConfig {
	cfg := styles.DarkStyleConfig
	return cfg
}
