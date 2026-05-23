package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gausszhou/gruff"
)

var sampleMD = "" +
	"# H1 — The quick brown fox jumps over the lazy dog\n\n" +
	"## H2 — The quick brown fox jumps over the lazy dog\n\n" +
	"### H3 — The quick brown fox jumps over the lazy dog\n\n" +
	"#### H4 — The quick brown fox jumps over the lazy dog\n\n" +
	"##### H5 — The quick brown fox jumps over the lazy dog\n\n" +
	"###### H6 — The quick brown fox jumps over the lazy dog\n\n" +

	"## Inline Formatting in Headings\n\n" +
	"# **Bold heading** with *italic* and `code`\n\n" +
	"## ***Bold italic*** and ~~strikethrough~~\n\n" +
	"### Link: [gruff](https://github.com/gausszhou/gruff)\n\n" +

	"## CJK & Unicode in Headings\n\n" +
	"# 你好世界 画龙点睛\n\n" +
	"## ✅ 🚀 ⭐ 🌟 Emoji in headings\n\n" +
	"### ∑ π ∫ ∞ λ Math symbols\n\n" +
	"#### éèê àáâ Accented characters\n\n" +
	"##### ★ ♥ ♦ ♣ Dingbats & symbols\n\n" +
	"###### 안녕하세요 こんにちは CJK\n"

func main() {
	light := flag.Bool("light", false, "use light theme")
	wrap := flag.Int("wrap", 60, "word wrap width")
	flag.Parse()

	var opts []gruff.Option
	if *light {
		opts = append(opts, gruff.WithLight())
	}
	if *wrap > 0 {
		opts = append(opts, gruff.WithWordWrap(*wrap))
	}

	out, err := gruff.Render(sampleMD, opts...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(out)
}
