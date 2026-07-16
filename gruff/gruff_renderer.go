package gruff

import (
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/yuin/goldmark/ast"
	extensionAst "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

func hasOSC8String(s string) bool {
	return len(s) >= 4 && s[0] == '\x1b' && s[1] == ']' && s[2] == '8' && s[3] == ';'
}

func endsWithOSC8End(s string) bool {
	return len(s) >= 7 && s[len(s)-7] == '\x1b' && s[len(s)-6] == ']' && s[len(s)-5] == '8' && s[len(s)-4] == ';' && s[len(s)-3] == ';' && s[len(s)-2:] == "\x1b\\"
}

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
		r.renderDocument(n)

	case *ast.Paragraph:
		r.renderChildrenMerged(n)
		r.buf.WriteByte('\n')

	case *ast.TextBlock:
		r.renderChildrenMerged(n)
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
		r.renderTextLike(n.Parent(), n.Value(r.source), n.SoftLineBreak())

	case *ast.String:
		r.renderTextLike(n.Parent(), n.Value, false)

	case *ast.Emphasis:
		r.renderEmphasis(n)

	case *ast.CodeSpan:
		r.renderCodeSpan(n)

	case *ast.Link:
		r.renderLink(n)

	case *ast.AutoLink:
		r.renderAutoLink(n)

	case *ast.Image:
		r.renderImage(n)

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
		r.renderBlockquote(n)

	case *extensionAst.Table:
		r.renderTable(n)

	case *extensionAst.TaskCheckBox:
		r.renderTaskCheckBox(n)

	default:
		r.renderChildren(n)
	}
}

func (r *nodeRenderer) renderTextLike(parent ast.Node, value []byte, softBreak bool) {
	_, isPara := parent.(*ast.Paragraph)
	_, isTB := parent.(*ast.TextBlock)
	_, isTC := parent.(*extensionAst.TableCell)
	if isPara || isTB || isTC {
		r.buf.WriteString(string(r.th.Paragraph.start()))
	}
	r.buf.Write(value)
	if softBreak {
		r.buf.WriteByte(' ')
	}
	if isPara || isTB || isTC {
		r.buf.WriteString(string(r.th.Paragraph.end()))
	}
}

func (r *nodeRenderer) renderEmphasis(n *ast.Emphasis) {
	st := r.th.Em
	if n.Level == 2 {
		st = r.th.Strong
	}
	r.buf.WriteString(string(st.start()))
	r.renderChildren(n)
	r.buf.WriteString(string(st.end()))
}

func (r *nodeRenderer) renderCodeSpan(n *ast.CodeSpan) {
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
}

func (r *nodeRenderer) renderLink(n *ast.Link) {
	st := r.th.Link
	url := string(n.Destination)
	r.buf.WriteString(osc8Link(url))
	r.buf.WriteString(string(st.start()))
	r.renderChildren(n)
	r.buf.WriteString(string(st.end()))
	if len(n.Destination) > 0 {
		r.buf.WriteByte(' ')
		uSt := r.th.LinkURL
		r.buf.WriteString(string(uSt.start()))
		r.buf.WriteByte('(')
		r.buf.WriteString(url)
		r.buf.WriteByte(')')
		r.buf.WriteString(string(uSt.end()))
	}
	r.buf.WriteString(osc8End)
}

func (r *nodeRenderer) renderImage(n *ast.Image) {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if text, ok := c.(*ast.Text); ok {
			r.buf.Write(text.Value(r.source))
		}
	}
}

func (r *nodeRenderer) renderAutoLink(n *ast.AutoLink) {
	url := string(n.URL(r.source))
	r.buf.WriteString(osc8Link(url))
	st := r.th.LinkURL
	r.buf.WriteString(string(st.start()))
	r.buf.Write(n.Label(r.source))
	r.buf.WriteString(string(st.end()))
	r.buf.WriteString(osc8End)
}

func (r *nodeRenderer) renderBlockquote(n *ast.Blockquote) {
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
}

func (r *nodeRenderer) renderTaskCheckBox(n *extensionAst.TaskCheckBox) {
	if n.IsChecked {
		r.buf.WriteString(string(r.th.TaskChecked.start()))
		r.buf.WriteString("[\u2713]")
		r.buf.WriteString(string(r.th.TaskChecked.end()))
	} else {
		r.buf.WriteString(string(r.th.TaskUnchecked.start()))
		r.buf.WriteString("[ ]")
		r.buf.WriteString(string(r.th.TaskUnchecked.end()))
	}
	r.buf.WriteByte(' ')
}

func (r *nodeRenderer) renderDocument(n *ast.Document) {
	r.buf.WriteString(string(ansiBg(r.th.Bg)))
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if c != n.FirstChild() && isBlockLevel(c) {
			r.buf.WriteByte('\n')
		}
		r.renderNode(c)
	}
}

func (r *nodeRenderer) renderChildren(node ast.Node) {
	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		r.renderNode(c)
	}
}

func (r *nodeRenderer) renderChildrenMerged(node ast.Node) {
	for c := node.FirstChild(); c != nil; {
		if _, ok := c.(*ast.Text); ok {
			var buf []byte
			for ; c != nil; c = c.NextSibling() {
				if t, ok := c.(*ast.Text); ok {
					buf = append(buf, t.Value(r.source)...)
				} else {
					break
				}
			}
			r.renderTextLike(node, buf, false)
		} else {
			r.renderNode(c)
			c = c.NextSibling()
		}
	}
}

func (r *nodeRenderer) renderSubtree(node ast.Node) string {
	var sub nodeRenderer
	sub.source = r.source
	sub.th = r.th
	sub.renderChildrenMerged(node)
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
	activeStyle := make([]byte, 0, 64)
	wordStartStyle := make([]byte, 0, 64)
	lineStartStyle := make([]byte, 0, 64)
	tempEsc := make([]byte, 0, 32)
	word := make([]byte, 0, 64)
	lineVisLen := 0
	wordVisLen := 0
	escSt := escNone

	flushWord := func() {
		if len(word) == 0 && wordVisLen == 0 {
			return
		}
		if lineVisLen > 0 && lineVisLen+1+wordVisLen > width {
			lines = append(lines, line.String())
			line.Reset()
			line.Write(wordStartStyle)
			lineStartStyle = append(lineStartStyle[:0], wordStartStyle...)
			lineVisLen = 0
		}
		if lineVisLen == 0 && wordVisLen > width {
			remaining := word
			splitStyle := make([]byte, len(wordStartStyle))
			copy(splitStyle, wordStartStyle)
			line.Reset()

			for len(remaining) > 0 {
				chunkStyle := make([]byte, len(splitStyle))
				copy(chunkStyle, splitStyle)
				head, tail := breakChunk(remaining, width, &splitStyle)
				if len(head) == 0 {
					line.Write(chunkStyle)
					line.Write(remaining)
					lineVisLen += ansiDisplayWidth(remaining)
					break
				}
				line.Write(chunkStyle)
				line.Write(head)
				lineVisLen += ansiDisplayWidth(head)
				lineStartStyle = append(lineStartStyle[:0], chunkStyle...)
			if len(tail) > 0 {
					line.WriteString(osc8End)
					lines = append(lines, line.String())
					line.Reset()
					lineVisLen = 0
				}
				remaining = tail
			}

			wordVisLen = 0
			word = word[:0]
			wordStartStyle = append(wordStartStyle[:0], activeStyle...)
			if len(lineStartStyle) == 0 {
				lineStartStyle = append(lineStartStyle[:0], activeStyle...)
			}
			return
		}
		if lineVisLen > 0 && wordVisLen > 0 {
			line.WriteByte(' ')
			lineVisLen++
		}
		line.Write(word)
		lineVisLen += wordVisLen
		wordVisLen = 0
		word = word[:0]
		wordStartStyle = append(wordStartStyle[:0], activeStyle...)
		if len(lineStartStyle) == 0 {
			lineStartStyle = append(lineStartStyle[:0], activeStyle...)
		}
	}

	for _, r := range content {
		processCellRune(r, &escSt, &word, &wordVisLen, &line, &lineVisLen, &lines, width, flushWord, &activeStyle, &tempEsc, &lineStartStyle)
	}
	flushWord()

	if line.Len() > 0 {
		lines = append(lines, line.String())
	} else if len(lines) == 0 {
		lines = []string{""}
	}

	return lines
}

func processCellRune(r rune, escSt *escapeState, word *[]byte, wordVisLen *int, line *strings.Builder, lineVisLen *int, lines *[]string, width int, flushWord func(), activeStyle, tempEsc, lineStartStyle *[]byte) {
	*word = utf8.AppendRune(*word, r)
	*tempEsc = utf8.AppendRune(*tempEsc, r)
	switch *escSt {
	case escStart:
		if r == '[' {
			*escSt = escCSI
		} else if r == ']' {
			*escSt = escOSC
		} else {
			*escSt = escNone
			*tempEsc = (*tempEsc)[:0]
		}
	case escCSI:
		if r >= 0x40 && r <= 0x7E {
			*escSt = escNone
			updateActiveStyle(activeStyle, *tempEsc)
			*tempEsc = (*tempEsc)[:0]
		}
	case escOSC:
		if r == '\x1b' {
			*escSt = escOSCSt
		} else if r == '\x07' {
			*escSt = escNone
			updateActiveStyle(activeStyle, *tempEsc)
			*tempEsc = (*tempEsc)[:0]
		}
	case escOSCSt:
		if r == '\\' {
			*escSt = escNone
			updateActiveStyle(activeStyle, *tempEsc)
			*tempEsc = (*tempEsc)[:0]
		} else {
			*escSt = escOSC
		}
	default:
		*word = (*word)[:len(*word)-utf8.RuneLen(r)]
		if r == '\x1b' {
			*word = utf8.AppendRune(*word, r)
			*escSt = escStart
			return
		}
		*tempEsc = (*tempEsc)[:0]
		if r == ' ' {
			flushWord()
			return
		}
		if r == '\n' {
			flushWord()
			if line.Len() > 0 {
				*lines = append(*lines, line.String())
				line.Reset()
				line.Write(*activeStyle)
				*lineStartStyle = append((*lineStartStyle)[:0], *activeStyle...)
				*lineVisLen = 0
			} else {
				*lines = append(*lines, "")
			}
			return
		}
		if runewidth.RuneWidth(r) > 1 {
			flushWord()
			rw := runewidth.RuneWidth(r)
			if *lineVisLen+rw > width && *lineVisLen > 0 {
				*lines = append(*lines, line.String())
				line.Reset()
				line.Write(*activeStyle)
				*lineStartStyle = append((*lineStartStyle)[:0], *activeStyle...)
				*lineVisLen = 0
			}
			line.WriteRune(r)
			*lineVisLen += rw
			return
		}
		*word = utf8.AppendRune(*word, r)
		*wordVisLen += runewidth.RuneWidth(r)
	}
}

func tableColWidths(headers []cellData, bodyRows [][]cellData, numCols, maxWidth int) ([]int, []extensionAst.Alignment) {
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

	return colWidths, colAligns
}

func (r *nodeRenderer) renderTable(table *extensionAst.Table) {
	var headers []cellData
	var bodyRows [][]cellData

	for child := table.FirstChild(); child != nil; child = child.NextSibling() {
		switch n := child.(type) {
		case *extensionAst.TableHeader:
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				if cell, ok := c.(*extensionAst.TableCell); ok {
					content := strings.TrimSuffix(r.renderSubtree(cell), "\n")
					headers = append(headers, cellData{
						content: content,
						align:   cell.Alignment,
					})
				}
			}
		case *extensionAst.TableRow:
			var rowCells []cellData
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				if cell, ok := c.(*extensionAst.TableCell); ok {
					content := strings.TrimSuffix(r.renderSubtree(cell), "\n")
					rowCells = append(rowCells, cellData{
						content: content,
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

	colWidths, colAligns := tableColWidths(headers, bodyRows, numCols, r.wordWrap-2*r.padding)
	if colWidths == nil {
		return
	}

	for i, h := range headers {
		headers[i].lines = wrapCellLines(h.content, colWidths[i])
	}
	for _, row := range bodyRows {
		for i, cell := range row {
			row[i].lines = wrapCellLines(cell.content, colWidths[i])
		}
	}


	if len(headers) > 0 {
		r.renderTableRow(headers, colWidths, colAligns)
		r.renderHRule(colWidths)
	}

	for i, row := range bodyRows {
		r.renderTableRow(row, colWidths, colAligns)
		if i < len(bodyRows)-1 {
			r.renderHRule(colWidths)
		}
	}
}

func (r *nodeRenderer) renderHRule(colWidths []int) {
	border := string(r.th.Border.start())
	reset := string(r.th.Border.end())

	r.buf.WriteString(border)
	for i, w := range colWidths {
		for j := 0; j < w+2; j++ {
			r.buf.WriteString("\u2500")
		}
		if i < len(colWidths)-1 {
			r.buf.WriteString("\u253c")
		}
	}
	r.buf.WriteString(reset)
	r.buf.WriteByte('\n')
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
			r.writeTableCell(lineIdx, i, cell, widths, aligns)
		}
		r.buf.WriteByte('\n')
	}
}

func (r *nodeRenderer) writeTableCell(lineIdx, i int, cell cellData, widths []int, aligns []extensionAst.Alignment) {
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

	writeReset := content != ""
	writeOSC8 := writeReset && hasOSC8String(content) && !endsWithOSC8End(content)

	switch aligns[i] {
	case extensionAst.AlignRight:
		for j := 0; j < padding; j++ {
			r.buf.WriteByte(' ')
		}
		r.buf.WriteString(content)
		if writeOSC8 {
			r.buf.WriteString(osc8End)
		}
		if writeReset {
			r.buf.WriteString("\x1b[22m\x1b[23m\x1b[24m\x1b[39m")
		}
	case extensionAst.AlignCenter:
		leftPad := padding / 2
		rightPad := padding - leftPad
		for j := 0; j < leftPad; j++ {
			r.buf.WriteByte(' ')
		}
		r.buf.WriteString(content)
		if writeOSC8 {
			r.buf.WriteString(osc8End)
		}
		if writeReset {
			r.buf.WriteString("\x1b[22m\x1b[23m\x1b[24m\x1b[39m")
		}
		for j := 0; j < rightPad; j++ {
			r.buf.WriteByte(' ')
		}
	default:
		r.buf.WriteString(content)
		if writeOSC8 {
			r.buf.WriteString(osc8End)
		}
		if writeReset {
			r.buf.WriteString("\x1b[22m\x1b[23m\x1b[24m\x1b[39m")
		}
		for j := 0; j < padding; j++ {
			r.buf.WriteByte(' ')
		}
	}
	r.buf.WriteByte(' ')
}
