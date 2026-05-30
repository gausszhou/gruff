package gruff

import (
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
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
	Bg:            "#141414",
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
	Link:          Style{Underline: true, Fg: "#5c9cf5"},
	LinkURL:       Style{Fg: "#5c9cf5"},
	Hr:            Style{Fg: "#808080"},
	Border:        Style{Fg: "#808080"},
	BlockQuote:    Style{Fg: "#808080"},
	TaskChecked:   Style{Fg: "#50fa7b"},
	TaskUnchecked: Style{Fg: "#808080"},
}

var lightTheme = Theme{
	Bg:            "#f0f0f0",
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
	Link:          Style{Underline: true, Fg: "#5c9cf5"},
	LinkURL:       Style{Fg: "#5c9cf5"},
	Hr:            Style{Fg: "#333333"},
	Border:        Style{Fg: "#333333"},
	BlockQuote:    Style{Fg: "#333333"},
	TaskChecked:   Style{Fg: "#008000"},
	TaskUnchecked: Style{Fg: "#333333"},
}

// displayWidth 返回字符串在终端中占用的显示宽度。
// 对 go-runewidth 的补充：
//   - U+FE0F（Variation Selector-16）单独不占宽度
//   - 后跟 U+FE0F 的码位强制为 emoji 宽度 2（弥补 runewidth 不处理 ambiguous + VS16 的不足）
func displayWidth(s string) int {
	w := 0
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		// U+FE0F（Variation Selector-16）单独出现时宽度为 0
		if r == 0xFE0F {
			i += size
			continue
		}
		// 若后跟 U+FE0F，则该码位为 emoji 呈现 → 宽度 2
		if i+size < len(s) {
			next, nextSize := utf8.DecodeRuneInString(s[i+size:])
			if next == 0xFE0F {
				w += 2
				i += size + nextSize
				continue
			}
		}
		w += runewidth.RuneWidth(r)
		i += size
	}
	return w
}

// stripANSI 移除字符串中所有 ANSI 转义序列，返回纯文本
func stripANSI(s string) string {
	var out []byte
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			for j := i + 2; j < len(s); j++ {
				if s[j] >= 0x40 && s[j] <= 0x7E {
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

// ansiDisplayWidth 计算包含 ANSI 码的字节切片的显示宽度（忽略 ANSI 码后的视觉宽度）
func ansiDisplayWidth(b []byte) int {
	w := 0
	for i := 0; i < len(b); {
		if b[i] == '\x1b' && i+1 < len(b) && b[i+1] == '[' {
			for j := i + 2; j < len(b); j++ {
				if b[j] >= 0x40 && b[j] <= 0x7E {
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
