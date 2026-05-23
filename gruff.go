package gruff

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Option func(*options)

type options struct {
	theme    theme
	wordWrap int
}

func WithDark() Option {
	return func(o *options) {
		o.theme = darkTheme
	}
}

func WithLight() Option {
	return func(o *options) {
		o.theme = lightTheme
	}
}

func WithWordWrap(n int) Option {
	return func(o *options) {
		o.wordWrap = n
	}
}

func Render(source string, opts ...Option) (string, error) {
	o := options{
		theme: darkTheme,
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

	out := renderMarkdown(sourceBytes, o.theme, doc)

	if o.wordWrap > 0 {
		out = wrapText(out, o.wordWrap)
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
	wordLen := 0
	inAnsi := false

	flushWord := func() {
		w := word.String()
		word.Reset()
		if w == "" {
			return
		}
		if lineLen > 0 && lineLen+1+wordLen > width {
			out.WriteByte('\n')
			lineLen = 0
		} else if lineLen > 0 {
			out.WriteByte(' ')
			lineLen++
		}
		out.WriteString(w)
		lineLen += wordLen
		wordLen = 0
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
		if r == ' ' || r == '\n' {
			flushWord()
			if r == '\n' {
				out.WriteByte('\n')
				lineLen = 0
			}
			continue
		}
		word.WriteRune(r)
		wordLen++
	}
	flushWord()

	return out.String()
}
