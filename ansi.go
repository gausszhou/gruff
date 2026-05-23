package gruff

import (
	"unicode/utf8"

	"charm.land/lipgloss/v2"
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

type Color string

const (
	cBlack   Color = "0"
	cMaroon  Color = "1"
	cGreen   Color = "2"
	cOlive   Color = "3"
	cNavy    Color = "4"
	cPurple  Color = "5"
	cTeal    Color = "6"
	cSilver  Color = "7"
	cGrey    Color = "8"
	cRed     Color = "9"
	cLime    Color = "10"
	cYellow  Color = "11"
	cBlue    Color = "12"
	cFuchsia Color = "13"
	cCyan    Color = "14"
	cWhite   Color = "15"
)

func is4bit(c Color) bool {
	return len(c) == 1 && c[0] >= '0' && c[0] <= '7'
}

func isHex(c Color) bool {
	return len(c) >= 4 && c[0] == '#'
}

func hexRGB(c Color) (r, g, b uint8) {
	cc := lipgloss.Color(string(c))
	rr, gg, bb, _ := cc.RGBA()
	return uint8(rr >> 8), uint8(gg >> 8), uint8(bb >> 8)
}

func ansiFg(c Color) ansiCode {
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

func ansiBg(c Color) ansiCode {
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
	Fg        Color
	Bg        Color
	Bold      bool
	Italic    bool
	Underline bool
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

func (s Style) end(bg Color) ansiCode {
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
		out = string(ansiReset)
	}
	return ansiCode(out)
}

type Theme struct {
	Background              Color
	H1, H2, H3, H4, H5, H6 Style
	Strong                  Style
	Em                      Style
	Code                    Style
	Link                    Style
	LinkURL                 Style
	Bullet                  Style
	Numbered                Style
}

var darkTheme = Theme{
	Background: "#141414",
	H1:         Style{Bold: true, Fg: cWhite},
	H2:         Style{Bold: true, Fg: cYellow},
	H3:         Style{Bold: true, Fg: cGreen},
	H4:         Style{Bold: true, Fg: cCyan},
	H5:         Style{Bold: true, Fg: cGrey},
	H6:         Style{Fg: cGrey},
	Strong:     Style{Bold: true},
	Em:         Style{Italic: true},
	Code:       Style{Fg: "#50865a"},
	Link:       Style{Underline: true, Fg: "#5c9cf5"},
	LinkURL:    Style{Fg: cGrey},
	Bullet:     Style{Fg: cYellow},
	Numbered:   Style{Fg: cYellow},
}

var lightTheme = Theme{
	Background: "",
	H1:         Style{Bold: true, Underline: true, Fg: cBlack},
	H2:         Style{Bold: true, Fg: cNavy},
	H3:         Style{Bold: true, Fg: cGreen},
	H4:         Style{Bold: true, Fg: cTeal},
	H5:         Style{Bold: true, Fg: cGrey},
	H6:         Style{Fg: cGrey},
	Strong:     Style{Bold: true},
	Em:         Style{Italic: true},
	Code:       Style{Bg: cSilver, Fg: cBlack},
	Link:       Style{Underline: true, Fg: cNavy},
	LinkURL:    Style{Fg: cGrey},
	Bullet:     Style{Fg: cMaroon},
	Numbered:   Style{Fg: cMaroon},
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

func truncateUTF8(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max])
}
