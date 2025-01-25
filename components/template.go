package components

import (
	"strings"

	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
)

var (
	// daily memo
	HEADING_NAME_TITLE             = markdown.NewHeading(1, "daily memo")
	HEADING_NAME_TODAYSMEMOARCHIVE = markdown.NewHeading(2, "today's memo archive")
	HEADING_NAME_TODOS             = markdown.NewHeading(2, "todos")
	HEADING_NAME_WANTTODOS         = markdown.NewHeading(2, "wanttodos")
	HEADING_NAME_MEMOS             = markdown.NewHeading(2, "memos")
	// weekly report
	HEADING_NAME_WEEKLYREPORT = markdown.NewHeading(1, "Weekly Report")
	// memo archives index
	HEADING_NAME_MEMOARCHIVES_INDEX = markdown.NewHeading(1, "Memo Archives Index")
)

var (
	dailymemoHeadings = []markdown.Heading{
		HEADING_NAME_TITLE,
		HEADING_NAME_TODAYSMEMOARCHIVE,
		HEADING_NAME_TODOS,
		HEADING_NAME_WANTTODOS,
		HEADING_NAME_MEMOS,
	}
	weeklyReportHeadings = []markdown.Heading{
		HEADING_NAME_WEEKLYREPORT,
	}
	memoArchivesHeadings = []markdown.Heading{
		markdown.NewHeading(1, "sushi (<- memo category)"),
		markdown.NewHeading(2, "how to eat sushi (<- memo title in heading level 2)"),
		markdown.NewHeading(2, "how to roll sushi (<- another memo)"),
	}
	memoArchivesIndexHeadings = []markdown.Heading{
		HEADING_NAME_MEMOARCHIVES_INDEX,
	}
)

var (
	TemplateDailymemo         = models.NewTemplate(dailymemoHeadings)
	TemplateWeeklyReport      = models.NewTemplate(weeklyReportHeadings)
	TemplateMemoArchives      = models.NewTemplate(memoArchivesHeadings)
	TemplateMemoArchivesIndex = models.NewTemplate(memoArchivesIndexHeadings)
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
