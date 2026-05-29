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
	"github.com/gausszhou/gruff"
)

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

	glamourDur time.Duration
	gruffDur   time.Duration
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
			if m.focus == focusGlamour {
				m.focus = focusGruff
			} else {
				m.focus = focusGlamour
			}
			return m, nil
		case "left":
			m.focus = focusGlamour
			return m, nil
		case "right":
			m.focus = focusGruff
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
		halfW := (m.termWidth - 4) / 2
		if msg.X < halfW {
			m.focus = focusGlamour
		} else {
			m.focus = focusGruff
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
			halfW := (msg.Width-4)/2 - 2
			h := msg.Height - 4

			m.glamourView = viewport.New(halfW, h)
			m.gruffView = viewport.New(halfW, h)

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
			m.glamourView.SetContent(out)
			m.glamourDur = time.Since(t0)

			t0 = time.Now()
			out2, err := gruff.Render(md,
				gruff.WithWordWrap(halfW-2),
			)
			if err != nil {
				log.Fatal(err)
			}
			m.gruffView.SetContent(out2)
			m.gruffDur = time.Since(t0)

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

	wxh := fmt.Sprintf("%dx%d", m.termWidth, m.termHeight)

	glamourInfo := wxh + "  " + m.glamourDur.Round(time.Microsecond).String()
	gruffInfo := wxh + "  " + m.gruffDur.Round(time.Microsecond).String()

	glamourHeader := makeHeader(" glamour ", glamourInfo, m.focus == focusGlamour,
		lipgloss.Color("#7c3aed"), lipgloss.Color("#3a1a6e"), halfW)
	gruffHeader := makeHeader(" gruff ", gruffInfo, m.focus == focusGruff,
		lipgloss.Color("#0891b2"), lipgloss.Color("#045a6e"), halfW)

	glamourContent := glamourBorder.Render(
		glamourHeader + "\n" + m.glamourView.View(),
	)
	gruffContent := gruffBorder.Render(
		gruffHeader + "\n" + m.gruffView.View(),
	)

	joined := lipgloss.JoinHorizontal(lipgloss.Top,
		glamourContent,
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
