package gruff

import (
	"strings"
	"unicode/utf8"

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
	o.Theme.inheritFg()

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

	out := renderMarkdown(sourceBytes, o.Theme, o.WordWrap, o.Theme.Document.Padding, doc)
	out = strings.TrimSpace(out)

	bgCode := string(ansiBg(o.Theme.Bg))
	if o.WordWrap > 0 {
		out = wrapText(out, o.WordWrap, o.Theme.Document.Padding, bgCode)
	}

	out += string(ansiDefaultBg)
	return out, nil
}

func RenderBytes(source []byte, opts ...Option) ([]byte, error) {
	s, err := Render(string(source), opts...)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

func wrapText(s string, width int, padding int, bgCode string) string {
	if width <= 0 {
		return s
	}

	var out strings.Builder
	out.Grow(len(s) + len(s)/(width+1) + 32)

	// ① 首行前：激活主题背景色，后续填充空格才有背景色
	out.WriteString(bgCode)
	// 首行左内边距（<padding> 个空格）
	for range padding {
		out.WriteByte(' ')
	}
	lineLen := padding

	word := make([]byte, 0, 64)
	spaces := 0
	inAnsi := false

	fillWidth := width

	flushWord := func() {
		if len(word) == 0 {
			return
		}
		wLen := ansiDisplayWidth(word)
		if lineLen > padding && lineLen+wLen+(b2i(spaces > 0)) > width-padding {
			// ② 换行前：填充本行剩余位置（含右内边距）的背景色
			out.WriteString(bgCode)
			for i := lineLen; i < fillWidth; i++ {
				out.WriteByte(' ')
			}
			out.WriteByte('\n')
			// ③ 换行后：下一行开头重新激活背景色
			out.WriteString(bgCode)
			// 下一行左内边距
			for range padding {
				out.WriteByte(' ')
			}
			lineLen = padding
			spaces = 0
		} else if spaces > 0 {
			for i := 0; i < spaces; i++ {
				out.WriteByte(' ')
			}
			lineLen += spaces
		}
		out.Write(word)
		lineLen += wLen
		spaces = 0
		word = word[:0]
	}

	for _, r := range s {
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
		if r == '\n' {
			flushWord()
			// ④ 硬换行前：填充本行剩余位置（含右内边距）的背景色
			out.WriteString(bgCode)
			for i := lineLen; i < fillWidth; i++ {
				out.WriteByte(' ')
			}
			out.WriteByte('\n')
			// ⑤ 硬换行后：下一行开头重新激活背景色
			out.WriteString(bgCode)
			// 下一行左内边距
			for range padding {
				out.WriteByte(' ')
			}
			lineLen = padding
			spaces = 0
			continue
		}
		if r == ' ' {
			flushWord()
			spaces++
			continue
		}
		word = utf8.AppendRune(word, r)
	}
	flushWord()

	// ⑥ 文档末尾：填充最后一行剩余位置（含右内边距）的背景色
	out.WriteString(bgCode)
	for i := lineLen; i < fillWidth; i++ {
		out.WriteByte(' ')
	}

	return out.String()
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
