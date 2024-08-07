package usecases

import (
	"strings"

	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
)

var (
	// daily memo
	HEADING_NAME_TITLE     = markdown.Heading{Level: 1, Text: "daily memo"}
	HEADING_NAME_TODAYSTIP = markdown.Heading{Level: 2, Text: "today's tip"}
	HEADING_NAME_TODOS     = markdown.Heading{Level: 2, Text: "todos"}
	HEADING_NAME_WANTTODOS = markdown.Heading{Level: 2, Text: "wanttodos"}
	HEADING_NAME_MEMOS     = markdown.Heading{Level: 2, Text: "memos"}

	// weekly report
	HEADING_NAME_WEEKLYREPORT = markdown.Heading{Level: 1, Text: "Weekly Report"}

	// tips index
	HEADING_NAME_TIPSINDEX = markdown.Heading{Level: 1, Text: "Tips Index"}
)

var (
	TemplateDailymemo = models.Template{
		Headings: []markdown.Heading{
			HEADING_NAME_TITLE,
			HEADING_NAME_TODAYSTIP,
			HEADING_NAME_TODOS,
			HEADING_NAME_WANTTODOS,
			HEADING_NAME_MEMOS,
		},
	}
	TemplateWeeklyReport = models.Template{
		Headings: []markdown.Heading{
			HEADING_NAME_WEEKLYREPORT,
		},
	}
	TemplateTips = models.Template{
		Headings: []markdown.Heading{
			{Level: 1, Text: "sushi (<- CATEGORY NAME HERE)"},
			{Level: 2, Text: "how to eat sushi (<- YOUR TIPS HERE)"},
			{Level: 2, Text: "how to roll sushi (<- ANOTHER RELATED TIPS HERE)"},
		},
	}
	TemplateTipsIndex = models.Template{
		Headings: []markdown.Heading{
			HEADING_NAME_TIPSINDEX,
		},
	}
)

func GenerateTemplateString(t models.Template) string {
	sb := strings.Builder{}
	for i, h := range t.Headings {
		sb.WriteString(markdown.BuildHeading(h.Level, h.Text) + "\n")
		if i < len(t.Headings)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
