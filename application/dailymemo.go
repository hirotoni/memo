package application

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hirotoni/memo/components"
	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
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

	b = app.gmw.InsertTextAfter(b, components.HEADING_NAME_TITLE, date)
	b = app.inheritHeading(b, components.HEADING_NAME_TODOS)
	b = app.inheritHeading(b, components.HEADING_NAME_WANTTODOS)
	b = app.appendMemoArchive(b)

	return b
}

// inheritHeading inherits information of the specified heading from previous day's memo
func (app *App) inheritHeading(tb []byte, heading models.Heading) []byte {
	// previous days
	today := time.Now()
	for i := range make([]int, app.Config.DaysToSeek) {
		previousDay := today.AddDate(0, 0, -1*(i+1)).Format(FULL_LAYOUT)
		md, err := app.repos.DailymemoRepo.FindByDate(previousDay)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if i+1 == app.Config.DaysToSeek {
					log.Printf("previous memos were not found in previous %d days.", app.Config.DaysToSeek)
				}
				continue
			}
			log.Fatal(err)
		}

		_, nodesToInsert := app.gmw.FindHeadingAndGetHangingNodes(md.Content, heading)
		tb = app.gmw.InsertNodesAfter(tb, heading, md.Content, nodesToInsert)
		break
	}

	return tb
}

// appendMemoArchive appends memo archives
func (app *App) appendMemoArchive(tb []byte) []byte {
	picked := app.saveMemoArchives(true)

	// insert todays memo archive
	if picked != nil && picked.Destination != "" {
		chosenMemoArchive := markdown.BuildList(markdown.BuildLink(
			picked.Text,
			picked.Destination,
		))
		tb = app.gmw.InsertTextAfter(tb, components.HEADING_NAME_TODAYSMEMOARCHIVE, chosenMemoArchive)
	}

	return tb
}
