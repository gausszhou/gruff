package main

import (
	"fmt"
	"log"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/gausszhou/gruff"
	"golang.org/x/term"
)

var docStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(0, 1).
	Width(72)

var sectionStyle = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Color("63")).
	Padding(0, 1).
	Width(72)

var noteStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("248")).
	Italic(true).
	Padding(0, 2)

func termWidth() int {
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		return w
	}
	return 0
}

func main() {
	title := "# gruff + lipgloss\n\n" +
		"Render markdown with **gruff** and wrap it with *lipgloss* styling.\n"
	out, err := gruff.Render(title, gruff.WithWordWrap(termWidth()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(docStyle.Render(out))

	md := "## Table\n\n" +
		"| Feature | Status |\n" +
		"|---------|--------|\n" +
		"| Headings | ✅ |\n" +
		"| Bold | ✅ |\n" +
		"| Italic | ✅ |\n" +
		"| Tables | ✅ |\n"
	out, err = gruff.Render(md, gruff.WithWordWrap(termWidth()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sectionStyle.Render(out))

	fmt.Println(noteStyle.Render("lipgloss provides borders, padding, colors, and width control around gruff's ANSI output."))
}
