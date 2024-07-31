package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
)

var (
	// Subcommands
	// NOTE underscore prefix is used to avoid shadowing by the same name in the package
	_create         = flag.NewFlagSet("create", flag.ExitOnError)
	_createTruncate = _create.Bool("truncate", false, "before creating today's memo, truncate the file if it exists")

	_weekly = flag.NewFlagSet("weekly", flag.ExitOnError)

	_tips = flag.NewFlagSet("tips", flag.ExitOnError)

	// Subcommand descriptions
	SUBCOMMANDS = []struct {
		subcommand *flag.FlagSet
		desc       string
	}{
		{_create, "create today's memo"},
		{_weekly, "generate weekly report"},
		{_tips, "generate tips index"},
	}
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tz, err := time.LoadLocation(TIMEZONE)
	if err != nil {
		panic(err)
	}
	time.Local = tz

	// Custom usage message
	flag.Usage = func() {
		sb := strings.Builder{}
		sb.WriteString("\n")
		sb.WriteString("Usage of ")
		sb.WriteString(flag.CommandLine.Name())
		sb.WriteString(":\n")
		sb.WriteString("  ")
		sb.WriteString(flag.CommandLine.Name())
		sb.WriteString(" [subcommand] [flags]\n")
		sb.WriteString("\n")
		sb.WriteString("Subcommands:\n")
		for _, sc := range SUBCOMMANDS {
			sb.WriteString("  ")
			s := fmt.Sprintf("%-10s", sc.subcommand.Name())
			sb.WriteString(s + "\t\t" + sc.desc + "\n")
			sc.subcommand.VisitAll(func(f *flag.Flag) {
				sb.WriteString("      -")
				sb.WriteString(f.Name)
				sb.WriteString("\t\t\t")
				sb.WriteString(f.Usage)
				sb.WriteString("\n")
			})
		}
		sb.WriteString("\n")
		fmt.Fprint(flag.CommandLine.Output(), sb.String())
	}
	flag.Parse()
}

func main() {
	if len(flag.Args()) == 0 {
		flag.Usage()
		return
	}

	app := NewApp()
	app.Initialize()

	switch flag.Args()[0] {
	case _create.Name():
		_create.Parse(flag.Args()[1:])
		app.OpenTodaysMemo(*_createTruncate)
		app.WeeklyReport(false)
	case _weekly.Name():
		app.WeeklyReport(true)
	case _tips.Name():
		app.SaveTips(true)
	default:
		fmt.Printf("\nInvalid subcommand: %s\n\n", flag.Args()[0])
		flag.Usage()
	}
}
