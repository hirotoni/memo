package models

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var FULL_LAYOUT = "2006-01-02-Mon"

type Dailymemo struct {
	Filepath string
	BaseName string
	Date     time.Time
	Content  []byte
}

func (dm Dailymemo) YearNum() int {
	year, _ := dm.Date.ISOWeek()
	return year
}
func (dm Dailymemo) WeekNum() int {
	_, week := dm.Date.ISOWeek()
	return week
}

func NewDailymemoFromFilepath(fpath string) Dailymemo {
	basename := filepath.Base(fpath)
	datestring, found := strings.CutSuffix(basename, ".md")
	if !found {
		log.Fatal("failed to cut suffix.")
	}
	date, err := time.Parse(FULL_LAYOUT, datestring)
	if err != nil {
		log.Fatal(err)
	}
	b, err := os.ReadFile(fpath)
	if err != nil {
		log.Fatal(err)
	}

	return Dailymemo{
		Filepath: fpath,
		BaseName: basename,
		Date:     date,
		Content:  b,
	}
}
