package main

import "strings"

type Heading struct {
	level int
	text  string
}

type Headings []Heading

type Template struct {
	Headings
}

func (t *Template) String() string {
	sb := strings.Builder{}
	for i, h := range t.Headings {
		sb.WriteString(buildHeading(h.level, h.text) + "\n")
		if i < len(t.Headings)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

const (
	HEADING_NAME_TITLE     = "daily memo"
	HEADING_NAME_TODOS     = "todos"
	HEADING_NAME_WANTTODOS = "wanttodos"
	HEADING_NAME_MEMOS     = "memos"
	HEADING_NAME_TIPSINDEX = "Tips Index"
)

var (
	TemplateDailymemo = Template{
		Headings: []Heading{
			{level: 1, text: HEADING_NAME_TITLE},
			{level: 2, text: HEADING_NAME_TODOS},
			{level: 2, text: HEADING_NAME_WANTTODOS},
			{level: 2, text: HEADING_NAME_MEMOS},
		},
	}
	TemplateTips = Template{
		Headings: []Heading{
			{level: 1, text: "sushi (<- CATEGORY NAME HERE)"},
			{level: 2, text: "how to eat sushi (<- YOUR TIPS HERE)"},
			{level: 2, text: "how to roll sushi (<- ANOTHER RELATED TIPS HERE)"},
		},
	}
	TemplateTipsIndex = Template{
		Headings: []Heading{
			{level: 1, text: HEADING_NAME_TIPSINDEX},
		},
	}
)
