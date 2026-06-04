package component

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestNew(t *testing.T) {
	v := NewViewportWithScrollbar(80, 24)
	if v.Width != 80 {
		t.Fatalf("Width = %d, want 80", v.Width)
	}
	if v.Height != 24 {
		t.Fatalf("Height = %d, want 24", v.Height)
	}
}

func TestSetContent(t *testing.T) {
	v := NewViewportWithScrollbar(40, 10)
	v.SetContent("line1\nline2\nline3")
	if v.TotalLineCount() < 3 {
		t.Fatalf("TotalLineCount = %d, want >= 3", v.TotalLineCount())
	}
}

func TestRenderContentFits(t *testing.T) {
	v := NewViewportWithScrollbar(10, 3)
	v.SetContent("hello\nworld")
	out := v.View()
	lines := splitLines(out)
	if len(lines) != 3 {
		t.Fatalf("Render produced %d lines, want 3", len(lines))
	}
}

func TestRenderEmpty(t *testing.T) {
	v := NewViewportWithScrollbar(10, 3)
	out := v.View()
	lines := splitLines(out)
	if len(lines) != 3 {
		t.Fatalf("Render empty produced %d lines, want 3", len(lines))
	}
}

func TestRenderZeroSize(t *testing.T) {
	v := NewViewportWithScrollbar(0, 0)
	out := v.View()
	if out != "" {
		t.Fatalf("Render zero size should return empty, got %q", out)
	}
}

func TestScrollDown(t *testing.T) {
	v := NewViewportWithScrollbar(10, 3)
	v.SetContent("line0\nline1\nline2\nline3\nline4\nline5")
	v.ScrollDown(2)
	if v.YOffset() != 2 {
		t.Fatalf("after ScrollDown(2) YOffset = %d, want 2", v.YOffset())
	}
}

func TestScrollUp(t *testing.T) {
	v := NewViewportWithScrollbar(10, 3)
	v.SetContent("line0\nline1\nline2\nline3\nline4\nline5")
	v.ScrollDown(3)
	if v.YOffset() != 3 {
		t.Fatalf("ScrollDown(3) YOffset = %d, want 3", v.YOffset())
	}
	v.ScrollUp(1)
	if v.YOffset() != 2 {
		t.Fatalf("after ScrollUp(1) YOffset = %d, want 2", v.YOffset())
	}
}

func TestSetYOffset(t *testing.T) {
	v := NewViewportWithScrollbar(10, 3)
	v.SetContent("line0\nline1\nline2\nline3\nline4\nline5")
	v.SetYOffset(2)
	if v.YOffset() != 2 {
		t.Fatalf("SetYOffset(2) YOffset = %d, want 2", v.YOffset())
	}
}

func TestSetYOffsetClamp(t *testing.T) {
	v := NewViewportWithScrollbar(10, 3)
	v.SetContent("a\nb\nc")
	v.SetYOffset(10)
	if v.YOffset() != 0 {
		t.Fatalf("SetYOffset(10) clamped to %d, want 0", v.YOffset())
	}
}

func TestSetYOffsetNegative(t *testing.T) {
	v := NewViewportWithScrollbar(10, 3)
	v.SetContent("line0\nline1\nline2\nline3\nline4\nline5")
	v.SetYOffset(-5)
	if v.YOffset() != 0 {
		t.Fatalf("SetYOffset(-5) clamped to %d, want 0", v.YOffset())
	}
}

func TestScrollPercent(t *testing.T) {
	v := NewViewportWithScrollbar(10, 3)
	v.SetContent("line0\nline1\nline2\nline3\nline4\nline5")
	if v.ScrollPercent() != 0 {
		t.Fatalf("initial ScrollPercent = %v, want 0", v.ScrollPercent())
	}
	v.ScrollDown(3)
	if v.ScrollPercent() <= 0 {
		t.Fatalf("after scroll ScrollPercent = %v, want > 0", v.ScrollPercent())
	}
}

func TestMouseClickOnScrollbar(t *testing.T) {
	v := NewViewportWithScrollbar(10, 4)
	v.OriginX = 0
	v.OriginY = 0
	v.SetContent(strings.Repeat("line\n", 20))

	click := tea.MouseClickMsg{X: 9, Y: 3, Button: tea.MouseLeft}
	v.Update(click)

	if v.YOffset() <= 0 {
		t.Fatalf("after click on scrollbar track, YOffset = %d, want > 0", v.YOffset())
	}
}

func TestMouseClickOnContent(t *testing.T) {
	v := NewViewportWithScrollbar(10, 4)
	v.OriginX = 0
	v.OriginY = 0
	v.SetContent(strings.Repeat("line\n", 20))

	oldY := v.YOffset()
	click := tea.MouseClickMsg{X: 0, Y: 0, Button: tea.MouseLeft}
	v.Update(click)

	if v.YOffset() != oldY {
		t.Fatalf("click on content should not scroll, YOffset changed from %d to %d", oldY, v.YOffset())
	}
}

func TestMouseClickNonLeft(t *testing.T) {
	v := NewViewportWithScrollbar(10, 4)
	v.OriginX = 0
	v.OriginY = 0
	v.SetContent(strings.Repeat("line\n", 20))

	oldY := v.YOffset()
	click := tea.MouseClickMsg{X: 9, Y: 0, Button: tea.MouseRight}
	v.Update(click)

	if v.YOffset() != oldY {
		t.Fatalf("right click on scrollbar should not scroll")
	}
}

func TestMouseReleaseEndsDrag(t *testing.T) {
	v := NewViewportWithScrollbar(10, 4)
	v.OriginX = 0
	v.OriginY = 0
	v.SetContent(strings.Repeat("line\n", 20))

	click := tea.MouseClickMsg{X: 9, Y: 0, Button: tea.MouseLeft}
	v.Update(click)
	if !v.dragging {
		t.Fatalf("click on thumb should start drag")
	}

	release := tea.MouseReleaseMsg{X: 9, Y: 0}
	v.Update(release)
	if v.dragging {
		t.Fatalf("after release, dragging should be false")
	}
}

func TestMouseDrag(t *testing.T) {
	v := NewViewportWithScrollbar(10, 6)
	v.OriginX = 0
	v.OriginY = 0
	v.SetContent(strings.Repeat("line\n", 30))

	thumbPos, _ := v.calcThumb()
	click := tea.MouseClickMsg{X: 9, Y: thumbPos, Button: tea.MouseLeft}
	v.Update(click)
	if !v.dragging {
		t.Fatalf("click on thumb should start drag")
	}

	oldY := v.YOffset()
	v.Update(tea.MouseMotionMsg{X: 9, Y: thumbPos + 2})
	if v.YOffset() == oldY {
		t.Fatalf("drag should change YOffset, old=%d new=%d", oldY, v.YOffset())
	}

	v.Update(tea.MouseReleaseMsg{})
	if v.dragging {
		t.Fatalf("dragging should be false after release")
	}
}

func TestMouseClickOutsideViewport(t *testing.T) {
	v := NewViewportWithScrollbar(10, 4)
	v.OriginX = 10
	v.OriginY = 10
	v.SetContent(strings.Repeat("line\n", 20))

	oldY := v.YOffset()
	click := tea.MouseClickMsg{X: 0, Y: 0, Button: tea.MouseLeft}
	v.Update(click)

	if v.YOffset() != oldY {
		t.Fatalf("click outside viewport should not scroll")
	}
}

func TestCalcThumbContentFits(t *testing.T) {
	v := NewViewportWithScrollbar(10, 10)
	v.SetContent("hello")
	pos, size := v.calcThumb()
	if pos != 0 || size != 0 {
		t.Fatalf("calcThumb when content fits = (%d,%d), want (0,0)", pos, size)
	}
}

func TestSetContentUpdatesRender(t *testing.T) {
	v := NewViewportWithScrollbar(10, 3)
	v.SetContent("old")
	out1 := v.View()
	if !strings.Contains(out1, "old") {
		t.Fatal("first render should contain 'old'")
	}

	v.SetContent("new content line")
	out2 := v.View()
	if !strings.Contains(out2, "new") {
		t.Fatal("second render should contain 'new'")
	}
}

func TestRenderFooterScrollbar(t *testing.T) {
	v := NewViewportWithScrollbar(5, 2)
	v.SetContent("a\nb\nc\nd\ne")
	out := v.View()
	lines := splitLines(out)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestDisplayWidth(t *testing.T) {
	if w := displayWidth("hello"); w != 5 {
		t.Fatalf("displayWidth('hello') = %d, want 5", w)
	}
	if w := displayWidth("中文"); w != 4 {
		t.Fatalf("displayWidth('中文') = %d, want 4", w)
	}
}

func TestStripANSI(t *testing.T) {
	in := "\x1b[31mhello\x1b[39m"
	out := stripANSI(in)
	if out != "hello" {
		t.Fatalf("stripANSI(%q) = %q, want 'hello'", in, out)
	}
}

func TestTruncateToWidth(t *testing.T) {
	s := truncateToWidth("hello world", 5)
	if displayWidth(stripANSI(s)) != 5 {
		t.Fatalf("truncateToWidth('hello world', 5) width = %d, want 5", displayWidth(stripANSI(s)))
	}
}

func TestTruncateToWidthShorter(t *testing.T) {
	s := truncateToWidth("hi", 10)
	if s != "hi" {
		t.Fatalf("truncateToWidth('hi', 10) = %q, want 'hi'", s)
	}
}

func TestSplitLines(t *testing.T) {
	lines := splitLines("a\nb\nc")
	if len(lines) != 3 {
		t.Fatalf("splitLines('a\\nb\\nc') = %d lines, want 3", len(lines))
	}

	lines = splitLines("a\nb\nc\n")
	if len(lines) != 3 {
		t.Fatalf("splitLines('a\\nb\\nc\\n') = %d lines, want 3", len(lines))
	}
}


