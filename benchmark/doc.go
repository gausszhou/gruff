// Package benchmark contains performance comparisons between gruff and glamour.
package benchmark

import (
	"charm.land/glamour/v2/ansi"
	"charm.land/glamour/v2/styles"
)

func GlamourMinimalStyle() ansi.StyleConfig {
	cfg := styles.DarkStyleConfig
	return cfg
}

func GlamourStandardStyle() ansi.StyleConfig {
	cfg := styles.DarkStyleConfig
	return cfg
}
