package main

import (
	"log"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gausszhou/gruff"
)

var viewportStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#202020")).
	Width(119)

var noBgTheme = gruff.Theme{
	H1:       gruff.Style{Bold: true, Fg: "#ffffff"},
	H2:       gruff.Style{Bold: true, Fg: "#ffff00"},
	H3:       gruff.Style{Bold: true, Fg: "#00ff00"},
	H4:       gruff.Style{Bold: true, Fg: "#00ffff"},
	H5:       gruff.Style{Bold: true, Fg: "#808080"},
	H6:       gruff.Style{Fg: "#808080"},
	Strong:   gruff.Style{Bold: true},
	Em:       gruff.Style{Italic: true},
	Code:     gruff.Style{Fg: "#50865a"},
	Link:     gruff.Style{Underline: true, Fg: "#5c9cf5"},
	LinkURL:  gruff.Style{Fg: "#808080"},
	Bullet:   gruff.Style{Fg: "#ffff00"},
	Numbered: gruff.Style{Fg: "#ffff00"},
}

type model struct {
	viewport            viewport.Model
	ready               bool
	termWidth, termHeight int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		if !m.ready {
			m.viewport = viewport.New(119, 36)
			var md strings.Builder
			md.WriteString("# Unicode Character Table\n\n")
			md.WriteString("All printable Unicode characters from U+0020 to U+10FFFF, rendered by **gruff**.\n\n")
			md.WriteString("---\n\n")
			for r := rune(0x20); r <= 0x10FFFF; r++ {
				if r >= 0x7F && r <= 0x9F {
					continue
				}
				if r >= 0xD800 && r <= 0xDFFF {
					continue
				}
				if !unicode.IsPrint(r) {
					continue
				}
				md.WriteRune(r)
				md.WriteRune(' ')
			}
			out, err := gruff.Render(md.String(),
				gruff.WithWordWrap(118),
				func(o *gruff.Options) { o.Theme = noBgTheme },
			)
			if err != nil {
				log.Fatal(err)
			}
			m.viewport.SetContent(viewportStyle.Render(out))
			m.ready = true
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

var (
	scrollTrack = lipgloss.NewStyle().Background(lipgloss.Color("#2a2a2a"))
	scrollThumb = lipgloss.NewStyle().Background(lipgloss.Color("#666666"))
)

func scrollbar(vp viewport.Model, height int) string {
	total := vp.TotalLineCount()
	if total <= height {
		return ""
	}
	thumb := int(vp.ScrollPercent() * float64(height-1))
	var lines []string
	for i := 0; i < height; i++ {
		if i == thumb {
			lines = append(lines, scrollThumb.Render(" "))
		} else {
			lines = append(lines, scrollTrack.Render(" "))
		}
	}
	return strings.Join(lines, "\n")
}

func (m model) View() string {
	if !m.ready {
		return "\n  Loading..."
	}
	vpView := m.viewport.View()
	if bar := scrollbar(m.viewport, m.viewport.Height); bar != "" {
		vpView = lipgloss.JoinHorizontal(lipgloss.Top, vpView, bar)
	}
	if m.termWidth > 0 {
		bg := lipgloss.NewStyle().
			Background(lipgloss.Color("#141414")).
			Width(m.termWidth).
			Height(m.termHeight)
		return bg.Render(vpView)
	}
	return vpView
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
