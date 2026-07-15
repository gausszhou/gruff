package gruff

import (
	"bytes"
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
			extension.GFM,
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

	out.WriteString(bgCode)
	for range padding {
		out.WriteByte(' ')
	}
	lineLen := padding

	activeStyle := make([]byte, 0, 128)
	wordStartStyle := make([]byte, 0, 128)
	tempEsc := make([]byte, 0, 32)
	word := make([]byte, 0, 64)
	spaces := 0
	escSt := escNone

	fillWidth := width

	newLine := func() {
		if len(activeStyle) > 0 && bytes.HasPrefix(activeStyle, []byte("\x1b]8;")) {
			out.WriteString(osc8End)
		}
		out.WriteString(bgCode)
		for i := lineLen; i < fillWidth; i++ {
			out.WriteByte(' ')
		}
		out.WriteByte('\n')
		out.WriteString(bgCode)
		for range padding {
			out.WriteByte(' ')
		}
		out.Write(activeStyle)
		lineLen = padding
		spaces = 0
	}

	flushWord := func() {
		if len(word) == 0 {
			return
		}
		wLen := ansiDisplayWidth(word)
		if lineLen > padding && lineLen+wLen+(b2i(spaces > 0)) > width-padding {
			savedStyle := activeStyle
			activeStyle = wordStartStyle
			newLine()
			activeStyle = savedStyle
		} else if spaces > 0 {
			for i := 0; i < spaces; i++ {
				out.WriteByte(' ')
			}
			lineLen += spaces
			spaces = 0
		}

		if lineLen == padding && wLen > width-2*padding {
			remaining := word
			splitStyle := make([]byte, len(wordStartStyle))
			copy(splitStyle, wordStartStyle)
			lineLen = breakLongWord(&out, remaining, width, padding, lineLen, &splitStyle, &activeStyle, newLine)
			word = word[:0]
			activeStyle = append(activeStyle[:0], splitStyle...)
			wordStartStyle = append(wordStartStyle[:0], splitStyle...)
			return
		}

		out.Write(word)
		lineLen += wLen
		spaces = 0
		word = word[:0]
		wordStartStyle = append(wordStartStyle[:0], activeStyle...)
	}

	for _, r := range s {
		processWrapRune(r, &escSt, &word, &tempEsc, &activeStyle, flushWord, newLine, &spaces)
	}
	flushWord()

	if len(activeStyle) > 0 {
		out.WriteString("\x1b[0m")
		out.WriteString(osc8End)
	}
	out.WriteString(bgCode)
	for i := lineLen; i < fillWidth; i++ {
		out.WriteByte(' ')
	}

	return out.String()
}

func breakLongWord(out *strings.Builder, word []byte, width, padding, lineLen int, splitStyle, activeStyle *[]byte, newLine func()) int {
	remaining := word
	for len(remaining) > 0 {
		available := width - padding - lineLen
		head, tail := breakChunk(remaining, available, splitStyle)
		if len(head) == 0 {
			out.Write(remaining)
			lineLen += ansiDisplayWidth(remaining)
			break
		}
		out.Write(head)
		lineLen += ansiDisplayWidth(head)
		remaining = tail
		if len(remaining) > 0 {
			savedStyle := *activeStyle
			*activeStyle = *splitStyle
			newLine()
			*activeStyle = savedStyle
			lineLen = padding
		}
	}
	return lineLen
}

func processWrapRune(r rune, escSt *escapeState, word, tempEsc, activeStyle *[]byte, flushWord, newLine func(), spaces *int) {
	switch *escSt {
	case escStart:
		*word = utf8.AppendRune(*word, r)
		*tempEsc = utf8.AppendRune(*tempEsc, r)
		if r == '[' {
			*escSt = escCSI
		} else if r == ']' {
			*escSt = escOSC
		} else {
			*escSt = escNone
			*tempEsc = (*tempEsc)[:0]
		}
	case escCSI:
		*word = utf8.AppendRune(*word, r)
		*tempEsc = utf8.AppendRune(*tempEsc, r)
		if r >= 0x40 && r <= 0x7E {
			*escSt = escNone
			updateActiveStyle(activeStyle, *tempEsc)
			*tempEsc = (*tempEsc)[:0]
		}
	case escOSC:
		*word = utf8.AppendRune(*word, r)
		*tempEsc = utf8.AppendRune(*tempEsc, r)
		if r == '\x1b' {
			*escSt = escOSCSt
		} else if r == '\x07' {
			*escSt = escNone
			updateActiveStyle(activeStyle, *tempEsc)
			*tempEsc = (*tempEsc)[:0]
		}
	case escOSCSt:
		*word = utf8.AppendRune(*word, r)
		*tempEsc = utf8.AppendRune(*tempEsc, r)
		if r == '\\' {
			*escSt = escNone
			updateActiveStyle(activeStyle, *tempEsc)
			*tempEsc = (*tempEsc)[:0]
		} else {
			*escSt = escOSC
		}
	default:
		if r == '\x1b' {
			*word = utf8.AppendRune(*word, r)
			*tempEsc = utf8.AppendRune(*tempEsc, r)
			*escSt = escStart
			return
		}
		if r == '\n' {
			flushWord()
			newLine()
			return
		}
		if r == ' ' {
			flushWord()
			*spaces++
			return
		}
		*word = utf8.AppendRune(*word, r)
	}
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
