package main

import (
	"strings"

	md "github.com/hirotoni/memo/markdown"
)

type Heading struct {
	Level int
	Text  string
}

type Template struct {
	Headings []Heading
}

func (t *Template) String() string {
	sb := strings.Builder{}
	for i, h := range t.Headings {
		sb.WriteString(md.BuildHeading(h.Level, h.Text) + "\n")
		if i < len(t.Headings)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

var (
	HEADING_NAME_TITLE     = Heading{Level: 1, Text: "daily memo"}
	HEADING_NAME_TODOS     = Heading{Level: 2, Text: "todos"}
	HEADING_NAME_WANTTODOS = Heading{Level: 2, Text: "wanttodos"}
	HEADING_NAME_MEMOS     = Heading{Level: 2, Text: "memos"}

	HEADING_NAME_TIPSINDEX = Heading{Level: 1, Text: "Tips Index"}
)

var (
	TemplateDailymemo = Template{
		Headings: []Heading{
			HEADING_NAME_TITLE,
			HEADING_NAME_TODOS,
			HEADING_NAME_WANTTODOS,
			HEADING_NAME_MEMOS,
		},
	}
	TemplateTips = Template{
		Headings: []Heading{
			{Level: 1, Text: "sushi (<- CATEGORY NAME HERE)"},
			{Level: 2, Text: "how to eat sushi (<- YOUR TIPS HERE)"},
			{Level: 2, Text: "how to roll sushi (<- ANOTHER RELATED TIPS HERE)"},
		},
	}
	TemplateTipsIndex = Template{
		Headings: []Heading{
			HEADING_NAME_TIPSINDEX,
		},
	}
)
