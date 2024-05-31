package markdown

import (
	"fmt"
	"strings"
)

type Heading struct {
	Level int
	Text  string
}

func Text2tag(text string) string {
	var tag = text
	tag = strings.ReplaceAll(tag, " ", "-")
	tag = strings.ReplaceAll(tag, "#", "")
	fullwidthchars := strings.Split("　！＠＃＄％＾＆＊（）＋｜〜＝￥｀「」｛｝；’：”、。・＜＞？【】『』《》〔〕［］‹›«»〘〙〚〛", "")
	for _, c := range fullwidthchars {
		tag = strings.ReplaceAll(tag, c, "")
	}
	return tag
}

func BuildHeading(level int, text string) string {
	return strings.Repeat("#", level) + " " + text
}

func BuildLink(text, destination string) string {
	return "[" + text + "]" + "(" + destination + ")"
}

func BuildList(text string) string {
	return "- " + text
}

func BuildOrderedList(order int, text string) string {
	return fmt.Sprint(order) + ". " + text
}

func BuildCheckbox(text string, checked bool) string {
	if checked {
		return "- [x] " + text
	} else {
		return "- [ ] " + text
	}
}
