package main

import (
	"log"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tz, err := time.LoadLocation(TIMEZONE)
	if err != nil {
		panic(err)
	}
	time.Local = tz
}

// var (
// 	flag1 = flag.Bool("flag1", false, "flag 1")
// 	flag2 = flag.Bool("flag2", false, "flag 2")
// )

func main() {
	app := NewApp()
	app.Initialize()
	app.OpenTodaysMemo()
	app.WeeklyReport()
}
