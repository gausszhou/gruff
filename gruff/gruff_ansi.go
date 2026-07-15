package gruff

import (
	"bytes"
	"unicode/utf8"

	"github.com/clipperhouse/displaywidth"
	"github.com/mattn/go-runewidth"
)

func init() {
	runewidth.DefaultCondition.EastAsianWidth = false
}

type escapeState byte

const (
	escNone escapeState = iota
	escStart
	escCSI
	escOSC
	escOSCSt
)

type ansiCode string

const (
	ansiReset       ansiCode = "\x1b[0m"
	ansiBold        ansiCode = "\x1b[1m"
	ansiItalic      ansiCode = "\x1b[3m"
	ansiUnderline   ansiCode = "\x1b[4m"
	ansiNoBold      ansiCode = "\x1b[22m"
	ansiNoItalic    ansiCode = "\x1b[23m"
	ansiNoUnderline ansiCode = "\x1b[24m"
	ansiDefaultFg   ansiCode = "\x1b[39m"
	ansiDefaultBg   ansiCode = "\x1b[49m"
)

func osc8Link(url string) string {
	return "\x1b]8;id=" + url + ";" + url + "\x1b\\"
}

const osc8End = "\x1b]8;;\x1b\\"

func is4bit(c string) bool {
	return len(c) == 1 && c[0] >= '0' && c[0] <= '7'
}

func isHex(c string) bool {
	return len(c) >= 4 && c[0] == '#'
}

func hexRGB(c string) (r, g, b uint8) {
	if len(c) < 7 || c[0] != '#' {
		return 0, 0, 0
	}
	hex := func(b1, b2 byte) uint8 {
		hn := func(b byte) uint8 {
			switch {
			case '0' <= b && b <= '9':
				return b - '0'
			case 'a' <= b && b <= 'f':
				return 10 + b - 'a'
			case 'A' <= b && b <= 'F':
				return 10 + b - 'A'
			default:
				return 0
			}
		}
		return hn(b1)<<4 | hn(b2)
	}
	return hex(c[1], c[2]), hex(c[3], c[4]), hex(c[5], c[6])
}

func ansiFg(c string) ansiCode {
	if c == "" {
		return ""
	}
	if isHex(c) {
		r, g, b := hexRGB(c)
		return ansiCode("\x1b[38;2;" + itoa(int(r)) + ";" + itoa(int(g)) + ";" + itoa(int(b)) + "m")
	}
	if is4bit(c) {
		return ansiCode("\x1b[3") + ansiCode(string(c)) + ansiCode("m")
	}
	return ansiCode("\x1b[38;5;" + string(c) + "m")
}

func ansiBg(c string) ansiCode {
	if c == "" {
		return ""
	}
	if isHex(c) {
		r, g, b := hexRGB(c)
		return ansiCode("\x1b[48;2;" + itoa(int(r)) + ";" + itoa(int(g)) + ";" + itoa(int(b)) + "m")
	}
	if is4bit(c) {
		return ansiCode("\x1b[4") + ansiCode(string(c)) + ansiCode("m")
	}
	return ansiCode("\x1b[48;5;" + string(c) + "m")
}

type Style struct {
	Fg        string
	Bold      bool
	Italic    bool
	Underline bool
	Padding   int
}

// start 输出样式起始 ANSI 码：加粗/斜体/下划线/前景色
func (s Style) start() ansiCode {
	var out string
	if s.Bold {
		out += string(ansiBold)
	}
	if s.Italic {
		out += string(ansiItalic)
	}
	if s.Underline {
		out += string(ansiUnderline)
	}
	if s.Fg != "" {
		out += string(ansiFg(s.Fg))
	}
	return ansiCode(out)
}

// end 关闭样式，前景重置
func (s Style) end() ansiCode {
	var out string
	if s.Italic {
		out += string(ansiNoItalic)
	}
	if s.Bold {
		out += string(ansiNoBold)
	}
	if s.Underline {
		out += string(ansiNoUnderline)
	}
	if s.Fg != "" {
		out += string(ansiDefaultFg)
	}
	return ansiCode(out)
}

type Theme struct {
	Bg                     string
	Document               Style
	Paragraph              Style
	H1, H2, H3, H4, H5, H6 Style
	Strong                 Style
	Em                     Style
	Code                   Style
	Link                   Style
	LinkURL                Style
	Hr                     Style
	Border                 Style
	BlockQuote             Style
	TaskChecked            Style
	TaskUnchecked          Style
}

func (th *Theme) inheritFg() {
	fg := th.Document.Fg
	inherit := func(s *Style) {
		if s.Fg == "" {
			s.Fg = fg
		}
	}
	inherit(&th.Paragraph)
	inherit(&th.H1)
	inherit(&th.H2)
	inherit(&th.H3)
	inherit(&th.H4)
	inherit(&th.H5)
	inherit(&th.H6)
	inherit(&th.Strong)
	inherit(&th.Em)
	inherit(&th.Code)
	inherit(&th.Link)
	inherit(&th.LinkURL)
	inherit(&th.Hr)
	inherit(&th.Border)
	inherit(&th.BlockQuote)
	inherit(&th.TaskChecked)
	inherit(&th.TaskUnchecked)
}

var darkTheme = Theme{
	Bg:            "",
	Document:      Style{Padding: 1, Fg: "#e0e0e0"},
	Paragraph:     Style{Fg: "#e0e0e0"},
	H1:            Style{Bold: true, Fg: "#FFFF87"},
	H2:            Style{Bold: true, Fg: "#00AFFF"},
	H3:            Style{Bold: true, Fg: "#00AFFF"},
	H4:            Style{Bold: true, Fg: "#00AFFF"},
	H5:            Style{Bold: true, Fg: "#00AFFF"},
	H6:            Style{Fg: "#00AFFF"},
	Strong:        Style{Bold: true},
	Em:            Style{Italic: true},
	Code:          Style{Fg: "#50fa7b"},
	Link:          Style{Bold: true, Fg: "#5c9cf5"},
	LinkURL:       Style{Underline: true, Fg: "#5c9cf5"},
	Hr:            Style{Fg: "#808080"},
	Border:        Style{Fg: "#808080"},
	BlockQuote:    Style{Fg: "#808080"},
	TaskChecked:   Style{Fg: "#50fa7b"},
	TaskUnchecked: Style{Fg: "#808080"},
}

var lightTheme = Theme{
	Bg:            "",
	Document:      Style{Padding: 1, Fg: "#333333"},
	Paragraph:     Style{Fg: "#333333"},
	H1:            Style{Bold: true, Underline: true, Fg: "#000000"},
	H2:            Style{Bold: true, Fg: "#00AFFF"},
	H3:            Style{Bold: true, Fg: "#00AFFF"},
	H4:            Style{Bold: true, Fg: "#00AFFF"},
	H5:            Style{Bold: true, Fg: "#00AFFF"},
	H6:            Style{Fg: "#00AFFF"},
	Strong:        Style{Bold: true},
	Em:            Style{Italic: true},
	Code:          Style{Fg: "#008000", Padding: 1},
	Link:          Style{Bold: true, Fg: "#5c9cf5"},
	LinkURL:       Style{Underline: true, Fg: "#5c9cf5"},
	Hr:            Style{Fg: "#333333"},
	Border:        Style{Fg: "#333333"},
	BlockQuote:    Style{Fg: "#333333"},
	TaskChecked:   Style{Fg: "#008000"},
	TaskUnchecked: Style{Fg: "#333333"},
}

func displayWidth(s string) int {
	return displaywidth.String(s)
}

// stripANSI 移除字符串中所有 ANSI 转义序列，返回纯文本
func stripANSI(s string) string {
	var out []byte
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' && i+1 < len(s) {
			switch s[i+1] {
			case '[':
				for j := i + 2; j < len(s); j++ {
					if s[j] >= 0x40 && s[j] <= 0x7E {
						i = j
						break
					}
				}
				continue
			case ']':
				found := false
				for j := i + 2; j < len(s); j++ {
					if s[j] == '\x07' {
						i = j
						found = true
						break
					}
					if s[j] == '\x1b' && j+1 < len(s) && s[j+1] == '\\' {
						i = j + 1
						found = true
						break
					}
				}
				if !found {
					i = len(s) - 1
				}
				continue
			}
		}
		out = append(out, s[i])
	}
	return string(out)
}

func ansiDisplayWidth(b []byte) int {
	w := 0
	for i := 0; i < len(b); {
		if b[i] == '\x1b' && i+1 < len(b) {
			switch b[i+1] {
			case '[':
				found := false
				for j := i + 2; j < len(b); j++ {
					if b[j] >= 0x40 && b[j] <= 0x7E {
						i = j + 1
						found = true
						break
					}
				}
				if !found {
					i = len(b)
				}
				continue
			case ']':
				found := false
				for j := i + 2; j < len(b); j++ {
					if b[j] == '\x07' {
						i = j + 1
						found = true
						break
					}
					if b[j] == '\x1b' && j+1 < len(b) && b[j+1] == '\\' {
						i = j + 2
						found = true
						break
					}
				}
				if !found {
					i = len(b)
				}
				continue
			}
		}
		r, size := utf8.DecodeRune(b[i:])
		w += runewidth.RuneWidth(r)
		i += size
	}
	return w
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [3]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

func breakChunk(b []byte, maxWidth int, style *[]byte) (head, tail []byte) {
	w := 0
	for i := 0; i < len(b); {
		if b[i] == '\x1b' && i+1 < len(b) {
			escStart := i
			switch b[i+1] {
			case '[':
				found := false
				for j := i + 2; j < len(b); j++ {
					if b[j] >= 0x40 && b[j] <= 0x7E {
						i = j + 1
						found = true
						break
					}
				}
				if !found {
					return b, nil
				}
				updateActiveStyle(style, b[escStart:i])
				continue
			case ']':
				found := false
				for j := i + 2; j < len(b); j++ {
					if b[j] == '\x07' {
						i = j + 1
						found = true
						break
					}
					if b[j] == '\x1b' && j+1 < len(b) && b[j+1] == '\\' {
						i = j + 2
						found = true
						break
					}
				}
				if !found {
					return b, nil
				}
				updateActiveStyle(style, b[escStart:i])
				continue
			}
		}
		r, size := utf8.DecodeRune(b[i:])
		rw := runewidth.RuneWidth(r)
		if w+rw > maxWidth {
			return b[:i], b[i:]
		}
		w += rw
		i += size
	}
	return b, nil
}

func removeExact(active []byte, seq string) []byte {
	b := []byte(seq)
	for i := 0; i <= len(active)-len(b); i++ {
		if bytes.Equal(active[i:i+len(b)], b) {
			return append(active[:i], active[i+len(b):]...)
		}
	}
	return active
}

func removeCSIPrefix(active []byte, prefix string) []byte {
	pb := []byte(prefix)
	for i := 0; i < len(active); i++ {
		if active[i] == '\x1b' && i+1 < len(active) && active[i+1] == '[' {
			paramStart := i + 2
			for j := paramStart; j < len(active); j++ {
				if active[j] >= 0x40 && active[j] <= 0x7E {
					if bytes.HasPrefix(active[paramStart:j], pb) {
						return append(active[:i], active[j+1:]...)
					}
					i = j
					break
				}
			}
		}
	}
	return active
}

func removeOSC(active []byte) []byte {
	for i := 0; i < len(active); i++ {
		if active[i] == '\x1b' && i+1 < len(active) && active[i+1] == ']' {
			for j := i + 2; j < len(active); j++ {
				if active[j] == '\x07' {
					return append(active[:i], active[j+1:]...)
				}
				if active[j] == '\x1b' && j+1 < len(active) && active[j+1] == '\\' {
					return append(active[:i], active[j+2:]...)
				}
			}
		}
	}
	return active
}

func updateActiveStyle(active *[]byte, esc []byte) {
	if bytes.HasPrefix(esc, []byte("\x1b]8;")) {
		*active = removeOSC(*active)
		if !bytes.Equal(esc, []byte(osc8End)) {
			*active = append(*active, esc...)
		}
		return
	}
	switch {
	case bytes.Equal(esc, []byte("\x1b[22m")):
		*active = removeExact(*active, "\x1b[1m")
	case bytes.Equal(esc, []byte("\x1b[23m")):
		*active = removeExact(*active, "\x1b[3m")
	case bytes.Equal(esc, []byte("\x1b[24m")):
		*active = removeExact(*active, "\x1b[4m")
	case bytes.Equal(esc, []byte("\x1b[39m")):
		*active = removeCSIPrefix(*active, "38;")
	case bytes.Equal(esc, []byte("\x1b[49m")):
		*active = removeCSIPrefix(*active, "48;")
	case bytes.Equal(esc, []byte("\x1b[0m")):
		*active = (*active)[:0]
	default:
		if bytes.HasPrefix(esc, []byte("\x1b[38;")) {
			*active = removeCSIPrefix(*active, "38;")
			*active = append(*active, esc...)
		} else if bytes.HasPrefix(esc, []byte("\x1b[48;")) {
			*active = removeCSIPrefix(*active, "48;")
			*active = append(*active, esc...)
		} else {
			*active = append(*active, esc...)
		}
	}
}
