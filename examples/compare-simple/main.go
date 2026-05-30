package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
	"time"

	"charm.land/glamour/v2"
	"charm.land/lipgloss/v2"
	"golang.org/x/term"

	"github.com/gausszhou/gruff/benchmark"
)

func main() {
	b, err := os.ReadFile("testdata/benchmark.md")
	if err != nil {
		log.Fatal(err)
	}
	md := strings.Repeat(string(b), 100)

	termW := 120
	if w, _, err := term.GetSize(1); err == nil {
		termW = w
	}
	halfW := (termW - 4) / 2

	r1, err := glamour.NewTermRenderer(
		glamour.WithStyles(benchmark.GlamourMinimalStyle()),
		glamour.WithWordWrap(halfW),
	)
	if err != nil {
		log.Fatal(err)
	}

	cleaned := benchmark.CleanInput(md)

	t1 := time.Now()
	out1, err := r1.Render(cleaned)
	if err != nil {
		log.Fatal(err)
	}
	dur1 := time.Since(t1)

	r2, err := glamour.NewTermRenderer(
		glamour.WithStyles(benchmark.GlamourStandardStyle()),
		glamour.WithWordWrap(halfW),
	)
	if err != nil {
		log.Fatal(err)
	}

	t2 := time.Now()
	out2, err := r2.Render(md)
	if err != nil {
		log.Fatal(err)
	}
	dur2 := time.Since(t2)

	info := func(title string, dur time.Duration, fg color.Color) string {
		return lipgloss.NewStyle().Background(fg).Foreground(lipgloss.Color("#fff")).Padding(0, 1).Width(halfW - 2).Render(
			fmt.Sprintf("%s    %s", title, dur.Round(time.Microsecond)),
		)
	}

	leftPane := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("#059669")).Render(
		out1 + "\n" + info("glamour minimal", dur1, lipgloss.Color("#059669")),
	)
	rightPane := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("#7c3aed")).Render(
		out2 + "\n" + info("glamour standard", dur2, lipgloss.Color("#7c3aed")),
	)

	fmt.Print(lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane))
}
