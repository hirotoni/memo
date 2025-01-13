package models

import (
	"time"
)

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
