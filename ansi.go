package gruff

import "unicode/utf8"

type ansiCode string

const (
	ansiReset     ansiCode = "\x1b[0m"
	ansiBold      ansiCode = "\x1b[1m"
	ansiItalic    ansiCode = "\x1b[3m"
	ansiUnderline ansiCode = "\x1b[4m"
	ansiNoBold    ansiCode = "\x1b[22m"
	ansiNoItalic  ansiCode = "\x1b[23m"
)

type Color uint8

const (
	cBlack   Color = 0
	cMaroon  Color = 1
	cGreen   Color = 2
	cOlive   Color = 3
	cNavy    Color = 4
	cPurple  Color = 5
	cTeal    Color = 6
	cSilver  Color = 7
	cGrey    Color = 8
	cRed     Color = 9
	cLime    Color = 10
	cYellow  Color = 11
	cBlue    Color = 12
	cFuchsia Color = 13
	cCyan    Color = 14
	cWhite   Color = 15
	cDarkBG  Color = 236
)

func ansiFg(c Color) ansiCode {
	if c <= 7 {
		return ansiCode("\x1b[3") + ansiCode('0'+c) + ansiCode("m")
	}
	return ansiCode("\x1b[38;5;") + ansiCode(itoa(int(c))) + ansiCode("m")
}

func ansiBg(c Color) ansiCode {
	if c <= 7 {
		return ansiCode("\x1b[4") + ansiCode('0'+c) + ansiCode("m")
	}
	return ansiCode("\x1b[48;5;") + ansiCode(itoa(int(c))) + ansiCode("m")
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
	if s.Fg != 0 || s.Bg != 0 {
		if s.Fg != 0 {
			out += string(ansiFg(s.Fg))
		}
		if s.Bg != 0 {
			out += string(ansiBg(s.Bg))
		}
	}
	return ansiCode(out)
}

var ansiResetStr = string(ansiReset)

func (s Style) end() ansiCode {
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
	if s.Fg != 0 {
		out += "\x1b[39m"
	}
	if s.Bg != 0 {
		out += "\x1b[49m"
	}
	if out == "" {
		out = string(ansiReset)
	}
	return ansiCode(out)
}

type Theme struct {
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
	H1:       Style{Bold: true, Underline: true, Fg: cWhite},
	H2:       Style{Bold: true, Fg: cYellow},
	H3:       Style{Bold: true, Fg: cGreen},
	H4:       Style{Bold: true, Fg: cCyan},
	H5:       Style{Bold: true, Fg: cGrey},
	H6:       Style{Fg: cGrey},
	Strong:   Style{Bold: true},
	Em:       Style{Italic: true},
	Code:     Style{Bg: cDarkBG, Fg: cWhite},
	Link:     Style{Underline: true, Fg: cCyan},
	LinkURL:  Style{Fg: cGrey},
	Bullet:   Style{Fg: cYellow},
	Numbered: Style{Fg: cYellow},
}

var lightTheme = Theme{
	H1:       Style{Bold: true, Underline: true, Fg: cBlack},
	H2:       Style{Bold: true, Fg: cNavy},
	H3:       Style{Bold: true, Fg: cGreen},
	H4:       Style{Bold: true, Fg: cTeal},
	H5:       Style{Bold: true, Fg: cGrey},
	H6:       Style{Fg: cGrey},
	Strong:   Style{Bold: true},
	Em:       Style{Italic: true},
	Code:     Style{Bg: 7, Fg: cBlack},
	Link:     Style{Underline: true, Fg: cNavy},
	LinkURL:  Style{Fg: cGrey},
	Bullet:   Style{Fg: cMaroon},
	Numbered: Style{Fg: cMaroon},
}

func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		if r >= 0x4E00 && r <= 0x9FFF ||
			r >= 0x3000 && r <= 0x303F ||
			r >= 0xFF00 && r <= 0xFFEF ||
			r >= 0x20000 && r <= 0x2FFFF {
			w += 2
		} else {
			w++
		}
	}
	return w
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
