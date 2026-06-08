package gruff

import (
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/yuin/goldmark/ast"
	extensionAst "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

type nodeRenderer struct {
	buf          strings.Builder
	source       []byte
	th           Theme
	wordWrap     int
	padding      int
	inBlockquote bool
}

func renderMarkdown(source []byte, th Theme, wordWrap int, padding int, node ast.Node) string {
	var r nodeRenderer
	r.source = source
	r.th = th
	r.wordWrap = wordWrap
	r.padding = padding
	r.renderNode(node)
	return r.buf.String()
}

func (r *nodeRenderer) renderNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.Document:
		r.buf.WriteString(string(ansiBg(r.th.Bg)))
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if c != n.FirstChild() && isBlockLevel(c) {
				r.buf.WriteByte('\n')
			}
			r.renderNode(c)
		}

	case *ast.Paragraph:
		r.renderChildren(n)
		r.buf.WriteByte('\n')

	case *ast.TextBlock:
		r.renderChildren(n)
		r.buf.WriteByte('\n')

	case *ast.Heading:
		st := r.headingStyle(n.Level)
		r.buf.WriteString(string(st.start()))
		r.renderChildren(n)
		r.buf.WriteString(string(st.end()))
		r.buf.WriteByte('\n')

	case *ast.List:
		r.renderChildren(n)

	case *ast.ListItem:
		r.renderListItem(n)

	case *ast.Text:
		_, isPara := n.Parent().(*ast.Paragraph)
		_, isTB := n.Parent().(*ast.TextBlock)
		_, isTC := n.Parent().(*extensionAst.TableCell)
		if isPara || isTB || isTC {
			r.buf.WriteString(string(r.th.Paragraph.start()))
		}
		v := n.Value(r.source)
		r.buf.Write(v)
		if n.SoftLineBreak() {
			r.buf.WriteByte(' ')
		}
		if isPara || isTB || isTC {
			r.buf.WriteString(string(r.th.Paragraph.end()))
		}

	case *ast.String:
		_, isPara := n.Parent().(*ast.Paragraph)
		_, isTB := n.Parent().(*ast.TextBlock)
		_, isTC := n.Parent().(*extensionAst.TableCell)
		if isPara || isTB || isTC {
			r.buf.WriteString(string(r.th.Paragraph.start()))
		}
		r.buf.Write(n.Value)
		if isPara || isTB || isTC {
			r.buf.WriteString(string(r.th.Paragraph.end()))
		}

	case *ast.Emphasis:
		st := r.th.Em
		if n.Level == 2 {
			st = r.th.Strong
		}
		r.buf.WriteString(string(st.start()))
		r.renderChildren(n)
		r.buf.WriteString(string(st.end()))

	case *ast.CodeSpan:
		r.buf.WriteString(string(r.th.Code.start()))
		for range r.th.Code.Padding {
			r.buf.WriteByte(' ')
		}
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if text, ok := c.(*ast.Text); ok {
				r.buf.Write(text.Value(r.source))
			}
		}
		for range r.th.Code.Padding {
			r.buf.WriteByte(' ')
		}
		r.buf.WriteString(string(r.th.Code.end()))

	case *ast.Link:
		st := r.th.Link
		r.buf.WriteString(string(st.start()))
		r.renderChildren(n)
		r.buf.WriteString(string(st.end()))
		if len(n.Destination) > 0 {
			url := string(n.Destination)
			uSt := r.th.LinkURL
			r.buf.WriteByte(' ')
			r.buf.WriteString(string(uSt.start()))
			r.buf.WriteByte('(')
			r.buf.WriteString(url)
			r.buf.WriteByte(')')
			r.buf.WriteString(string(uSt.end()))
		}

	case *ast.AutoLink:
		st := r.th.Link
		r.buf.WriteString(string(st.start()))
		r.buf.Write(n.Label(r.source))
		r.buf.WriteString(string(st.end()))

	case *ast.Image:
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if text, ok := c.(*ast.Text); ok {
				r.buf.Write(text.Value(r.source))
			}
		}

	case *ast.FencedCodeBlock:
		r.renderCodeBlock(n.Lines(), n.Language(r.source))

	case *ast.CodeBlock:
		r.renderCodeBlock(n.Lines(), nil)

	case *ast.ThematicBreak:
		r.buf.WriteString(string(r.th.Hr.start()))
		r.buf.WriteString("────────────────────")
		r.buf.WriteString(string(r.th.Hr.end()))
		r.buf.WriteByte('\n')

	case *ast.Blockquote:
		st := r.th.BlockQuote
		prefix := string(st.start()) + "│ " + string(st.end())
		r.inBlockquote = true
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			r.buf.WriteString(prefix)
			r.renderNode(c)
			if c.NextSibling() != nil {
				r.buf.WriteString(prefix)
				r.buf.WriteByte('\n')
			}
		}
		r.inBlockquote = false

	case *extensionAst.Table:
		r.renderTable(n)

	case *extensionAst.TaskCheckBox:
		if n.IsChecked {
			r.buf.WriteString(string(r.th.TaskChecked.start()))
			r.buf.WriteString("[\u2713]")
			r.buf.WriteString(string(r.th.TaskChecked.end()))
			r.buf.WriteByte(' ')
		} else {
			r.buf.WriteString(string(r.th.TaskUnchecked.start()))
			r.buf.WriteString("[ ]")
			r.buf.WriteString(string(r.th.TaskUnchecked.end()))
			r.buf.WriteByte(' ')
		}

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

func (r *nodeRenderer) renderCodeBlock(lines *text.Segments, lang []byte) {
	st := r.th.Code
	codeStyleStart := string(st.start())
	const padding = 0

	for j := 0; j < padding; j++ {
		r.buf.WriteByte(' ')
	}
	r.buf.WriteString("```")
	ls := r.th.LinkURL
	r.buf.WriteString(string(ls.start()))
	if len(lang) > 0 {
		r.buf.Write(lang)
	}
	r.buf.WriteString(string(ls.end()))
	r.buf.WriteByte('\n')

	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)
		content := string(seg.Value(r.source))
		content = strings.TrimRight(content, "\n\r")
		r.buf.WriteString(codeStyleStart)
		for j := 0; j < padding; j++ {
			r.buf.WriteByte(' ')
		}
		r.buf.WriteString(content)
		r.buf.WriteString(string(st.end()))
		r.buf.WriteByte('\n')
	}
	for j := 0; j < padding; j++ {
		r.buf.WriteByte(' ')
	}
	r.buf.WriteString("```")
	r.buf.WriteByte('\n')
}

func (r *nodeRenderer) listDepth(node ast.Node) int {
	depth := 0
	for p := node.Parent(); p != nil; p = p.Parent() {
		if _, ok := p.(*ast.List); ok {
			depth++
		}
	}
	return depth
}

func (r *nodeRenderer) renderListItem(node ast.Node) {
	depth := r.listDepth(node)
	indent := strings.Repeat("  ", depth)

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

	isTask := r.isTaskItem(node)
	if isTask {
		// no thing
	} else if list.IsOrdered() {
		num := list.Start + index
		r.buf.WriteString(indent)
		r.buf.WriteString(itoa(num))
		r.buf.WriteString(". ")
	} else {
		r.buf.WriteString(indent)
		r.buf.WriteString("• ")
	}

	r.renderChildren(node)
}

func (r *nodeRenderer) isTaskItem(node ast.Node) bool {
	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		if _, ok := c.(*extensionAst.TaskCheckBox); ok {
			return true
		}
		if tb, ok := c.(*ast.TextBlock); ok {
			for t := tb.FirstChild(); t != nil; t = t.NextSibling() {
				if _, ok := t.(*extensionAst.TaskCheckBox); ok {
					return true
				}
			}
		}
	}
	return false
}

func isBlockLevel(n ast.Node) bool {
	switch n.(type) {
	case *ast.Paragraph, *ast.Heading, *ast.List, *ast.Blockquote,
		*ast.FencedCodeBlock, *ast.CodeBlock, *ast.ThematicBreak,
		*extensionAst.Table:
		return true
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
	word := make([]byte, 0, 64)
	lineVisLen := 0
	wordVisLen := 0
	inAnsi := false

	flushWord := func() {
		if len(word) == 0 && wordVisLen == 0 {
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
		line.Write(word)
		lineVisLen += wordVisLen
		wordVisLen = 0
		word = word[:0]
	}

	for _, r := range content {
		if inAnsi {
			word = utf8.AppendRune(word, r)
			if r == 'm' {
				inAnsi = false
			}
			continue
		}
		if r == '\x1b' {
			inAnsi = true
			word = utf8.AppendRune(word, r)
			continue
		}
		if r == ' ' {
			flushWord()
			continue
		}
		if r == '\n' {
			flushWord()
			if line.Len() > 0 {
				lines = append(lines, line.String())
				line.Reset()
				lineVisLen = 0
			} else {
				lines = append(lines, "")
			}
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
		word = utf8.AppendRune(word, r)
		wordVisLen += runewidth.RuneWidth(r)
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

	overhead := 3*(numCols-1) + 2

	totalNat := 0
	for _, w := range colWidths {
		totalNat += w
	}
	totalNat += overhead

	maxWidth := r.wordWrap - r.padding
	if totalNat <= maxWidth {
		for i := range colWidths {
			if colWidths[i] < 3 {
				colWidths[i] = 3
			}
		}
	} else {
		equal := (maxWidth - overhead) / numCols
		if equal < 3 {
			equal = 3
		}
		for i := range colWidths {
			colWidths[i] = equal
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

	border := string(r.th.Border.start())
	reset := string(r.th.Border.end())

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

	for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
		for i, cell := range cells {
			if i >= len(widths) {
				break
			}

			if i > 0 {
				r.buf.WriteString(string(r.th.Border.start()))
				r.buf.WriteString("\u2502")
				r.buf.WriteString(string(r.th.Border.end()))
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
			case extensionAst.AlignCenter:
				leftPad := padding / 2
				rightPad := padding - leftPad
				for j := 0; j < leftPad; j++ {
					r.buf.WriteByte(' ')
				}
				r.buf.WriteString(content)
				r.buf.WriteString("\x1b[39m")
				for j := 0; j < rightPad; j++ {
					r.buf.WriteByte(' ')
				}
			default:
				r.buf.WriteString(content)
				r.buf.WriteString("\x1b[39m")
				for j := 0; j < padding; j++ {
					r.buf.WriteByte(' ')
				}
			}
			r.buf.WriteByte(' ')
		}
		r.buf.WriteByte('\n')
	}
}
