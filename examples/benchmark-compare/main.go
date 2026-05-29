package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"charm.land/glamour/v2"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gausszhou/gruff/benchmark"
)

type focus int

const (
	focusStandard focus = iota
	focusMinimal
)

type model struct {
	standardView viewport.Model
	minimalView  viewport.Model
	ready        bool
	termWidth    int
	termHeight   int
	focus        focus

	standardDur time.Duration
	minimalDur  time.Duration
}

func readInput() string {
	b, err := os.ReadFile("testdata/benchmark.md")
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
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
			if m.focus == focusStandard {
				m.focus = focusMinimal
			} else {
				m.focus = focusStandard
			}
			return m, nil
		case "left":
			m.focus = focusStandard
			return m, nil
		case "right":
			m.focus = focusMinimal
			return m, nil
		case "up", "down", "pgup", "pgdown":
			var cmd tea.Cmd
			if m.focus == focusStandard {
				m.standardView, cmd = m.standardView.Update(msg)
			} else {
				m.minimalView, cmd = m.minimalView.Update(msg)
			}
			return m, cmd
		}
	case tea.MouseMsg:
		halfW := (m.termWidth - 4) / 2
		if msg.X < halfW {
			m.focus = focusStandard
		} else {
			m.focus = focusMinimal
		}
		var cmd tea.Cmd
		if m.focus == focusStandard {
			m.standardView, cmd = m.standardView.Update(msg)
		} else {
			m.minimalView, cmd = m.minimalView.Update(msg)
		}
		return m, cmd
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		if !m.ready {
			halfW := (msg.Width-4)/2 - 2
			h := msg.Height - 4

			m.standardView = viewport.New(halfW, h)
			m.minimalView = viewport.New(halfW, h)

			md := readInput()

			t0 := time.Now()
			r, err := glamour.NewTermRenderer(
				glamour.WithStandardStyle("dark"),
				glamour.WithWordWrap(halfW - 2),
			)
			if err != nil {
				log.Fatal(err)
			}
			out, err := r.Render(md)
			if err != nil {
				log.Fatal(err)
			}
			m.standardView.SetContent(out)
			m.standardDur = time.Since(t0)

			t0 = time.Now()
			r2, err := glamour.NewTermRenderer(
				glamour.WithStyles(benchmark.GruffMinimalStyle()),
				glamour.WithWordWrap(halfW-2),
				glamour.WithTableWrap(false),
				glamour.WithInlineTableLinks(true),
			)
			if err != nil {
				log.Fatal(err)
			}
			out2, err := r2.Render(md)
			if err != nil {
				log.Fatal(err)
			}
			m.minimalView.SetContent(out2)
			m.minimalDur = time.Since(t0)

			m.ready = true
		}
	}
	return m, nil
}

func makeHeader(title, info string, active bool, activeBg, inactiveBg lipgloss.Color, width int) string {
	bg := inactiveBg
	if active {
		bg = activeBg
	}
	base := lipgloss.NewStyle().Background(bg)
	left := base.Foreground(lipgloss.Color("#ffffff")).Render(title)
	gap := width - lipgloss.Width(left) - lipgloss.Width(info)
	if gap < 0 {
		gap = 0
	}
	middle := base.Width(gap).Render("")
	right := base.Foreground(lipgloss.Color("#888888")).Render(info)
	return left + middle + right
}

func (m model) View() string {
	if !m.ready {
		return "\n  Loading..."
	}
	width := m.termWidth
	if width == 0 {
		width = 80
	}
	halfW := (width-4)/2 - 2

	standardBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(lipgloss.Color("#444444"))
	minimalBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(lipgloss.Color("#444444"))
	if m.focus == focusStandard {
		standardBorder = standardBorder.BorderForeground(lipgloss.Color("#7c3aed"))
	}
	if m.focus == focusMinimal {
		minimalBorder = minimalBorder.BorderForeground(lipgloss.Color("#059669"))
	}

	wxh := fmt.Sprintf("%dx%d", m.termWidth, m.termHeight)

	standardInfo := wxh + "  " + m.standardDur.Round(time.Microsecond).String()
	minimalInfo := wxh + "  " + m.minimalDur.Round(time.Microsecond).String()

	standardHeader := makeHeader(" glamour standard ", standardInfo, m.focus == focusStandard,
		lipgloss.Color("#7c3aed"), lipgloss.Color("#3a1a6e"), halfW)
	minimalHeader := makeHeader(" glamour minimal ", minimalInfo, m.focus == focusMinimal,
		lipgloss.Color("#059669"), lipgloss.Color("#065f46"), halfW)

	standardContent := standardBorder.Render(
		standardHeader + "\n" + m.standardView.View(),
	)
	minimalContent := minimalBorder.Render(
		minimalHeader + "\n" + m.minimalView.View(),
	)

	joined := lipgloss.JoinHorizontal(lipgloss.Top,
		standardContent,
		minimalContent,
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
