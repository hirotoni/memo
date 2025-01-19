package usecases

import (
	"strings"

	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
)

var (
	// daily memo
	HEADING_NAME_TITLE             = models.NewHeading(1, "daily memo")
	HEADING_NAME_TODAYSMEMOARCHIVE = models.NewHeading(2, "today's memo archive")
	HEADING_NAME_TODOS             = models.NewHeading(2, "todos")
	HEADING_NAME_WANTTODOS         = models.NewHeading(2, "wanttodos")
	HEADING_NAME_MEMOS             = models.NewHeading(2, "memos")
	// weekly report
	HEADING_NAME_WEEKLYREPORT = models.NewHeading(1, "Weekly Report")
	// memo archives index
	HEADING_NAME_MEMOARCHIVES_INDEX = models.NewHeading(1, "Memo Archives Index")
)

var (
	dailymemoHeadings = []models.Heading{
		HEADING_NAME_TITLE,
		HEADING_NAME_TODAYSMEMOARCHIVE,
		HEADING_NAME_TODOS,
		HEADING_NAME_WANTTODOS,
		HEADING_NAME_MEMOS,
	}
	weeklyReportHeadings = []models.Heading{
		HEADING_NAME_WEEKLYREPORT,
	}
	memoArchivesHeadings = []models.Heading{
		models.NewHeading(1, "sushi (<- memo category)"),
		models.NewHeading(2, "how to eat sushi (<- memo title in heading level 2)"),
		models.NewHeading(2, "how to roll sushi (<- another memo)"),
	}
	memoArchivesIndexHeadings = []models.Heading{
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
