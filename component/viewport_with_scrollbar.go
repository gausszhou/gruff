package component

import (
	"image/color"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
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

	scrollbar := v.renderScrollbar()
	return lipgloss.JoinHorizontal(lipgloss.Top, v.inner.View(), scrollbar)
}

func (v ViewportWithScrollbar) renderScrollbar() string {
	thumbPos, thumbSize := v.calcThumb()
	hasScrollbar := v.inner.TotalLineCount() > v.Height

	var buf strings.Builder
	for i := 0; i < v.Height; i++ {
		if i > 0 {
			buf.WriteByte('\n')
		}
		if !hasScrollbar {
			buf.WriteString(v.trackStyle.Render(" "))
		} else if i >= thumbPos && i < thumbPos+thumbSize {
			buf.WriteString(v.thumbStyle.Render(" "))
		} else {
			buf.WriteString(v.trackStyle.Render(" "))
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


