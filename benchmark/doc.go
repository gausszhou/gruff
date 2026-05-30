// Package benchmark contains performance comparisons between gruff and glamour.
package benchmark

import (
	"charm.land/glamour/v2/ansi"
	"charm.land/glamour/v2/styles"
)

// GlamourMinimalStyle returns a style config with all non-essential features
// disabled for maximum rendering speed.
func GlamourMinimalStyle() ansi.StyleConfig {
	cfg := styles.DarkStyleConfig
	cfg.Strikethrough = ansi.StylePrimitive{}
	cfg.DefinitionList = ansi.StyleBlock{}
	cfg.DefinitionTerm = ansi.StylePrimitive{}
	cfg.DefinitionDescription = ansi.StylePrimitive{}
	cfg.HTMLBlock = ansi.StyleBlock{}
	cfg.HTMLSpan = ansi.StyleBlock{}

	return cfg
}

func GlamourStandardStyle() ansi.StyleConfig {
	cfg := styles.DarkStyleConfig
	return cfg
}
