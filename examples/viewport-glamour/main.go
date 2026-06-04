package main

import (
	"log"
	"os"
	"time"

	"charm.land/glamour/v2"
	tea "charm.land/bubbletea/v2"
	"github.com/gausszhou/gruff/benchmark"
	"github.com/gausszhou/gruff/component"
)

type renderTickMsg time.Time

const renderInterval = 200 * time.Millisecond

func renderTick() tea.Cmd {
	return tea.Tick(renderInterval, func(t time.Time) tea.Msg {
		return renderTickMsg(t)
	})
}

type model struct {
	viewport  component.ViewportWithScrollbar
	mdContent string
	dirty     bool
	termWidth int
}

func (m model) Init() tea.Cmd {
	return renderTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case renderTickMsg:
		if m.dirty && m.termWidth > 0 {
			m.dirty = false
			r, err := glamour.NewTermRenderer(
				glamour.WithStyles(benchmark.GlamourStandardStyle()),
				glamour.WithWordWrap(m.termWidth-1),
			)
			if err != nil {
				log.Fatal(err)
			}
			out, err := r.Render(benchmark.CleanInput(m.mdContent))
			if err != nil {
				log.Fatal(err)
			}
			m.viewport.SetContent(out)
		}
		return m, renderTick()

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.viewport = component.NewViewportWithScrollbar(msg.Width, msg.Height)
		m.viewport.Inner().MouseWheelEnabled = true
		m.dirty = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tea.MouseMsg:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	if m.dirty {
		return tea.NewView("Loading...")
	}
	v := tea.NewView(m.viewport.View())
	v.AltScreen = true
	v.MouseMode = tea.MouseModeAllMotion
	return v
}

func main() {
	b, err := os.ReadFile("testdata/benchmark.md")
	if err != nil {
		log.Fatal(err)
	}
	p := tea.NewProgram(model{mdContent: string(b), dirty: true})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
