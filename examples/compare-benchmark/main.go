package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
	"time"

	"charm.land/glamour/v2"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	flex "github.com/gausszhou/bubbleflex"
	"github.com/gausszhou/gruff/benchmark"
	"github.com/gausszhou/gruff/component"
	"github.com/gausszhou/gruff/gruff"
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
	leftView  component.ViewportWithScrollbar
	rightView component.ViewportWithScrollbar

	termWidth  int
	termHeight int

	viewWidth  int
	viewHeight int

	focus focus
	dirty bool

	md string

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
		}
		var cmd tea.Cmd
		switch m.focus {
		case focusLeft:
			m.leftView, cmd = m.leftView.Update(msg)
		case focusRight:
			m.rightView, cmd = m.rightView.Update(msg)
		}
		return m, cmd
	case tea.MouseMsg:
		halfW := m.termWidth / 2
		if msg.Mouse().X < halfW {
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
		m.viewWidth = (msg.Width - 4) / 2
		m.viewHeight = msg.Height - 3

		m.leftView = component.NewViewportWithScrollbar(m.viewWidth, m.viewHeight)
		m.rightView = component.NewViewportWithScrollbar(m.viewWidth, m.viewHeight)
		m.dirty = true
		return m, nil
	default:
		return m, nil
	}
}

func makeHeader(width int, left, right string) string {
	return flex.New(flex.Row).JustifyContent(flex.SpaceBetween).Width(width).Join(left, right)
}

func (m model) paneBorder(active bool, activeColor color.Color) lipgloss.Style {
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
	wrapWidth := m.viewWidth - 1

	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(benchmark.GlamourStandardStyle()),
		glamour.WithWordWrap(wrapWidth),
	)
	if err != nil {
		log.Fatal(err)
	}

	cleaned := benchmark.CleanInput(m.md)

	t1 := time.Now()
	out, err := r.Render(cleaned)
	if err != nil {
		log.Fatal(err)
	}
	m.glamourContent = out
	m.glamourDur = time.Since(t1)

	t2 := time.Now()
	out2, err := gruff.Render(m.md,
		gruff.WithWordWrap(wrapWidth),
	)
	if err != nil {
		log.Fatal(err)
	}
	m.gruffContent = out2
	m.gruffDur = time.Since(t2)

	m.leftView.SetContent("\n" + m.gruffContent + "\n")
	m.rightView.SetContent(m.glamourContent)
	return m
}

func (m model) headerFor(left bool) string {
	halfW := (m.termWidth - 4) / 2
	var title, info string
	var activeBg, inactiveBg color.Color
	if left {
		title = "gruff"
		info = m.wxhInfo() + "  " + m.gruffDur.Round(time.Microsecond).String()
		activeBg = lipgloss.Color("#0891b2")
		inactiveBg = lipgloss.Color("#056982")
	} else {
		title = "glamour minimal"
		info = m.wxhInfo() + "  " + m.glamourDur.Round(time.Microsecond).String()
		activeBg = lipgloss.Color("#059669")
		inactiveBg = lipgloss.Color("#056f4d")
	}
	active := (left && m.focus == focusLeft) || (!left && m.focus == focusRight)
	bg := inactiveBg
	if active {
		bg = activeBg
	}
	return lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color("#ffffff")).Padding(0, 1).Render(
		makeHeader(halfW-2, title, info),
	)
}

func (m model) View() tea.View {
	width := m.termWidth
	if width == 0 {
		width = 80
	}

	leftPane := m.paneBorder(m.focus == focusLeft, lipgloss.Color("#0891b2")).Render(
		m.headerFor(true) + "\n" + m.leftView.View(),
	)

	rightPane := m.paneBorder(m.focus == focusRight, lipgloss.Color("#059669")).Render(
		m.headerFor(false) + "\n" + m.rightView.View(),
	)

	joined := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	rendered := lipgloss.NewStyle().
		Width(width).
		Height(m.termHeight).
		Render(joined)

	v := tea.NewView(rendered)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeAllMotion
	return v
}

func main() {
	b, err := os.ReadFile("testdata/benchmark.md")
	if err != nil {
		log.Fatal(err)
	}
	md := strings.TrimSpace(strings.Repeat(string(b), 100))
	p := tea.NewProgram(NewModel(md))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
