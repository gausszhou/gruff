package gruff

import (
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/yuin/goldmark/ast"
	extensionAst "github.com/yuin/goldmark/extension/ast"
)

type nodeRenderer struct {
	buf      strings.Builder
	source   []byte
	th       Theme
	wordWrap int
}

func renderMarkdown(source []byte, th Theme, wordWrap int, node ast.Node) string {
	var r nodeRenderer
	r.source = source
	r.th = th
	r.wordWrap = wordWrap
	r.renderNode(node)
	return r.buf.String()
}

func (r *nodeRenderer) renderNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.Document:
		if r.th.Background != "" {
			r.buf.WriteString(string(ansiBg(r.th.Background)))
		}
		r.renderChildren(n)

	case *ast.Paragraph:
		r.renderChildren(n)
		if !r.isInsideList(n) && !r.isInsideTable(n) {
			r.buf.WriteString("\n\n")
		}

	case *ast.Heading:
		st := r.headingStyle(n.Level)
		r.buf.WriteString(string(st.start()))
		r.renderChildren(n)
		r.buf.WriteString(string(st.end(r.th.Background)))
		r.buf.WriteString("\n\n")

	case *ast.List:
		r.renderChildren(n)
		r.buf.WriteByte('\n')

	case *ast.ListItem:
		r.renderListItem(n)

	case *ast.Text:
		r.buf.Write(n.Value(r.source))
		if n.SoftLineBreak() {
			r.buf.WriteByte('\n')
		}

	case *ast.String:
		r.buf.Write(n.Value)

	case *ast.Emphasis:
		st := r.th.Em
		if n.Level == 2 {
			st = r.th.Strong
		}
		r.buf.WriteString(string(st.start()))
		r.renderChildren(n)
		r.buf.WriteString(string(st.end(r.th.Background)))

	case *ast.CodeSpan:
		r.buf.WriteString(string(r.th.Code.start()))
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if text, ok := c.(*ast.Text); ok {
				r.buf.Write(text.Value(r.source))
			}
		}
		r.buf.WriteString(string(r.th.Code.end(r.th.Background)))

	case *ast.Link:
		st := r.th.Link
		r.buf.WriteString(string(st.start()))
		r.renderChildren(n)
		r.buf.WriteString(string(st.end(r.th.Background)))
		if len(n.Destination) > 0 {
			url := string(n.Destination)
			uSt := r.th.LinkURL
			r.buf.WriteByte(' ')
			r.buf.WriteString(string(uSt.start()))
			r.buf.WriteByte('(')
			r.buf.WriteString(url)
			r.buf.WriteByte(')')
			r.buf.WriteString(string(uSt.end(r.th.Background)))
		}

	case *ast.Image:
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if text, ok := c.(*ast.Text); ok {
				r.buf.Write(text.Value(r.source))
			}
		}

	case *ast.ThematicBreak:
		r.buf.WriteString("\x1b[90m────────────────────\x1b[39m\n\n")

	case *extensionAst.Table:
		r.renderTable(n)

	default:
		r.renderChildren(n)
	}
}

func (r *nodeRenderer) renderChildren(node ast.Node) {
	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		r.renderNode(c)
	}
}

func (r *nodeRenderer) renderSubtree(node ast.Node) string {
	var sub nodeRenderer
	sub.source = r.source
	sub.th = r.th
	sub.renderChildren(node)
	return sub.buf.String()
}

func (r *nodeRenderer) renderListItem(node ast.Node) {
	parent := node.Parent()
	list, ok := parent.(*ast.List)
	if !ok {
		r.renderChildren(node)
		return
	}

	var index int
	for c := list.FirstChild(); c != nil; c = c.NextSibling() {
		if c == node {
			break
		}
		index++
	}

	if list.IsOrdered() {
		num := list.Start + index
		r.buf.WriteString("  ")
		r.buf.WriteString(string(r.th.Numbered.start()))
		r.buf.WriteString(itoa(num))
		r.buf.WriteString(". ")
		r.buf.WriteString(string(r.th.Numbered.end(r.th.Background)))
	} else {
		r.buf.WriteString("  ")
		r.buf.WriteString(string(r.th.Bullet.start()))
		r.buf.WriteString("• ")
		r.buf.WriteString(string(r.th.Bullet.end(r.th.Background)))
	}

	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		r.renderNode(c)
	}
	r.buf.WriteByte('\n')
}

func (r *nodeRenderer) isInsideList(node ast.Node) bool {
	for p := node.Parent(); p != nil; p = p.Parent() {
		if _, ok := p.(*ast.ListItem); ok {
			return true
		}
	}
	return false
}

func (r *nodeRenderer) isInsideTable(node ast.Node) bool {
	for p := node.Parent(); p != nil; p = p.Parent() {
		if _, ok := p.(*extensionAst.Table); ok {
			return true
		}
	}
	return false
}

func (r *nodeRenderer) headingStyle(level int) Style {
	switch level {
	case 1:
		return r.th.H1
	case 2:
		return r.th.H2
	case 3:
		return r.th.H3
	case 4:
		return r.th.H4
	case 5:
		return r.th.H5
	default:
		return r.th.H6
	}
}

type cellData struct {
	content string
	align   extensionAst.Alignment
	lines   []string
}

func wrapCellLines(content string, width int) []string {
	if displayWidth(stripANSI(content)) <= width {
		return []string{content}
	}

	var lines []string
	var line strings.Builder
	var word strings.Builder
	lineVisLen := 0
	wordVisLen := 0
	inAnsi := false

	flushWord := func() {
		w := word.String()
		word.Reset()
		if len(w) == 0 && wordVisLen == 0 {
			return
		}
		if lineVisLen > 0 && lineVisLen+1+wordVisLen > width {
			lines = append(lines, line.String())
			line.Reset()
			lineVisLen = 0
		}
		if lineVisLen > 0 && wordVisLen > 0 {
			line.WriteByte(' ')
			lineVisLen++
		}
		line.WriteString(w)
		lineVisLen += wordVisLen
		wordVisLen = 0
	}

	for _, r := range content {
		if inAnsi {
			word.WriteRune(r)
			if r == 'm' {
				inAnsi = false
			}
			continue
		}
		if r == '\x1b' {
			inAnsi = true
			word.WriteRune(r)
			continue
		}
		if r == ' ' || r == '\n' {
			flushWord()
			continue
		}
		if runewidth.RuneWidth(r) > 1 {
			flushWord()
			rw := runewidth.RuneWidth(r)
			if lineVisLen+rw > width && lineVisLen > 0 {
				lines = append(lines, line.String())
				line.Reset()
				lineVisLen = 0
			}
			line.WriteRune(r)
			lineVisLen += rw
			continue
		}
		word.WriteRune(r)
		wordVisLen++
	}
	flushWord()

	if line.Len() > 0 {
		lines = append(lines, line.String())
	} else if len(lines) == 0 {
		lines = []string{""}
	}

	return lines
}

func (r *nodeRenderer) renderTable(table *extensionAst.Table) {
	var headers []cellData
	var bodyRows [][]cellData

	for child := table.FirstChild(); child != nil; child = child.NextSibling() {
		switch n := child.(type) {
		case *extensionAst.TableHeader:
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				if cell, ok := c.(*extensionAst.TableCell); ok {
					headers = append(headers, cellData{
						content: r.renderSubtree(cell),
						align:   cell.Alignment,
					})
				}
			}
		case *extensionAst.TableRow:
			var rowCells []cellData
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				if cell, ok := c.(*extensionAst.TableCell); ok {
					rowCells = append(rowCells, cellData{
						content: r.renderSubtree(cell),
						align:   cell.Alignment,
					})
				}
			}
			bodyRows = append(bodyRows, rowCells)
		}
	}

	numCols := len(headers)
	if numCols == 0 && len(bodyRows) > 0 {
		numCols = len(bodyRows[0])
	}
	if numCols == 0 {
		return
	}

	colWidths := make([]int, numCols)
	colAligns := make([]extensionAst.Alignment, numCols)

	for i, h := range headers {
		if i >= numCols {
			break
		}
		w := displayWidth(stripANSI(h.content))
		if w > colWidths[i] {
			colWidths[i] = w
		}
		if h.align != 0 {
			colAligns[i] = h.align
		}
	}

	for _, row := range bodyRows {
		for i, cell := range row {
			if i >= numCols {
				break
			}
			w := displayWidth(stripANSI(cell.content))
			if w > colWidths[i] {
				colWidths[i] = w
			}
			if cell.align != 0 {
				colAligns[i] = cell.align
			}
		}
	}

	overhead := 3 * (numCols - 1)
	maxCol := (r.wordWrap - overhead) / numCols
	if maxCol < 20 {
		maxCol = 20
	}

	for i := range colWidths {
		if colWidths[i] < 3 {
			colWidths[i] = 3
		} else if colWidths[i] > maxCol {
			colWidths[i] = maxCol
		}
	}

	for i, h := range headers {
		headers[i].lines = wrapCellLines(h.content, colWidths[i])
	}
	for _, row := range bodyRows {
		for i, cell := range row {
			row[i].lines = wrapCellLines(cell.content, colWidths[i])
		}
	}

	border := "\x1b[38;5;8m"
	reset := "\x1b[39m"

	seg := func(w int) string {
		s := ""
		for j := 0; j < w+2; j++ {
			s += "\u2500"
		}
		return s
	}

	hline := func() {
		r.buf.WriteString(border)
		for i, w := range colWidths {
			r.buf.WriteString(seg(w))
			if i < len(colWidths)-1 {
				r.buf.WriteString("\u253c") // ┼
			}
		}
		r.buf.WriteString(reset)
		r.buf.WriteByte('\n')
	}

	if len(headers) > 0 {
		r.renderTableRow(headers, colWidths, colAligns)
		hline()
	}

	for i, row := range bodyRows {
		r.renderTableRow(row, colWidths, colAligns)
		if i < len(bodyRows)-1 {
			hline()
		}
	}
}

func (r *nodeRenderer) renderTableRow(cells []cellData, widths []int, aligns []extensionAst.Alignment) {
	maxLines := 1
	for _, cell := range cells {
		if len(cell.lines) > maxLines {
			maxLines = len(cell.lines)
		}
	}

	bgReset := "\x1b[49m"
	if r.th.Background != "" {
		bgReset = string(ansiBg(r.th.Background))
	}

	for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
		for i, cell := range cells {
			if i >= len(widths) {
				break
			}

			if i > 0 {
				r.buf.WriteString("\x1b[38;5;8m\u2502\x1b[39m") // │
			}

			var content string
			if lineIdx < len(cell.lines) {
				content = cell.lines[lineIdx]
			}

			visLen := displayWidth(stripANSI(content))
			padding := widths[i] - visLen
			if padding < 0 {
				padding = 0
			}

			r.buf.WriteByte(' ')
			switch aligns[i] {
			case extensionAst.AlignRight:
				for j := 0; j < padding; j++ {
					r.buf.WriteByte(' ')
				}
				r.buf.WriteString(content)
				r.buf.WriteString("\x1b[39m")
				r.buf.WriteString(bgReset)
			case extensionAst.AlignCenter:
				leftPad := padding / 2
				rightPad := padding - leftPad
				for j := 0; j < leftPad; j++ {
					r.buf.WriteByte(' ')
				}
				r.buf.WriteString(content)
				r.buf.WriteString("\x1b[39m")
				r.buf.WriteString(bgReset)
				for j := 0; j < rightPad; j++ {
					r.buf.WriteByte(' ')
				}
			default:
				r.buf.WriteString(content)
				r.buf.WriteString("\x1b[39m")
				r.buf.WriteString(bgReset)
				for j := 0; j < padding; j++ {
					r.buf.WriteByte(' ')
				}
			}
			r.buf.WriteByte(' ')
		}
		r.buf.WriteByte('\n')
	}
}
