package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

const (
	ENV_NAME_DEFAULT_BASE_DIR = "MEMOAPP_BASE_DIR"

	FOLDER_NAME_CONFIG    = ".config/memoapp/"
	FOLDER_NAME_DAILYMEMO = "dailymemo/"
	FOLDER_NAME_TIPS      = "tips"

	FILE_NAME_DAILYMEMO_TEMPLATE = "template.md"
	FILE_NAME_TIPS_TEMPLATE      = "template.md"
	FILE_NAME_TIPS_INDEX         = "index.md"
	FILE_NAME_WEEKLY_REPORT      = "weekly_report.md"

	LAYOUT   = "2006-01-02-Mon"
	TIMEZONE = "Asia/Tokyo"

	HEADING_NAME_TITLE     = "daily memo"
	HEADING_NAME_TODOS     = "todos"
	HEADING_NAME_WANTTODOS = "wanttodos"
	HEADING_NAME_MEMOS     = "memos"

	// number of dates to seek back when inheriting todos from previous days
	DAYS_TO_SEEK = 10
)

type Heading struct {
	text  string
	level int
}

var DailyemoTemplate = struct {
	Title     Heading
	Todos     Heading
	WantToDos Heading
	Memos     Heading
}{
	Title:     Heading{text: HEADING_NAME_TITLE, level: 1},
	Todos:     Heading{text: HEADING_NAME_TODOS, level: 2},
	WantToDos: Heading{text: HEADING_NAME_WANTTODOS, level: 2},
	Memos:     Heading{text: HEADING_NAME_MEMOS, level: 2},
}

var (
	HOME_DIR                = os.Getenv("HOME")
	DEFAULT_BASE_DIR        = filepath.Join(HOME_DIR, FOLDER_NAME_CONFIG)                                          // .config/memoapp/
	DAILYMEMO_DIR           = filepath.Join(DEFAULT_BASE_DIR, FOLDER_NAME_DAILYMEMO)                               // .config/memoapp/dailymemo/
	DAILYMEMO_TEMPLATE_FILE = filepath.Join(DEFAULT_BASE_DIR, FOLDER_NAME_DAILYMEMO, FILE_NAME_DAILYMEMO_TEMPLATE) // .config/memoapp/dailymemo/template.md
	WEEKLY_REPORT_FILE      = filepath.Join(DEFAULT_BASE_DIR, FOLDER_NAME_DAILYMEMO, FILE_NAME_WEEKLY_REPORT)      // .config/memoapp/dailymemo/weekly_report.md
	TIPS_DIR                = filepath.Join(DEFAULT_BASE_DIR, FOLDER_NAME_TIPS)                                    // .config/memoapp/tips/
	TIPS_TEMPLATE_FILE      = filepath.Join(DEFAULT_BASE_DIR, FOLDER_NAME_TIPS, FILE_NAME_TIPS_TEMPLATE)           // .config/memoapp/tips/template.md
	TIPS_INDEX_FILE         = filepath.Join(DEFAULT_BASE_DIR, FOLDER_NAME_TIPS, FILE_NAME_TIPS_INDEX)              // .config/memoapp/tips/index.md
)

type AppConfig struct {
	BaseDir               string
	DailymemoDir          string
	DailymemoTemplateFile string
	WeeklyReportFile      string
	TipsDir               string
	TipsTemplateFile      string
	TipsIndexFile         string
}

func NewAppConfig() AppConfig {
	ac := AppConfig{
		BaseDir:               DEFAULT_BASE_DIR,
		DailymemoDir:          DAILYMEMO_DIR,
		DailymemoTemplateFile: DAILYMEMO_TEMPLATE_FILE,
		WeeklyReportFile:      WEEKLY_REPORT_FILE,
		TipsDir:               TIPS_DIR,
		TipsTemplateFile:      TIPS_TEMPLATE_FILE,
		TipsIndexFile:         TIPS_INDEX_FILE,
	}

	v, found := os.LookupEnv(ENV_NAME_DEFAULT_BASE_DIR)
	if found {
		if _, err := os.Stat(v); errors.Is(err, os.ErrNotExist) {
			log.Printf("the directory specified by $%s(%s) does not exist. Using default value(%s).", ENV_NAME_DEFAULT_BASE_DIR, v, ac.BaseDir)
			return ac
		}

		ac.BaseDir = v
		ac.DailymemoDir = filepath.Join(ac.BaseDir, FOLDER_NAME_DAILYMEMO)
		ac.DailymemoTemplateFile = filepath.Join(ac.BaseDir, FOLDER_NAME_DAILYMEMO, FILE_NAME_DAILYMEMO_TEMPLATE)
		ac.WeeklyReportFile = filepath.Join(ac.BaseDir, FOLDER_NAME_DAILYMEMO, FILE_NAME_WEEKLY_REPORT)
		ac.TipsDir = filepath.Join(ac.BaseDir, FOLDER_NAME_TIPS)
		ac.TipsTemplateFile = filepath.Join(ac.BaseDir, FOLDER_NAME_TIPS, FILE_NAME_TIPS_TEMPLATE)
		ac.TipsIndexFile = filepath.Join(ac.BaseDir, FOLDER_NAME_TIPS, FILE_NAME_TIPS_INDEX)
	}

	return ac
}
