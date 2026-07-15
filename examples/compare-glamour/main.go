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
	leftView   component.ViewportWithScrollbar
	rightView  component.ViewportWithScrollbar
	dirty      bool
	termWidth  int
	termHeight int
	focus      focus
	md         string

	viewWidth  int
	viewHeight int

	leftContent   string
	rightContent  string
	leftDur       time.Duration
	rightDur      time.Duration
	leftRenderer  *glamour.TermRenderer
	rightRenderer *glamour.TermRenderer
	renderWidth   int
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
		case "q", "ctrl+c":
			return m, tea.Quit
		case "left", "right":
			if m.focus == focusLeft {
				m.focus = focusRight
			} else {
				m.focus = focusLeft
			}
			return m, nil
		case "up", "down", "pgup", "pgdown", "home", "end", "g", "G", "ctrl+u", "ctrl+d":
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
		if msg.Mouse().X < m.viewWidth {
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
		m.leftView.OriginX = 1
		m.leftView.OriginY = 2
		m.rightView = component.NewViewportWithScrollbar(m.viewWidth, m.viewHeight)
		m.rightView.OriginX = m.viewWidth + 3
		m.rightView.OriginY = 2

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
	halfW := m.viewWidth - 1

	if m.leftRenderer == nil || m.renderWidth != halfW {
		var err error
		m.leftRenderer, err = glamour.NewTermRenderer(
			glamour.WithStyles(benchmark.GlamourMinimalStyle()),
			glamour.WithWordWrap(halfW),
			glamour.WithTableWrap(false),
			glamour.WithInlineTableLinks(false),
		)
		if err != nil {
			log.Fatal(err)
		}
		m.rightRenderer, err = glamour.NewTermRenderer(
			glamour.WithStyles(benchmark.GlamourStandardStyle()),
			glamour.WithWordWrap(halfW),
		)
		if err != nil {
			log.Fatal(err)
		}
		m.renderWidth = halfW
	}

	cleaned := benchmark.CleanInput(m.md)

	t1 := time.Now()
	out, err := m.leftRenderer.Render(cleaned)
	if err != nil {
		log.Fatal(err)
	}
	m.leftContent = out
	m.leftDur = time.Since(t1)

	t2 := time.Now()
	out2, err := m.rightRenderer.Render(m.md)
	if err != nil {
		log.Fatal(err)
	}
	m.rightContent = out2
	m.rightDur = time.Since(t2)

	m.leftView.SetContent(m.leftContent)
	m.rightView.SetContent(m.rightContent)
	return m
}

func (m model) headerFor(left bool) string {
	halfW := (m.termWidth - 4) / 2
	var title, info string
	var activeBg, inactiveBg color.Color
	if left {
		title = "glamour minimal"
		info = m.wxhInfo() + "  " + m.leftDur.Round(time.Microsecond).String()
		activeBg = lipgloss.Color("#059669")
		inactiveBg = lipgloss.Color("#056f4d")
	} else {
		title = "glamour standard"
		info = m.wxhInfo() + "  " + m.rightDur.Round(time.Microsecond).String()
		activeBg = lipgloss.Color("#7c3aed")
		inactiveBg = lipgloss.Color("#3a1a6e")
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

	leftPane := m.paneBorder(m.focus == focusLeft, lipgloss.Color("#059669")).Render(
		m.headerFor(true) + "\n" + m.leftView.View(),
	)

	rightPane := m.paneBorder(m.focus == focusRight, lipgloss.Color("#7c3aed")).Render(
		m.headerFor(false) + "\n" + m.rightView.View(),
	)

	joined := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	v := tea.NewView(joined)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeAllMotion
	return v
}

func main() {
	b, err := os.ReadFile("testdata/benchmark.md")
	if err != nil {
		log.Fatal(err)
	}
	md := strings.Repeat(string(b), 100)
	p := tea.NewProgram(NewModel(md))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
