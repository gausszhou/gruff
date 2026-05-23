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

type color8bit uint8

const (
	cBlack   color8bit = 0
	cMaroon  color8bit = 1
	cGreen   color8bit = 2
	cOlive   color8bit = 3
	cNavy    color8bit = 4
	cPurple  color8bit = 5
	cTeal    color8bit = 6
	cSilver  color8bit = 7
	cGrey    color8bit = 8
	cRed     color8bit = 9
	cLime    color8bit = 10
	cYellow  color8bit = 11
	cBlue    color8bit = 12
	cFuchsia color8bit = 13
	cCyan    color8bit = 14
	cWhite   color8bit = 15
	cDarkBG  color8bit = 236
)

func ansiFg(c color8bit) ansiCode {
	if c <= 7 {
		return ansiCode("\x1b[3") + ansiCode('0'+c) + ansiCode("m")
	}
	return ansiCode("\x1b[38;5;") + ansiCode(itoa(int(c))) + ansiCode("m")
}

func ansiBg(c color8bit) ansiCode {
	if c <= 7 {
		return ansiCode("\x1b[4") + ansiCode('0'+c) + ansiCode("m")
	}
	return ansiCode("\x1b[48;5;") + ansiCode(itoa(int(c))) + ansiCode("m")
}

type style struct {
	fg        color8bit
	bg        color8bit
	bold      bool
	italic    bool
	underline bool
	prefix    string
	suffix    string
}

func (s style) start() ansiCode {
	var out string
	if s.bold {
		out += string(ansiBold)
	}
	if s.italic {
		out += string(ansiItalic)
	}
	if s.underline {
		out += string(ansiUnderline)
	}
	if s.fg != 0 || s.bg != 0 {
		if s.fg != 0 {
			out += string(ansiFg(s.fg))
		}
		if s.bg != 0 {
			out += string(ansiBg(s.bg))
		}
	}
	return ansiCode(out)
}

var ansiResetStr = string(ansiReset)

func (s style) end() ansiCode {
	var out string
	if s.italic {
		out += string(ansiNoItalic)
	}
	if s.bold {
		out += string(ansiNoBold)
	}
	if s.underline {
		out += "\x1b[24m"
	}
	if s.fg != 0 {
		out += "\x1b[39m"
	}
	if s.bg != 0 {
		out += "\x1b[49m"
	}
	if out == "" {
		out = string(ansiReset)
	}
	return ansiCode(out)
}

type theme struct {
	h1, h2, h3, h4, h5, h6 style
	strong                  style
	em                      style
	code                    style
	link                    style
	linkURL                 style
	bullet                  style
	numbered                style
}

var darkTheme = theme{
	h1:       style{bold: true, underline: true, fg: cWhite},
	h2:       style{bold: true, fg: cYellow},
	h3:       style{bold: true, fg: cGreen},
	h4:       style{bold: true, fg: cCyan},
	h5:       style{bold: true, fg: cGrey},
	h6:       style{fg: cGrey},
	strong:   style{bold: true},
	em:       style{italic: true},
	code:     style{bg: cDarkBG, fg: cWhite},
	link:     style{underline: true, fg: cCyan},
	linkURL:  style{fg: cGrey},
	bullet:   style{fg: cYellow},
	numbered: style{fg: cYellow},
}

var lightTheme = theme{
	h1:       style{bold: true, underline: true, fg: cBlack},
	h2:       style{bold: true, fg: cNavy},
	h3:       style{bold: true, fg: cGreen},
	h4:       style{bold: true, fg: cTeal},
	h5:       style{bold: true, fg: cGrey},
	h6:       style{fg: cGrey},
	strong:   style{bold: true},
	em:       style{italic: true},
	code:     style{bg: 7, fg: cBlack},
	link:     style{underline: true, fg: cNavy},
	linkURL:  style{fg: cGrey},
	bullet:   style{fg: cMaroon},
	numbered: style{fg: cMaroon},
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
