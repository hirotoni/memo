package main

import (
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
