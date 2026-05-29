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
	"github.com/gausszhou/gruff/benchmark"
)

type focus int

const (
	focusLeft focus = iota
	focusRight
)

type rightMode int

const (
	rightMinimal rightMode = iota
	rightGruff
)

type model struct {
	leftView   viewport.Model
	rightView  viewport.Model
	ready      bool
	termWidth  int
	termHeight int
	focus      focus
	rightMode  rightMode

	leftDur    time.Duration
	minimalDur time.Duration
	gruffDur   time.Duration
	minimalOut string
	gruffOut   string
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
			if m.rightMode == rightMinimal {
				m.rightMode = rightGruff
				m.rightView.SetContent(m.gruffOut)
			} else {
				m.rightMode = rightMinimal
				m.rightView.SetContent(m.minimalOut)
			}
			return m, nil
		case "left":
			m.focus = focusLeft
			return m, nil
		case "right":
			m.focus = focusRight
			return m, nil
		case "up", "down", "pgup", "pgdown":
			var cmd tea.Cmd
			if m.focus == focusLeft {
				m.leftView, cmd = m.leftView.Update(msg)
			} else {
				m.rightView, cmd = m.rightView.Update(msg)
			}
			return m, cmd
		}
	case tea.MouseMsg:
		halfW := (m.termWidth - 4) / 2
		if msg.X < halfW {
			m.focus = focusLeft
		} else {
			m.focus = focusRight
		}
		var cmd tea.Cmd
		if m.focus == focusLeft {
			m.leftView, cmd = m.leftView.Update(msg)
		} else {
			m.rightView, cmd = m.rightView.Update(msg)
		}
		return m, cmd
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		if !m.ready {
			halfW := (msg.Width-4)/2 - 2
			h := msg.Height - 4

			m.leftView = viewport.New(halfW, h)
			m.rightView = viewport.New(halfW, h)

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
			m.minimalOut, err = r2.Render(md)
			if err != nil {
				log.Fatal(err)
			}
			m.minimalDur = time.Since(t0)

			t0 = time.Now()
			m.gruffOut, err = gruff.Render(md,
				gruff.WithWordWrap(halfW-2),
			)
			if err != nil {
				log.Fatal(err)
			}
			m.gruffDur = time.Since(t0)

			m.rightView.SetContent(m.minimalOut)

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

func (m model) leftBorderStyle() lipgloss.Style {
	s := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true)
	if m.focus == focusLeft {
		return s.BorderForeground(lipgloss.Color("#7c3aed"))
	}
	return s.BorderForeground(lipgloss.Color("#444444"))
}

func (m model) rightBorderStyle() lipgloss.Style {
	s := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true)
	if m.focus != focusRight {
		return s.BorderForeground(lipgloss.Color("#444444"))
	}
	if m.rightMode == rightGruff {
		return s.BorderForeground(lipgloss.Color("#0891b2"))
	}
	return s.BorderForeground(lipgloss.Color("#059669"))
}

func (m model) wxhInfo() string {
	return fmt.Sprintf("%dx%d", m.termWidth, m.termHeight)
}

func (m model) leftInfo() string {
	return m.wxhInfo() + "  " + m.leftDur.Round(time.Microsecond).String()
}

func (m model) rightInfo() (title, info string, activeBg, inactiveBg lipgloss.Color) {
	if m.rightMode == rightGruff {
		return " gruff ", m.wxhInfo() + "  " + m.gruffDur.Round(time.Microsecond).String(),
			lipgloss.Color("#0891b2"), lipgloss.Color("#056982")
	}
	return " glamour minimal ", m.wxhInfo() + "  " + m.minimalDur.Round(time.Microsecond).String(),
		lipgloss.Color("#059669"), lipgloss.Color("#056f4d")
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

	leftPane := m.leftBorderStyle().Render(
		makeHeader(" glamour standard ", m.leftInfo(), m.focus == focusLeft,
			lipgloss.Color("#7c3aed"), lipgloss.Color("#3a1a6e"), halfW) + "\n" + m.leftView.View(),
	)

	rightTitle, rightInfo, rightBg, rightBgInactive := m.rightInfo()
	rightPane := m.rightBorderStyle().Render(
		makeHeader(rightTitle, rightInfo, m.focus == focusRight, rightBg, rightBgInactive, halfW) + "\n" + m.rightView.View(),
	)

	joined := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

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
