package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	flex "github.com/gausszhou/bubbleflex"
	"github.com/gausszhou/gruff/component"
	"github.com/gausszhou/gruff/gruff"
)

type focus int

const (
	focusDark focus = iota
	focusLight
)

type renderTickMsg time.Time

const renderInterval = 1 * time.Second

func renderTick() tea.Cmd {
	return tea.Tick(renderInterval, func(t time.Time) tea.Msg {
		return renderTickMsg(t)
	})
}

type model struct {
	md string

	termWidth  int
	termHeight int
	viewWidth  int
	viewHeight int

	darkView  component.ViewportWithScrollbar
	lightView component.ViewportWithScrollbar
	renderW   int

	dirty bool

	focus    focus
	darkDur  time.Duration
	lightDur time.Duration
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
		case "tab":
			if m.focus == focusDark {
				m.focus = focusLight
			} else {
				m.focus = focusDark
			}
			return m, nil
		case "up", "down", "pgup", "pgdown", "home", "end", "g", "G", "ctrl+u", "ctrl+d":
			var cmd tea.Cmd
			switch m.focus {
			case focusDark:
				m.darkView, cmd = m.darkView.Update(msg)
			case focusLight:
				m.lightView, cmd = m.lightView.Update(msg)
			}
			return m, cmd
		}
		return m, nil
	case tea.MouseMsg:
		halfW := m.termWidth / 2
		if msg.Mouse().X < halfW {
			m.focus = focusDark
		} else {
			m.focus = focusLight
		}
		var cmd tea.Cmd
		switch m.focus {
		case focusDark:
			m.darkView, cmd = m.darkView.Update(msg)
		case focusLight:
			m.lightView, cmd = m.lightView.Update(msg)
		}
		return m, cmd
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		m.viewWidth = (msg.Width - 4) / 2
		m.viewHeight = msg.Height - 4
		m.darkView = component.NewViewportWithScrollbar(m.viewWidth, m.viewHeight)
		m.darkView.OriginX = 1
		m.darkView.OriginY = 2
		m.darkView.Inner().Style = lipgloss.NewStyle().Background(lipgloss.Color("#141414")).Foreground(lipgloss.Color("#ffffff")).Padding(1)
		m.lightView = component.NewViewportWithScrollbar(m.viewWidth, m.viewHeight)
		m.lightView.OriginX = m.viewWidth + 3
		m.lightView.OriginY = 2
		m.lightView.Inner().Style = lipgloss.NewStyle().Background(lipgloss.Color("#f0f0f0")).Foreground(lipgloss.Color("#000000")).Padding(1)
		m.dirty = true
	}
	return m, nil
}

func (m model) renderAll() model {
	wrapWidth := m.viewWidth - 1

	t0 := time.Now()
	dark, err := gruff.Render(m.md,
		gruff.WithDark(),
		gruff.WithWordWrap(wrapWidth),
	)
	if err != nil {
		log.Fatal(err)
	}
	m.darkView.SetContent(dark)
	m.darkDur = time.Since(t0)

	t1 := time.Now()
	light, err := gruff.Render(m.md,
		gruff.WithLight(),
		gruff.WithWordWrap(wrapWidth),
	)
	if err != nil {
		log.Fatal(err)
	}
	m.lightView.SetContent(light)
	m.lightDur = time.Since(t1)
	return m
}

func (m model) wxhInfo() string {
	return fmt.Sprintf("%dx%d", m.termWidth, m.termHeight)
}

func makeHeader(width int, left, right string) string {
	return flex.New(flex.Row).JustifyContent(flex.SpaceBetween).Width(width).Join(left, right)
}

func (m model) View() tea.View {
	darkActive := m.focus == focusDark
	lightActive := m.focus == focusLight

	darkHeader := lipgloss.NewStyle().Background(lipgloss.Color("#222222")).Foreground(lipgloss.Color("#ffffff")).Padding(0, 1).Render(
		makeHeader(m.viewWidth-2, "dark", m.wxhInfo()+"  "+m.darkDur.Round(time.Microsecond).String()),
	)
	lightHeader := lipgloss.NewStyle().Background(lipgloss.Color("#dddddd")).Foreground(lipgloss.Color("#000000")).Padding(0, 1).Render(
		makeHeader(m.viewWidth-2, "light", m.wxhInfo()+"  "+m.lightDur.Round(time.Microsecond).String()),
	)

	darkBorderColor := lipgloss.Color("#444444")
	lightBorderColor := lipgloss.Color("#444444")
	if darkActive {
		darkBorderColor = lipgloss.Color("#7c3aed")
	}
	if lightActive {
		lightBorderColor = lipgloss.Color("#f5a623")
	}

	border := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true)
	darkPane := border.BorderForeground(darkBorderColor).
		Background(lipgloss.Color("#141414")).
		Render(darkHeader + "\n" + m.darkView.View())
	lightPane := border.BorderForeground(lightBorderColor).
		Background(lipgloss.Color("#ffffff")).
		Render(lightHeader + "\n" + m.lightView.View())

	joined := lipgloss.JoinHorizontal(lipgloss.Top, darkPane, lightPane)

	rendered := lipgloss.NewStyle().
		Width(m.termWidth).
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
	md := strings.Repeat(string(b), 100)
	p := tea.NewProgram(model{md: md})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
