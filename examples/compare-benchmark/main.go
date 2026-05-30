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
	focusRight
)

type renderTickMsg time.Time

const renderInterval = 1 * time.Second

func renderTick() tea.Cmd {
	return tea.Tick(renderInterval, func(t time.Time) tea.Msg {
		return renderTickMsg(t)
	})
}

type model struct {
	leftView   viewport.Model
	rightView  viewport.Model
	dirty      bool
	termWidth  int
	termHeight int
	focus      focus
	md         string

	glamourContent string
	gruffContent   string
	glamourDur     time.Duration
	gruffDur       time.Duration
}

func NewModel(md string) model {
	return model{md: md, dirty: true}
}

func (m model) Init() tea.Cmd {
	return renderTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case renderTickMsg:
		if m.dirty {
			m = m.renderAll()
			m.dirty = false
		}
		return m, renderTick()
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "left", "right":
			if m.focus == focusLeft {
				m.focus = focusRight
			} else {
				m.focus = focusLeft
			}
			return m, nil
		case "up", "down", "pgup", "pgdown":
			var cmd tea.Cmd
			switch m.focus {
			case focusLeft:
				m.leftView, cmd = m.leftView.Update(msg)
			case focusRight:
				m.rightView, cmd = m.rightView.Update(msg)
			}
			return m, cmd
		}
		return m, nil
	case tea.MouseMsg:
		halfW := (m.termWidth - 4) / 2
		if msg.X < halfW {
			m.focus = focusLeft
		} else {
			m.focus = focusRight
		}
		var cmd tea.Cmd
		switch m.focus {
		case focusLeft:
			m.leftView, cmd = m.leftView.Update(msg)
		case focusRight:
			m.rightView, cmd = m.rightView.Update(msg)
		}
		return m, cmd
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		halfW := (msg.Width - 4) / 2

		m.leftView = viewport.New(halfW, msg.Height-4)
		m.leftView.Style = lipgloss.NewStyle().Background(lipgloss.Color("#141414"))
		m.rightView = viewport.New(halfW, msg.Height-4)
		m.rightView.Style = lipgloss.NewStyle().Background(lipgloss.Color("#141414"))

		m.dirty = true
		return m, nil
	default:
		return m, nil
	}
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

func (m model) renderAll() model {
	halfW := (m.termWidth - 4) / 2

	t0 := time.Now()
	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(benchmark.GruffStandradStyle()),
		glamour.WithWordWrap(halfW),
	)
	if err != nil {
		log.Fatal(err)
	}
	out, err := r.Render(m.md)
	if err != nil {
		log.Fatal(err)
	}
	m.glamourContent = out
	m.glamourDur = time.Since(t0)

	t0 = time.Now()
	out2, err := gruff.Render(m.md,
		gruff.WithWordWrap(halfW),
	)
	if err != nil {
		log.Fatal(err)
	}
	m.gruffContent = out2
	m.gruffDur = time.Since(t0)

	m.leftView.SetContent(m.glamourContent)
	m.rightView.SetContent("\n" + m.gruffContent + "\n")
	return m
}

func (m model) headerFor(left bool) string {
	halfW := (m.termWidth - 4) / 2
	if left {
		active := m.focus == focusLeft
		return makeHeader(" glamour standard ", m.wxhInfo()+"  "+m.glamourDur.Round(time.Microsecond).String(),
			active, lipgloss.Color("#7c3aed"), lipgloss.Color("#3a1a6e"), halfW)
	}
	active := m.focus == focusRight
	return makeHeader(" gruff ", m.wxhInfo()+"  "+m.gruffDur.Round(time.Microsecond).String(),
		active, lipgloss.Color("#0891b2"), lipgloss.Color("#056982"), halfW)
}

func (m model) View() string {
	width := m.termWidth
	if width == 0 {
		width = 80
	}

	leftPane := m.paneBorder(m.focus == focusLeft, lipgloss.Color("#7c3aed")).Render(
		m.headerFor(true) + "\n" + m.leftView.View(),
	)

	rightPane := m.paneBorder(m.focus == focusRight, lipgloss.Color("#0891b2")).Render(
		m.headerFor(false) + "\n" + m.rightView.View(),
	)

	joined := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	return lipgloss.NewStyle().
		Background(lipgloss.Color("#141414")).
		Width(width).
		Height(m.termHeight).
		Render(joined)
}

func main() {
	b, err := os.ReadFile("testdata/benchmark.md")
	if err != nil {
		log.Fatal(err)
	}
	md := strings.TrimSpace(strings.Repeat(string(b), 100))
	p := tea.NewProgram(NewModel(md), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
