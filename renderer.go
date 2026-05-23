package gruff

import (
	"strings"

	"github.com/yuin/goldmark/ast"
)

type nodeRenderer struct {
	buf    strings.Builder
	source []byte
	th     theme
}

func renderMarkdown(source []byte, th theme, node ast.Node) string {
	var r nodeRenderer
	r.source = source
	r.th = th
	r.renderNode(node)
	return r.buf.String()
}

func (r *nodeRenderer) renderNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.Document:
		r.renderChildren(n)

	case *ast.Paragraph:
		r.renderChildren(n)
		if !r.isInsideList(n) {
			r.buf.WriteString("\n\n")
		}

	case *ast.Heading:
		st := r.headingStyle(n.Level)
		r.buf.WriteString(string(st.start()))
		r.renderChildren(n)
		r.buf.WriteString(string(st.end()))
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
		st := r.th.em
		if n.Level == 2 {
			st = r.th.strong
		}
		r.buf.WriteString(string(st.start()))
		r.renderChildren(n)
		r.buf.WriteString(string(st.end()))

	case *ast.CodeSpan:
		r.buf.WriteString(string(r.th.code.start()))
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if text, ok := c.(*ast.Text); ok {
				r.buf.Write(text.Value(r.source))
			}
		}
		r.buf.WriteString(string(r.th.code.end()))

	case *ast.Link:
		st := r.th.link
		r.buf.WriteString(string(st.start()))
		r.renderChildren(n)
		r.buf.WriteString(string(st.end()))
		if len(n.Destination) > 0 {
			url := string(n.Destination)
			uSt := r.th.linkURL
			r.buf.WriteByte(' ')
			r.buf.WriteString(string(uSt.start()))
			r.buf.WriteByte('(')
			r.buf.WriteString(url)
			r.buf.WriteByte(')')
			r.buf.WriteString(string(uSt.end()))
		}

	case *ast.Image:
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if text, ok := c.(*ast.Text); ok {
				r.buf.Write(text.Value(r.source))
			}
		}

	case *ast.ThematicBreak:
		r.buf.WriteString("\x1b[90m────────────────────\x1b[0m\n\n")

	default:
		r.renderChildren(n)
	}
}

func (r *nodeRenderer) renderChildren(node ast.Node) {
	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		r.renderNode(c)
	}
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
		r.buf.WriteString(string(r.th.numbered.start()))
		r.buf.WriteString(itoa(num))
		r.buf.WriteString(". ")
		r.buf.WriteString(string(r.th.numbered.end()))
	} else {
		r.buf.WriteString("  ")
		r.buf.WriteString(string(r.th.bullet.start()))
		r.buf.WriteString("• ")
		r.buf.WriteString(string(r.th.bullet.end()))
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

func (r *nodeRenderer) headingStyle(level int) style {
	switch level {
	case 1:
		return r.th.h1
	case 2:
		return r.th.h2
	case 3:
		return r.th.h3
	case 4:
		return r.th.h4
	case 5:
		return r.th.h5
	default:
		return r.th.h6
	}
}
