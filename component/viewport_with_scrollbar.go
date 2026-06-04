package component

import (
	"image/color"
	"strings"
	"unicode/utf8"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mattn/go-runewidth"
)

const (
	ansiDefaultFg = "\x1b[39m"
	ansiNoBold    = "\x1b[22m"
	ansiNoItalic  = "\x1b[23m"

	scrollbarWidth = 1
)

	type ViewportWithScrollbar struct {
	inner viewport.Model

	OriginX int
	OriginY int
	Width   int
	Height  int

	TrackColor color.Color
	ThumbColor color.Color

	dragging   bool
	dragOffset int

	trackStyle lipgloss.Style
	thumbStyle lipgloss.Style
}

func NewViewportWithScrollbar(width, height int) ViewportWithScrollbar {
	contentWidth := width - scrollbarWidth
	if contentWidth < 0 {
		contentWidth = 0
	}
	m := ViewportWithScrollbar{
		inner: viewport.New(
			viewport.WithWidth(contentWidth),
			viewport.WithHeight(height),
		),
		Width:      width,
		Height:     height,
		TrackColor: lipgloss.Color("236"),
		ThumbColor: lipgloss.Color("248"),
	}
	m.inner.FillHeight = true
	m.rebuildStyles()
	return m
}

func (v *ViewportWithScrollbar) rebuildStyles() {
	v.trackStyle = lipgloss.NewStyle().Background(v.TrackColor)
	v.thumbStyle = lipgloss.NewStyle().Background(v.ThumbColor)
}

func (v ViewportWithScrollbar) View() string {
	if v.Height <= 0 || v.Width <= 0 {
		return ""
	}

	contentWidth := v.Width - scrollbarWidth
	if contentWidth < 0 {
		contentWidth = 0
	}

	innerContent := v.inner.View()
	inLines := splitLines(innerContent)
	thumbPos, thumbSize := v.calcThumb()
	hasScrollbar := v.inner.TotalLineCount() > v.Height

	var buf strings.Builder
	for i := 0; i < v.Height; i++ {
		if i < len(inLines) {
			line := inLines[i]
			rendered := truncateToWidth(line, contentWidth)
			buf.WriteString(rendered)
			w := displayWidth(stripANSI(rendered))
			if w < contentWidth {
				buf.WriteString(strings.Repeat(" ", contentWidth-w))
			}
		} else {
			buf.WriteString(strings.Repeat(" ", contentWidth))
		}

		buf.WriteString(ansiNoBold + ansiNoItalic + ansiDefaultFg)
		if !hasScrollbar {
			buf.WriteString(v.trackStyle.Render(" "))
		} else if i >= thumbPos && i < thumbPos+thumbSize {
			buf.WriteString(v.thumbStyle.Render(" "))
		} else {
			buf.WriteString(v.trackStyle.Render(" "))
		}

		if i < v.Height-1 {
			buf.WriteByte('\n')
		}
	}

	return buf.String()
}

func (v *ViewportWithScrollbar) Update(msg tea.Msg) (ViewportWithScrollbar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseClickMsg:
		m := msg.Mouse()
		if v.isScrollbarX(m.X) && m.Button == tea.MouseLeft {
			v.handleScrollbarClick(m.Y)
		}
		return *v, nil

	case tea.MouseMotionMsg:
		if v.dragging {
			m := msg.Mouse()
			v.handleScrollbarDrag(m.Y)
		}
		return *v, nil

	case tea.MouseReleaseMsg:
		v.dragging = false
		return *v, nil

	default:
		var cmd tea.Cmd
		v.inner, cmd = v.inner.Update(msg)
		return *v, cmd
	}
}

func (v *ViewportWithScrollbar) isScrollbarX(x int) bool {
	contentWidth := v.Width - scrollbarWidth
	scrollbarStart := v.OriginX + contentWidth
	scrollbarEnd := v.OriginX + v.Width
	return x >= scrollbarStart && x < scrollbarEnd
}

func (v *ViewportWithScrollbar) handleScrollbarClick(termY int) {
	relY := termY - v.OriginY
	total := v.inner.TotalLineCount()

	if total <= v.Height {
		return
	}

	thumbPos, thumbSize := v.calcThumb()
	if relY >= thumbPos && relY < thumbPos+thumbSize {
		v.dragging = true
		v.dragOffset = relY - thumbPos
		return
	}

	maxScroll := total - v.Height
	if maxScroll <= 0 {
		return
	}
	newScrollY := int(float64(relY) / float64(v.Height) * float64(total))
	if newScrollY > maxScroll {
		newScrollY = maxScroll
	}
	if newScrollY < 0 {
		newScrollY = 0
	}
	if newScrollY != v.inner.YOffset() {
		v.inner.SetYOffset(newScrollY)
	}
}

func (v *ViewportWithScrollbar) handleScrollbarDrag(termY int) {
	total := v.inner.TotalLineCount()
	if total <= v.Height {
		v.dragging = false
		return
	}

	_, thumbSize := v.calcThumb()
	trackHeight := v.Height - thumbSize
	if trackHeight <= 0 {
		v.dragging = false
		return
	}

	relY := termY - v.OriginY - v.dragOffset
	if relY < 0 {
		relY = 0
	}
	if relY > trackHeight {
		relY = trackHeight
	}

	maxScroll := total - v.Height
	newScrollY := int(float64(relY) / float64(trackHeight) * float64(maxScroll))
	if newScrollY > maxScroll {
		newScrollY = maxScroll
	}
	if newScrollY < 0 {
		newScrollY = 0
	}
	if newScrollY != v.inner.YOffset() {
		v.inner.SetYOffset(newScrollY)
	}
}

func (v *ViewportWithScrollbar) calcThumb() (pos, size int) {
	total := v.inner.TotalLineCount()
	height := v.Height
	if total <= height || height <= 0 {
		return 0, 0
	}

	visibleRatio := float64(height) / float64(total)
	size = int(visibleRatio * float64(height))
	if size < 1 {
		size = 1
	}

	maxScroll := total - height
	scrollRatio := float64(v.inner.YOffset()) / float64(maxScroll)
	trackHeight := height - size
	pos = int(scrollRatio * float64(trackHeight))
	if pos < 0 {
		pos = 0
	}
	if pos > trackHeight {
		pos = trackHeight
	}
	return pos, size
}

func (v *ViewportWithScrollbar) IsDragging() bool {
	return v.dragging
}

func (v *ViewportWithScrollbar) SetContent(s string) {
	v.inner.SetContent(s)
}

func (v *ViewportWithScrollbar) ScrollPercent() float64 {
	return v.inner.ScrollPercent()
}

func (v *ViewportWithScrollbar) TotalLineCount() int {
	return v.inner.TotalLineCount()
}

func (v *ViewportWithScrollbar) YOffset() int {
	return v.inner.YOffset()
}

func (v *ViewportWithScrollbar) SetYOffset(n int) {
	v.inner.SetYOffset(n)
}

func (v *ViewportWithScrollbar) ScrollDown(n int) {
	v.inner.ScrollDown(n)
}

func (v *ViewportWithScrollbar) ScrollUp(n int) {
	v.inner.ScrollUp(n)
}

func (v *ViewportWithScrollbar) Inner() *viewport.Model {
	return &v.inner
}

func splitLines(s string) []string {
	lines := strings.Split(s, "\n")
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func displayWidth(s string) int {
	w := 0
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == 0xFE0F {
			i += size
			continue
		}
		if i+size < len(s) {
			next, nextSize := utf8.DecodeRuneInString(s[i+size:])
			if next == 0xFE0F {
				w += 2
				i += size + nextSize
				continue
			}
		}
		w += runewidth.RuneWidth(r)
		i += size
	}
	return w
}

func stripANSI(s string) string {
	var out []byte
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			for j := i + 2; j < len(s); j++ {
				if s[j] >= 0x40 && s[j] <= 0x7E {
					i = j
					break
				}
			}
			continue
		}
		out = append(out, s[i])
	}
	return string(out)
}

func truncateToWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if displayWidth(stripANSI(s)) <= width {
		return s
	}
	var out strings.Builder
	w := 0
	for i := 0; i < len(s); {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && !(s[j] >= 0x40 && s[j] <= 0x7E) {
				j++
			}
			if j < len(s) {
				j++
			}
			out.WriteString(s[i:j])
			i = j
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		rw := runewidth.RuneWidth(r)
		if w+rw > width {
			break
		}
		out.WriteRune(r)
		w += rw
		i += size
	}
	return out.String()
}
