package application

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/usecases"
	"github.com/yuin/goldmark/ast"
)

// WeeklyReport generates weekly report file
func (app *App) WeeklyReport() {
	wr := app.buildWeeklyReport()

	f, err := os.Create(app.Config.WeeklyReportFile())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(usecases.GenerateTemplateString(usecases.TemplateWeeklyReport) + "\n")
	f.WriteString(wr)
}

func weekSpliter(date time.Time) string {
	year, week := date.ISOWeek()
	return "## " + fmt.Sprint(year) + " | Week " + fmt.Sprint(week) + "\n\n"
}

// buildWeeklyReport builds weekly report
func (app *App) buildWeeklyReport() string {
	var sb strings.Builder
	var curWeekNum int

	dms := app.repos.DailymemoRepo.Entires()
	for _, dm := range dms {
		if curWeekNum != dm.WeekNum() {
			sb.WriteString(weekSpliter(dm.Date))
			curWeekNum = dm.WeekNum()
		}

		sb.WriteString(markdown.BuildHeading(3, dm.BaseName+"\n\n"))
		_, hangingNodes := app.gmw.FindHeadingAndGetHangingNodes(dm.Content, usecases.HEADING_NAME_MEMOS)

		var order = 0
		for _, node := range hangingNodes {
			if n, ok := node.(*ast.Heading); ok {
				relpath, err := filepath.Rel(app.Config.DailymemoDir(), dm.Filepath)
				if err != nil {
					log.Fatal(err)
				}

				order++

				title := markdown.BuildHeading(n.Level-2, string(node.Text(dm.Content)))
				tag := markdown.Text2tag(string(node.Text(dm.Content)))
				s := markdown.BuildOrderedList(order, markdown.BuildLink(title, relpath+"#"+tag)) + "\n"
				sb.WriteString(s)
			}
		}

		if order > 0 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
