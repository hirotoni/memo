package models

import (
	"strings"

	md "github.com/hirotoni/memo/markdown"
)

type Template struct {
	Headings []md.Heading
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
	// daily memo
	HEADING_NAME_TITLE     = md.Heading{Level: 1, Text: "daily memo"}
	HEADING_NAME_TODOS     = md.Heading{Level: 2, Text: "todos"}
	HEADING_NAME_WANTTODOS = md.Heading{Level: 2, Text: "wanttodos"}
	HEADING_NAME_MEMOS     = md.Heading{Level: 2, Text: "memos"}

	// weekly report
	HEADING_NAME_WEEKLYREPORT = md.Heading{Level: 1, Text: "Weekly Report"}

	// tips index
	HEADING_NAME_TIPSINDEX = md.Heading{Level: 1, Text: "Tips Index"}
)

var (
	TemplateDailymemo = Template{
		Headings: []md.Heading{
			HEADING_NAME_TITLE,
			HEADING_NAME_TODOS,
			HEADING_NAME_WANTTODOS,
			HEADING_NAME_MEMOS,
		},
	}
	TemplateWeeklyReport = Template{
		Headings: []md.Heading{
			HEADING_NAME_WEEKLYREPORT,
		},
	}
	TemplateTips = Template{
		Headings: []md.Heading{
			{Level: 1, Text: "sushi (<- CATEGORY NAME HERE)"},
			{Level: 2, Text: "how to eat sushi (<- YOUR TIPS HERE)"},
			{Level: 2, Text: "how to roll sushi (<- ANOTHER RELATED TIPS HERE)"},
		},
	}
	TemplateTipsIndex = Template{
		Headings: []md.Heading{
			HEADING_NAME_TIPSINDEX,
		},
	}
)
