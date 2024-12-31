package usecases

import (
	"strings"

	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
)

var (
	// daily memo
	HEADING_NAME_TITLE     = models.NewHeading(1, "daily memo")
	HEADING_NAME_TODAYSTIP = models.NewHeading(2, "today's tip")
	HEADING_NAME_TODOS     = models.NewHeading(2, "todos")
	HEADING_NAME_WANTTODOS = models.NewHeading(2, "wanttodos")
	HEADING_NAME_MEMOS     = models.NewHeading(2, "memos")
	// weekly report
	HEADING_NAME_WEEKLYREPORT = models.NewHeading(1, "Weekly Report")
	// tips index
	HEADING_NAME_TIPSINDEX = models.NewHeading(1, "Tips Index")
)

var (
	dailymemoHeadings = []models.Heading{
		HEADING_NAME_TITLE,
		HEADING_NAME_TODAYSTIP,
		HEADING_NAME_TODOS,
		HEADING_NAME_WANTTODOS,
		HEADING_NAME_MEMOS,
	}
	weeklyReportHeadings = []models.Heading{
		HEADING_NAME_WEEKLYREPORT,
	}
	tipsHeadings = []models.Heading{
		models.NewHeading(1, "sushi (<- tip category)"),
		models.NewHeading(2, "how to eat sushi (<- tip title in heading level 2)"),
		models.NewHeading(2, "how to roll sushi (<- another tip)"),
	}
	tipsIndexHeadings = []models.Heading{
		HEADING_NAME_TIPSINDEX,
	}
)

var (
	TemplateDailymemo    = models.NewTemplate(dailymemoHeadings)
	TemplateWeeklyReport = models.NewTemplate(weeklyReportHeadings)
	TemplateTips         = models.NewTemplate(tipsHeadings)
	TemplateTipsIndex    = models.NewTemplate(tipsIndexHeadings)
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
