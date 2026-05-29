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

func benchGlamourStandard(b *testing.B, file string) {
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

func benchGlamourMinimal(b *testing.B, file string) {
	source, err := os.ReadFile("../" + file)
	if err != nil {
		b.Fatal(err)
	}
	input := string(source)

	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(GruffMinimalStyle()),
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

func BenchmarkGruff(b *testing.B) { benchGruff(b, "testdata/benchmark.md") }

func BenchmarkGlamourMinimal(b *testing.B) { benchGlamourMinimal(b, "testdata/benchmark.md") }

func BenchmarkGlamourStandard(b *testing.B) { benchGlamourStandard(b, "testdata/benchmark.md") }

func BenchmarkLargeGruff(b *testing.B) { benchGruff(b, "testdata/_data.md") }

func BenchmarkLargeGlamourMinimal(b *testing.B) { benchGlamourMinimal(b, "testdata/_data.md") }

func BenchmarkLargeGlamourStandard(b *testing.B) { benchGlamourStandard(b, "testdata/_data.md") }
