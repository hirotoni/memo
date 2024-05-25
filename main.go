package main

import (
	"flag"
	"log"
	"time"
)

var (
	trancate = flag.Bool("trancate", false, "trancate todays file")
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tz, err := time.LoadLocation(TIMEZONE)
	if err != nil {
		panic(err)
	}
	time.Local = tz

	flag.Parse()
}

func main() {
	app := NewApp()
	app.Initialize()
	app.OpenTodaysMemo(*trancate)
	app.WeeklyReport()
}
