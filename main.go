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
				Name:  "new",
				Usage: "create memo",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "date",
						Aliases:     []string{"d"},
						Usage:       "specify the date to create memo: `YYYY-MM-DD`",
						DefaultText: "today",
					},
					&cli.BoolFlag{
						Name:    "truncate",
						Aliases: []string{"t"},
						Usage:   "before creating memo, truncate the file if it exists",
					},
				},
				Action: func(c *cli.Context) error {
					var date string

					arg := c.String("date")
					if arg != "" {
						d, err := time.Parse(SHORT_LAYOUT, arg)
						if err != nil {
							log.Fatalf("Invalid date format: %s", arg)
						}
						date = d.Format(FULL_LAYOUT)
					} else {
						date = time.Now().Format(FULL_LAYOUT) // default to today
					}

					targetFile := app.GenerateMemo(date, c.Bool("truncate"))
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
				Name:  "config",
				Usage: "edit configuration information",
				Action: func(c *cli.Context) error {
					app.EditConfig()
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:  "show",
						Usage: "show configuration information",
						Action: func(c *cli.Context) error {
							app.ShowConfig()
							return nil
						},
					},
				},
			},
		},
	}

	if err := cliapp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
