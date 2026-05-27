package gruff

import (
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

type ansiCode string

const (
	ansiReset     ansiCode = "\x1b[0m"
	ansiBold      ansiCode = "\x1b[1m"
	ansiItalic    ansiCode = "\x1b[3m"
	ansiUnderline ansiCode = "\x1b[4m"
	ansiNoBold    ansiCode = "\x1b[22m"
	ansiNoItalic  ansiCode = "\x1b[23m"
)

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
	Bg        string
	Bold      bool
	Italic    bool
	Underline bool
	Padding   int
}

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
	if s.Fg != "" || s.Bg != "" {
		if s.Fg != "" {
			out += string(ansiFg(s.Fg))
		}
		if s.Bg != "" {
			out += string(ansiBg(s.Bg))
		}
	}
	return ansiCode(out)
}

var ansiResetStr = string(ansiReset)

func (s Style) end(bg string) ansiCode {
	var out string
	if s.Italic {
		out += string(ansiNoItalic)
	}
	if s.Bold {
		out += string(ansiNoBold)
	}
	if s.Underline {
		out += "\x1b[24m"
	}
	if s.Fg != "" {
		out += "\x1b[39m"
	}
	if s.Bg != "" {
		if bg != "" {
			out += string(ansiBg(bg))
		} else {
			out += "\x1b[49m"
		}
	}
	if out == "" {
		out = "\x1b[39m\x1b[49m"
	}
	return ansiCode(out)
}

type Theme struct {
	Document                Style
	H1, H2, H3, H4, H5, H6 Style
	Strong                  Style
	Em                      Style
	Code                    Style
	Link                    Style
	LinkURL                 Style
	Bullet                  Style
	Numbered                Style
	Hr                      Style
	Border                  Style
}

var darkTheme = Theme{
	Document:    Style{Padding: 2},
	H1:          Style{Bold: true, Fg: "#FFFF87"},
	H2:         Style{Bold: true, Fg: "#00AFFF"},
	H3:         Style{Bold: true, Fg: "#00AFFF"},
	H4:         Style{Bold: true, Fg: "#00AFFF"},
	H5:         Style{Bold: true, Fg: "#00AFFF"},
	H6:         Style{Fg: "#00AF5F"},
	Strong:     Style{Bold: true},
	Em:         Style{Italic: true},
	Code:       Style{Fg: "#A6E22E"},
	Link:       Style{Underline: true, Fg: "#5c9cf5"},
	LinkURL:    Style{Fg: "#808080"},
	Bullet:     Style{Fg: "#FFFF00"},
	Numbered:   Style{Fg: "#FFFF00"},
	Hr:         Style{Fg: "#808080"},
	Border:     Style{Fg: "#808080"},
}

var lightTheme = Theme{
	Document:    Style{Padding: 2},
	H1:          Style{Bold: true, Underline: true, Fg: "#000000"},
	H2:         Style{Bold: true, Fg: "#000080"},
	H3:         Style{Bold: true, Fg: "#008000"},
	H4:         Style{Bold: true, Fg: "#008080"},
	H5:         Style{Bold: true, Fg: "#808080"},
	H6:         Style{Fg: "#808080"},
	Strong:     Style{Bold: true},
	Em:         Style{Italic: true},
	Code:       Style{Fg: "#000000", Padding: 1},
	Link:       Style{Underline: true, Fg: "#000080"},
	LinkURL:    Style{Fg: "#808080"},
	Bullet:     Style{Fg: "#800000"},
	Numbered:   Style{Fg: "#800000"},
	Hr:         Style{Fg: "#808080"},
	Border:     Style{Fg: "#808080"},
}

func displayWidth(s string) int {
	return runewidth.StringWidth(s)
}

func stripANSI(s string) string {
	var out []byte
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			for j := i + 2; j < len(s); j++ {
				if s[j] == 'm' {
					i = j
					break
				}
			}
			continue
		}
		out = append(out, s[i])
	}
	return string(out)
}

func ansiDisplayWidth(b []byte) int {
	w := 0
	for i := 0; i < len(b); {
		if b[i] == '\x1b' && i+1 < len(b) && b[i+1] == '[' {
			for j := i + 2; j < len(b); j++ {
				if b[j] == 'm' {
					i = j + 1
					break
				}
			}
			continue
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


