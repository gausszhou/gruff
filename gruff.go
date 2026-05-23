package gruff

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

const defaultWordWrap = 80

type Options struct {
	Theme    Theme
	WordWrap int
}

type Option func(*Options)

func WithDark() Option {
	return func(o *Options) {
		o.Theme = darkTheme
	}
}

func WithLight() Option {
	return func(o *Options) {
		o.Theme = lightTheme
	}
}

func WithWordWrap(n int) Option {
	return func(o *Options) {
		o.WordWrap = n
	}
}

func Render(source string, opts ...Option) (string, error) {
	o := Options{
		Theme:    darkTheme,
		WordWrap: defaultWordWrap,
	}
	for _, opt := range opts {
		opt(&o)
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	sourceBytes := []byte(source)
	reader := text.NewReader(sourceBytes)
	doc := md.Parser().Parse(reader)

	out := renderMarkdown(sourceBytes, o.Theme, o.WordWrap, doc)

	if o.WordWrap > 0 {
		out = wrapText(out, o.WordWrap)
	}

	return out, nil
}

func RenderBytes(source []byte, opts ...Option) ([]byte, error) {
	s, err := Render(string(source), opts...)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

func wrapText(s string, width int) string {
	if width <= 0 {
		return s
	}

	var out strings.Builder
	var word strings.Builder
	lineLen := 0
	inAnsi := false

	flushWord := func() {
		w := word.String()
		word.Reset()
		if w == "" {
			return
		}
		wLen := displayWidth(stripANSI(w))
		if lineLen > 0 && lineLen+wLen > width {
			out.WriteByte('\n')
			lineLen = 0
		}
		out.WriteString(w)
		lineLen += wLen
	}

	for _, r := range s {
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
		if r == '\n' {
			flushWord()
			out.WriteByte('\n')
			lineLen = 0
			continue
		}
		if r == ' ' {
			flushWord()
			out.WriteByte(' ')
			lineLen++
			continue
		}
		word.WriteRune(r)
	}
	flushWord()

	return out.String()
}
