package application

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
	"github.com/hirotoni/memo/usecases"
)

// GenerateMemo generates memo file
func (app *App) GenerateMemo(date string, truncate bool) string {
	filename := fmt.Sprintf(FILENAME_FORMAT, date)
	targetFile := filepath.Join(app.Config.DailymemoDir(), filename)

	log.Default().Printf("truncate: %v", truncate)

	_, err := os.Stat(targetFile)
	if errors.Is(err, os.ErrNotExist) || truncate {
		f, err := os.Create(targetFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.Write(app.generateMemo(date))
	}

	return targetFile
}

// generateMemo generates memo file
func (app *App) generateMemo(date string) []byte {
	b, err := os.ReadFile(app.Config.DailymemoTemplateFile())
	if err != nil {
		log.Fatal(err)
	}

	b = app.gmw.InsertTextAfter(b, usecases.HEADING_NAME_TITLE, date)
	b = app.inheritHeading(b, usecases.HEADING_NAME_TODOS)
	b = app.inheritHeading(b, usecases.HEADING_NAME_WANTTODOS)
	b = app.appendTips(b)

	return b
}

// inheritHeading inherits information of the specified heading from previous day's memo
func (app *App) inheritHeading(tb []byte, heading models.Heading) []byte {
	// previous days
	today := time.Now()
	for i := range make([]int, app.Config.DaysToSeek) {
		previousDay := today.AddDate(0, 0, -1*(i+1)).Format(FULL_LAYOUT)
		pb, err := os.ReadFile(filepath.Join(app.Config.DailymemoDir(), previousDay+".md"))
		if errors.Is(err, os.ErrNotExist) {
			if i+1 == app.Config.DaysToSeek {
				log.Printf("previous memos were not found in previous %d days.", app.Config.DaysToSeek)
			}
			continue
		} else if err != nil {
			log.Fatal(err)
		}

		_, nodesToInsert := app.gmw.FindHeadingAndGetHangingNodes(pb, heading)
		tb = app.gmw.InsertNodesAfter(tb, heading, pb, nodesToInsert)
		break
	}

	return tb
}

// appendTips appends tips
func (app *App) appendTips(tb []byte) []byte {
	// tips are the things that you want to remember periodically such as
	// - ER diagrams, component diagrams, constants of application you are in charge
	// - product management, development process knowledge
	// - bookmarks, web links
	// - life sayings, someone's sayings

	picked := app.saveTips(true)

	// insert todays tip
	if picked != nil && picked.Destination != "" {
		chosenTip := markdown.BuildList(markdown.BuildLink(
			picked.Text,
			picked.Destination,
		))
		tb = app.gmw.InsertTextAfter(tb, usecases.HEADING_NAME_TODAYSTIP, chosenTip)
	}

	return tb
}
