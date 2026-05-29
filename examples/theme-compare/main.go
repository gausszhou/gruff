package main

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gausszhou/gruff"
)

var sampleMD = "# Gruff\n\n" +
	"A **lightweight** markdown renderer for the terminal.\n\n" +
	"## Text Formatting\n\n" +
	"Markdown supports **bold** and *italic* text, as well as ***bold italic***.\n" +
	"You can also use `inline code` for short snippets, or ~~strikethrough~~ for crossed-out text.\n" +
	"Standard paragraphs are the most common element in any document.\n" +
	"## Text\n\n" +
	"- *Italic* and **bold**\n" +
	"- `inline code`\n" +
	"- ***bold italic***\n" +
	"- ~~strikethrough~~\n\n" +
	"## Table\n\n" +
	"| Left | Center | Right |\n" +
	"|:-----|:------:|------:|\n" +
	"| a | b | c |\n" +
	"| 1 | 2 | 3 |\n\n" +
	"## Code\n\n" +
	"```\n" +
	"func main() {\n" +
	"    fmt.Println(\"Hello\")\n" +
	"}\n" +
	"```\n\n" +
	"> A wise quote.\n\n" +
	"\n\n" +
	"- [x] done\n" +
	"- [ ] todo\n\n" +
	"---\n\n" +
	"Visit [Gruff](https://github.com/gausszhou/gruff) for more.\n"

type focus int

const (
	focusDark focus = iota
	focusLight
)

type model struct {
	darkView   viewport.Model
	lightView  viewport.Model
	ready      bool
	termWidth  int
	termHeight int
	focus      focus
	darkDur    time.Duration
	lightDur   time.Duration
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
			if m.focus == focusDark {
				m.focus = focusLight
			} else {
				m.focus = focusDark
			}
			return m, nil
		case "up", "down", "pgup", "pgdown":
			var cmd tea.Cmd
			switch m.focus {
			case focusDark:
				m.darkView, cmd = m.darkView.Update(msg)
			case focusLight:
				m.lightView, cmd = m.lightView.Update(msg)
			}
			return m, cmd
		}
	case tea.MouseMsg:
		halfW := m.termWidth / 2
		if msg.X < halfW {
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
		if !m.ready {
			w := (msg.Width-4)/2 - 2
			h := msg.Height - 4
			m.darkView = viewport.New(w, h)
			m.darkView.Style = lipgloss.NewStyle().Background(lipgloss.Color("#141414")).Foreground(lipgloss.Color("#ffffff"))
			m.lightView = viewport.New(w, h)
			m.lightView.Style = lipgloss.NewStyle().Background(lipgloss.Color("#ffffff")).Foreground(lipgloss.Color("#000000"))

			t0 := time.Now()
			dark, err := gruff.Render(sampleMD,
				gruff.WithDark(),
				gruff.WithWordWrap(w-2),
			)
			if err != nil {
				log.Fatal(err)
			}
			m.darkView.SetContent(dark)
			m.darkDur = time.Since(t0)

			t0 = time.Now()
			light, err := gruff.Render(sampleMD,
				gruff.WithLight(),
				gruff.WithWordWrap(w-2),
			)
			if err != nil {
				log.Fatal(err)
			}
			m.lightView.SetContent(light)
			m.lightDur = time.Since(t0)

			m.ready = true
		}
	}
	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "\n  Loading..."
	}
	w := (m.termWidth-4)/2 - 2

	darkActive := m.focus == focusDark
	lightActive := m.focus == focusLight

	darkHeader := lipgloss.NewStyle().
		Background(lipgloss.Color("#222222")).
		Foreground(lipgloss.Color("#ffffff")).
		Width(w).
		Render(" dark  " + m.darkDur.Round(time.Microsecond).String())
	lightHeader := lipgloss.NewStyle().
		Background(lipgloss.Color("#dddddd")).
		Foreground(lipgloss.Color("#000000")).
		Width(w).
		Render(" light " + m.lightDur.Round(time.Microsecond).String())

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

	return lipgloss.NewStyle().
		Width(m.termWidth).
		Height(m.termHeight).
		Render(joined)
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
