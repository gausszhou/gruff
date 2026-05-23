package benchmark

import (
	"os"
	"testing"

	"charm.land/glamour/v2"
	"github.com/gausszhou/gruff"
)

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

func BenchmarkGruff(b *testing.B) { benchGruff(b, "testdata/_data.md") }

func benchGlamour(b *testing.B, file string) {
	source, err := os.ReadFile("../" + file)
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

func BenchmarkGlamour(b *testing.B) { benchGlamour(b, "testdata/_data.md") }
