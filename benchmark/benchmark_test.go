package benchmark

import (
	"os"
	"testing"

	"github.com/charmbracelet/glamour"
	"github.com/gausszhou/gruff"
)

func BenchmarkGruff(b *testing.B) {
	source, err := os.ReadFile("../testdata/sample.md")
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

func BenchmarkGlamour(b *testing.B) {
	source, err := os.ReadFile("../testdata/sample.md")
	if err != nil {
		b.Fatal(err)
	}
	input := string(source)

	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
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
