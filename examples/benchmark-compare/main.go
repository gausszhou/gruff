package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"charm.land/glamour/v2"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gausszhou/gruff"
	"github.com/gausszhou/gruff/benchmark"
)

type focus int

const (
	focusLeft focus = iota
	focusMinimal
	focusGruff
)

type model struct {
	leftView    viewport.Model
	minimalView viewport.Model
	gruffView   viewport.Model
	ready       bool
	termWidth   int
	termHeight  int
	focus       focus

	leftDur    time.Duration
	minimalDur time.Duration
	gruffDur   time.Duration
}

func readInput() string {
	b, err := os.ReadFile("testdata/benchmark.md")
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(strings.Repeat(string(b), 100))
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
		case "left":
			m.focus = focusLeft
			return m, nil
		case "right":
			if m.focus == focusLeft {
				m.focus = focusMinimal
			}
			return m, nil
		case "up", "down", "pgup", "pgdown":
			var cmd tea.Cmd
			switch m.focus {
			case focusLeft:
				m.leftView, cmd = m.leftView.Update(msg)
			case focusMinimal:
				m.minimalView, cmd = m.minimalView.Update(msg)
			case focusGruff:
				m.gruffView, cmd = m.gruffView.Update(msg)
			}
			return m, cmd
		}
	case tea.MouseMsg:
		halfW := (m.termWidth - 4) / 2
		halfH := (m.termHeight - 4) / 2
		if msg.X < halfW {
			m.focus = focusLeft
		} else if msg.Y < halfH {
			m.focus = focusMinimal
		} else {
			m.focus = focusGruff
		}
		var cmd tea.Cmd
		switch m.focus {
		case focusLeft:
			m.leftView, cmd = m.leftView.Update(msg)
		case focusMinimal:
			m.minimalView, cmd = m.minimalView.Update(msg)
		case focusGruff:
			m.gruffView, cmd = m.gruffView.Update(msg)
		}
		return m, cmd
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		if !m.ready {
			halfW := (msg.Width-4)/2 - 2
			rightH := (msg.Height - 6) / 2

			m.leftView = viewport.New(halfW, msg.Height-4)
			m.minimalView = viewport.New(halfW, rightH)
			m.gruffView = viewport.New(halfW, rightH)

			md := readInput()

			t0 := time.Now()
			r, err := glamour.NewTermRenderer(
				glamour.WithStandardStyle("dark"),
				glamour.WithWordWrap(halfW-2),
			)
			if err != nil {
				log.Fatal(err)
			}
			out, err := r.Render(md)
			if err != nil {
				log.Fatal(err)
			}
			m.leftView.SetContent(out)
			m.leftDur = time.Since(t0)

			t0 = time.Now()
			r2, err := glamour.NewTermRenderer(
				glamour.WithStyles(benchmark.GruffMinimalStyle()),
				glamour.WithChromaFormatter("noop"),
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

			t0 = time.Now()
			out3, err := gruff.Render(md,
				gruff.WithWordWrap(halfW-2),
			)
			if err != nil {
				log.Fatal(err)
			}
			m.gruffView.SetContent(out3)
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
	right := base.Foreground(lipgloss.Color("#ffffff")).Render(info)
	return left + middle + right
}

func (m model) paneBorder(active bool, activeColor lipgloss.Color) lipgloss.Style {
	s := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true)
	if active {
		return s.BorderForeground(activeColor)
	}
	return s.BorderForeground(lipgloss.Color("#444444"))
}

func (m model) wxhInfo() string {
	return fmt.Sprintf("%dx%d", m.termWidth, m.termHeight)
}

func (m model) leftInfo() string {
	return m.wxhInfo() + "  " + m.leftDur.Round(time.Microsecond).String()
}

func (m model) minimalInfo() string {
	return m.wxhInfo() + "  " + m.minimalDur.Round(time.Microsecond).String()
}

func (m model) gruffInfo() string {
	return m.wxhInfo() + "  " + m.gruffDur.Round(time.Microsecond).String()
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

	leftPane := m.paneBorder(m.focus == focusLeft, lipgloss.Color("#7c3aed")).Render(
		makeHeader(" glamour standard ", m.leftInfo(), m.focus == focusLeft,
			lipgloss.Color("#7c3aed"), lipgloss.Color("#3a1a6e"), halfW) + "\n" + m.leftView.View(),
	)

	minimalPane := m.paneBorder(m.focus == focusMinimal, lipgloss.Color("#059669")).Render(
		makeHeader(" glamour minimal ", m.minimalInfo(), m.focus == focusMinimal,
			lipgloss.Color("#059669"), lipgloss.Color("#056f4d"), halfW) + "\n" + m.minimalView.View(),
	)

	gruffPane := m.paneBorder(m.focus == focusGruff, lipgloss.Color("#0891b2")).Render(
		makeHeader(" gruff ", m.gruffInfo(), m.focus == focusGruff,
			lipgloss.Color("#0891b2"), lipgloss.Color("#056982"), halfW) + "\n" + m.gruffView.View(),
	)

	rightSide := lipgloss.JoinVertical(lipgloss.Top, minimalPane, gruffPane)

	joined := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightSide)

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
