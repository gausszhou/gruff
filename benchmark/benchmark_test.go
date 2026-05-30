package benchmark

import (
	"os"
	"strings"
	"testing"

	"charm.land/glamour/v2"
	"github.com/gausszhou/gruff/gruff"
)

func benchGruff(b *testing.B, file string) {
	source, err := os.ReadFile("../" + file)
	if err != nil {
		b.Fatal(err)
	}
	input := strings.Repeat(string(source), 100)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gruff.Render(input)
	}
}

func benchGlamourMinimal(b *testing.B, file string) {
	source, err := os.ReadFile("../" + file)
	if err != nil {
		b.Fatal(err)
	}
	input := strings.Repeat(string(source), 100)

	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(GlamourMinimalStyle()),
		glamour.WithWordWrap(0),
		glamour.WithTableWrap(false),
		glamour.WithInlineTableLinks(true),
		glamour.WithChromaFormatter("noop"),
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

func benchGlamourStandard(b *testing.B, file string) {
	source, err := os.ReadFile("../" + file)
	if err != nil {
		b.Fatal(err)
	}
	input := strings.Repeat(string(source), 100)

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

func BenchmarkGruff(b *testing.B) { benchGruff(b, "testdata/benchmark.md") }

func BenchmarkGlamourMinimal(b *testing.B) { benchGlamourMinimal(b, "testdata/benchmark.md") }

func BenchmarkGlamourStandard(b *testing.B) { benchGlamourStandard(b, "testdata/benchmark.md") }
