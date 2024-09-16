package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tz, err := time.LoadLocation(TIMEZONE)
	if err != nil {
		panic(err)
	}
	time.Local = tz
}

func main() {
	app := NewApp()
	app.Initialize()

	cliapp := &cli.App{
		EnableBashCompletion: true,
		Name:                 "memo",
		Usage:                "A CLI tool for managing daily memo",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "create memo",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "truncate",
						Aliases: []string{"t"},
						Usage:   "before creating memo, truncate the file if it exists",
					},
				},
				Action: func(c *cli.Context) error {
					today := time.Now().Format(LAYOUT)
					targetFile := app.GenerateMemo(today, c.Bool("truncate"))
					app.WeeklyReport()
					app.OpenEditor(targetFile)
					return nil
				},
			},
			{
				Name:  "weekly",
				Usage: "generate weekly report",
				Action: func(c *cli.Context) error {
					app.WeeklyReport()
					app.OpenEditor(app.config.WeeklyReportFile())
					return nil
				},
			},
			{
				Name:  "tips",
				Usage: "generate tips index",
				Action: func(c *cli.Context) error {
					app.SaveTips()
					app.OpenEditor(app.config.TipsIndexFile())
					return nil
				},
			},
			{
				Name:  "env",
				Usage: "print environment information",
				Action: func(c *cli.Context) error {
					app.ShowEnv()
					return nil
				},
			},
		},
	}

	if err := cliapp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
