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
	focusTop
	focusBottom
)

type slot int

const (
	slotStandard slot = iota
	slotMinimal
	slotGruff
)

var slotNames = map[slot]string{
	slotStandard: " glamour standard ",
	slotMinimal:  " glamour minimal ",
	slotGruff:    " gruff ",
}

var slotColors = map[slot]struct {
	activeBg, inactiveBg lipgloss.Color
}{
	slotStandard: {lipgloss.Color("#7c3aed"), lipgloss.Color("#3a1a6e")},
	slotMinimal:  {lipgloss.Color("#059669"), lipgloss.Color("#056f4d")},
	slotGruff:    {lipgloss.Color("#0891b2"), lipgloss.Color("#056982")},
}

var slotFocusColor = map[focus]lipgloss.Color{
	focusLeft:   lipgloss.Color("#7c3aed"),
	focusTop:    lipgloss.Color("#059669"),
	focusBottom: lipgloss.Color("#0891b2"),
}

type renderTickMsg time.Time

const renderInterval = 1 * time.Second

func renderTick() tea.Cmd {
	return tea.Tick(renderInterval, func(t time.Time) tea.Msg {
		return renderTickMsg(t)
	})
}

type model struct {
	leftView   viewport.Model
	topView    viewport.Model
	bottomView viewport.Model
	dirty      bool
	termWidth  int
	termHeight int
	focus      focus
	perm       int
	md         string

	contents  [3]string
	durations [3]time.Duration
}

func (m model) slotAt(pos int) slot {
	return slot((pos + m.perm) % 3)
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
		case "tab":
			m.perm = (m.perm + 1) % 3
			m.leftView.SetContent(m.contents[m.slotAt(0)])
			m.topView.SetContent(m.contents[m.slotAt(1)])
			m.bottomView.SetContent(m.contents[m.slotAt(2)])
			return m, nil
		case "left":
			m.focus = focusLeft
			return m, nil
		case "right":
			if m.focus == focusLeft {
				m.focus = focusTop
			}
			return m, nil
		case "up", "down", "pgup", "pgdown":
			var cmd tea.Cmd
			switch m.focus {
			case focusLeft:
				m.leftView, cmd = m.leftView.Update(msg)
			case focusTop:
				m.topView, cmd = m.topView.Update(msg)
			case focusBottom:
				m.bottomView, cmd = m.bottomView.Update(msg)
			}
			return m, cmd
		}
		return m, nil
	case tea.MouseMsg:
		halfW := (m.termWidth - 4) / 2
		halfH := (m.termHeight - 4) / 2
		if msg.X < halfW {
			m.focus = focusLeft
		} else if msg.Y < halfH {
			m.focus = focusTop
		} else {
			m.focus = focusBottom
		}
		var cmd tea.Cmd
		switch m.focus {
		case focusLeft:
			m.leftView, cmd = m.leftView.Update(msg)
		case focusTop:
			m.topView, cmd = m.topView.Update(msg)
		case focusBottom:
			m.bottomView, cmd = m.bottomView.Update(msg)
		}
		return m, cmd
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		halfW := (msg.Width-2)/2 - 2
		rightH := (msg.Height - 6) / 2

		m.leftView = viewport.New(halfW, msg.Height-3)
		m.leftView.Style = lipgloss.NewStyle().Background(lipgloss.Color("#141414"))
		m.topView = viewport.New(halfW, rightH)
		m.topView.Style = lipgloss.NewStyle().Background(lipgloss.Color("#141414"))
		m.bottomView = viewport.New(halfW, rightH)
		m.bottomView.Style = lipgloss.NewStyle().Background(lipgloss.Color("#141414"))

		m.dirty = true

		return m, nil
		// return m, func() tea.Msg { return renderTickMsg(time.Now()) }
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
	halfW := (m.termWidth-2)/2 - 1

	t0 := time.Now()
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(halfW-2),
	)
	if err != nil {
		log.Fatal(err)
	}
	out, err := r.Render(m.md)
	if err != nil {
		log.Fatal(err)
	}
	m.contents[slotStandard] = out
	m.durations[slotStandard] = time.Since(t0)

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
	out2, err := r2.Render(m.md)
	if err != nil {
		log.Fatal(err)
	}
	m.contents[slotMinimal] = out2
	m.durations[slotMinimal] = time.Since(t0)

	t0 = time.Now()
	out3, err := gruff.Render(m.md,
		gruff.WithWordWrap(halfW-2),
	)
	if err != nil {
		log.Fatal(err)
	}
	m.contents[slotGruff] = out3
	m.durations[slotGruff] = time.Since(t0)

	m.leftView.SetContent(m.contents[m.slotAt(0)])
	m.topView.SetContent(m.contents[m.slotAt(1)])
	m.bottomView.SetContent(m.contents[m.slotAt(2)])
	return m
}

func (m model) headerFor(pos int) string {
	s := m.slotAt(pos)
	colors := slotColors[s]
	active := (pos == 0 && m.focus == focusLeft) ||
		(pos == 1 && m.focus == focusTop) ||
		(pos == 2 && m.focus == focusBottom)
	halfW := ((m.termWidth - 4) / 2) - 2
	return makeHeader(slotNames[s], m.wxhInfo()+"  "+m.durations[s].Round(time.Microsecond).String(),
		active, colors.activeBg, colors.inactiveBg, halfW)
}

func (m model) View() string {
	width := m.termWidth
	if width == 0 {
		width = 80
	}

	leftPane := m.paneBorder(m.focus == focusLeft, slotFocusColor[focusLeft]).Render(
		m.headerFor(0) + "\n" + m.leftView.View(),
	)

	topPane := m.paneBorder(m.focus == focusTop, slotFocusColor[focusTop]).Render(
		m.headerFor(1) + "\n" + m.topView.View(),
	)

	bottomPane := m.paneBorder(m.focus == focusBottom, slotFocusColor[focusBottom]).Render(
		m.headerFor(2) + "\n" + m.bottomView.View(),
	)

	rightSide := lipgloss.JoinVertical(lipgloss.Top, topPane, bottomPane)
	joined := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightSide)

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
