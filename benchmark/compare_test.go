package benchmark

import (
	"os"
	"strings"
	"testing"

	"charm.land/glamour/v2"
	"github.com/gausszhou/gruff"
)

func stripANSI(s string) string {
	var out []byte
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			for j := i + 2; j < len(s); j++ {
				if s[j] == 'm' {
					i = j
					break
				}
			}
			continue
		}
		out = append(out, s[i])
	}
	return string(out)
}

func TestCompareWrapWidth(t *testing.T) {
	b, err := os.ReadFile("../testdata/_data.md")
	if err != nil {
		t.Fatal(err)
	}
	md := string(b)

	width := 76

	// gruff
	gOut, err := gruff.Render(md, gruff.WithWordWrap(width))
	if err != nil {
		t.Fatal(err)
	}
	gStripped := stripANSI(gOut)
	gLines := strings.Split(strings.TrimRight(gStripped, "\n"), "\n")

	// glamour
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		t.Fatal(err)
	}
	glOut, err := r.Render(md)
	if err != nil {
		t.Fatal(err)
	}
	glStripped := stripANSI(glOut)
	glLines := strings.Split(strings.TrimRight(glStripped, "\n"), "\n")

	// Compare line counts
	t.Logf("gruff lines: %d, glamour lines: %d", len(gLines), len(glLines))

	// Check for lines exceeding width
	gMax, glMax := 0, 0
	for _, line := range gLines {
		if len(line) > gMax {
			gMax = len(line)
		}
		if len(line) > width+2 {
			t.Errorf("GRUFF line exceeds width %d: len=%d [%s]", width, len(line), line[:40])
		}
	}
	for _, line := range glLines {
		if len(line) > glMax {
			glMax = len(line)
		}
		if len(line) > width+2 {
			t.Errorf("GLAMOUR line exceeds width %d: len=%d [%s]", width, len(line), line[:40])
		}
	}
	t.Logf("gruff max line: %d, glamour max line: %d", gMax, glMax)

	// Find first 10 lines where lengths differ
	diffCount := 0
	for i := 0; i < len(gLines) && i < len(glLines) && diffCount < 10; i++ {
		if len(gLines[i]) != len(glLines[i]) {
			t.Logf("DIFF line %d: gruff len=%d [%s]", i, len(gLines[i]), gLines[i])
			t.Logf("       glamour len=%d [%s]", len(glLines[i]), glLines[i])
			diffCount++
		}
	}
}
