package main

import (
	"log"
	"strings"
	"unicode"

	"charm.land/glamour/v2"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gausszhou/gruff"
	"github.com/gausszhou/gruff/benchmark"
)

var noBgTheme = gruff.Theme{
	Document: gruff.Style{Padding: 2},
	H1:       gruff.Style{Bold: true, Fg: "#FFFF87", Bg: "#5F5FFF"},
	H2:       gruff.Style{Bold: true, Fg: "#00AFFF"},
	H3:       gruff.Style{Bold: true, Fg: "#00AFFF"},
	H4:       gruff.Style{Bold: true, Fg: "#00AFFF"},
	H5:       gruff.Style{Bold: true, Fg: "#00AFFF"},
	H6:       gruff.Style{Fg: "#00AF5F"},
	Strong:   gruff.Style{Bold: true},
	Em:       gruff.Style{Italic: true},
	Code:     gruff.Style{Fg: "#FF5F5F"},
	Link:     gruff.Style{Underline: true, Fg: "#5c9cf5"},
	LinkURL:  gruff.Style{Fg: "#808080"},
	Bullet:   gruff.Style{Fg: "#ffff00"},
	Numbered: gruff.Style{Fg: "#ffff00"},
}

type focus int

const (
	focusGlamour focus = iota
	focusGruff
)

type model struct {
	glamourView viewport.Model
	gruffView   viewport.Model
	ready       bool
	termWidth   int
	termHeight  int
	focus       focus
}

func generateUnicodeMD(title string) string {
	var md strings.Builder
	md.WriteString("# Unicode Character Table\n\n")
	md.WriteString(title)
	md.WriteString("\n\n---\n\n")
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
	return md.String()
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
		case "tab":
			if m.focus == focusGlamour {
				m.focus = focusGruff
			} else {
				m.focus = focusGlamour
			}
			return m, nil
		case "up", "down", "pgup", "pgdown":
			var cmd tea.Cmd
			if m.focus == focusGlamour {
				m.glamourView, cmd = m.glamourView.Update(msg)
			} else {
				m.gruffView, cmd = m.gruffView.Update(msg)
			}
			return m, cmd
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			separatorY := 1 + m.glamourView.Height
			if msg.Y < separatorY {
				m.focus = focusGlamour
			} else {
				m.focus = focusGruff
			}
		}
		var cmd tea.Cmd
		if m.focus == focusGlamour {
			m.glamourView, cmd = m.glamourView.Update(msg)
		} else {
			m.gruffView, cmd = m.gruffView.Update(msg)
		}
		return m, cmd
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		if !m.ready {
			w := msg.Width - 2
			h := (msg.Height - 5) / 2
			m.glamourView = viewport.New(w, h)
			m.gruffView = viewport.New(w, h)

			md := generateUnicodeMD("All printable Unicode characters from U+0020 to U+10FFFF, rendered by **glamour**.")

			r, err := glamour.NewTermRenderer(
				glamour.WithStyles(benchmark.GruffMinimalStyle()),
				glamour.WithWordWrap(w-2),
				glamour.WithTableWrap(false),
				glamour.WithInlineTableLinks(true),
			)
			if err != nil {
				log.Fatal(err)
			}
			out, err := r.Render(md)
			if err != nil {
				log.Fatal(err)
			}
			m.glamourView.SetContent(
				lipgloss.NewStyle().Width(w).Render(out),
			)

			md2 := generateUnicodeMD("All printable Unicode characters from U+0020 to U+10FFFF, rendered by **gruff**.")

			out2, err := gruff.Render(md2,
				gruff.WithWordWrap(w-2),
				func(o *gruff.Options) { o.Theme = noBgTheme },
			)
			if err != nil {
				log.Fatal(err)
			}
			m.gruffView.SetContent(
				lipgloss.NewStyle().Width(w).Render(out2),
			)

			m.ready = true
		}
	}
	return m, nil
}

func makeLabel(text string, active bool, activeBg, inactiveBg lipgloss.Color, width int) string {
	bg := inactiveBg
	if active {
		bg = activeBg
	}
	return lipgloss.NewStyle().
		Background(bg).
		Foreground(lipgloss.Color("#ffffff")).
		Width(width).
		Align(lipgloss.Center).
		Render(text)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Loading..."
	}
	width := m.termWidth
	if width == 0 {
		width = 80
	}
	w := width - 2

	sepBg := lipgloss.NewStyle().Background(lipgloss.Color("#333333"))
	fullSep := strings.Repeat("─", w)
	sepLine := sepBg.Width(w).Render(fullSep)

	glamourBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(lipgloss.Color("#444444"))
	gruffBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(lipgloss.Color("#444444"))
	if m.focus == focusGlamour {
		glamourBorder = glamourBorder.BorderForeground(lipgloss.Color("#7c3aed"))
	}
	if m.focus == focusGruff {
		gruffBorder = gruffBorder.BorderForeground(lipgloss.Color("#0891b2"))
	}

	glamourHeader := makeLabel(" glamour ", m.focus == focusGlamour,
		lipgloss.Color("#7c3aed"), lipgloss.Color("#3a1a6e"), w)

	gruffHeader := makeLabel(" gruff ", m.focus == focusGruff,
		lipgloss.Color("#0891b2"), lipgloss.Color("#045a6e"), w)

	glamourContent := glamourBorder.Render(
		glamourHeader + "\n" + m.glamourView.View(),
	)
	gruffContent := gruffBorder.Render(
		gruffHeader + "\n" + m.gruffView.View(),
	)

	joined := lipgloss.JoinVertical(lipgloss.Top,
		glamourContent,
		sepLine,
		gruffContent,
	)

	return lipgloss.NewStyle().
		Background(lipgloss.Color("#141414")).
		Width(width).
		Height(m.termHeight).
		Render(joined)
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
