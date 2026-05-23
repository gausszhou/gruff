package gruff

import (
	_ "embed"
	"testing"

	"github.com/charmbracelet/glamour"
)

//go:embed testdata/sample.md
var sampleMD string

func BenchmarkGruff(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Render(sampleMD)
	}
}

func BenchmarkGlamour(b *testing.B) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r.Render(sampleMD)
	}
}
