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

	// Subcommand descriptions
	SUBCOMMANDS = map[*flag.FlagSet]string{
		_create: "create today's memo",
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
		for sc, desc := range SUBCOMMANDS {
			sb.WriteString("  ")
			sb.WriteString(sc.Name() + "\t" + desc)
		}
		sb.WriteString("\n")
		sb.WriteString("\n")
		sb.WriteString("Use \"")
		sb.WriteString(flag.CommandLine.Name())
		sb.WriteString(" [subcommand] -h\" for more information about a subcommand.\n")
		sb.WriteString("\n")
		fmt.Fprintf(flag.CommandLine.Output(), sb.String())
	}
	for sc, desc := range SUBCOMMANDS {
		sc.Usage = func() {
			sb := strings.Builder{}
			sb.WriteString("\n")
			sb.WriteString("Subcommand: " + sc.Name())
			sb.WriteString("\n")
			sb.WriteString("  ")
			sb.WriteString(desc)
			sb.WriteString("\n")
			sb.WriteString("\n")
			sb.WriteString("Flags:\n")
			sc.VisitAll(func(f *flag.Flag) {
				sb.WriteString("  -")
				sb.WriteString(f.Name)
				sb.WriteString("\t")
				sb.WriteString(f.Usage)
				sb.WriteString("\n")
			})
			sb.WriteString("\n")
			fmt.Fprintf(flag.CommandLine.Output(), sb.String())
		}
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
		app.WeeklyReport()
	default:
		fmt.Printf("\nInvalid subcommand: %s\n\n", flag.Args()[0])
		flag.Usage()
	}
}
