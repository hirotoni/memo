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

	TIMEZONE     = "Asia/Tokyo"
	LAYOUT       = "2006-01-02-Mon"
	LAYOUT_REGEX = `\d{4}-\d{2}-\d{2}-\S{3}\.md`

	// number of dates to seek back when inheriting todos from previous days
	DAYS_TO_SEEK = 10
)

var (
	HOME_DIR         = os.Getenv("HOME")
	DEFAULT_BASE_DIR = filepath.Join(HOME_DIR, FOLDER_NAME_CONFIG) // .config/memoapp/
)

type AppConfig struct {
	baseDir string
}

func NewAppConfig() AppConfig {
	ac := AppConfig{
		baseDir: DEFAULT_BASE_DIR,
	}

	v, found := os.LookupEnv(ENV_NAME_DEFAULT_BASE_DIR)
	if found {
		if _, err := os.Stat(v); errors.Is(err, os.ErrNotExist) {
			log.Printf("the directory specified by $%s(%s) does not exist. Using default value(%s).", ENV_NAME_DEFAULT_BASE_DIR, v, ac.BaseDir())
			return ac
		}

		ac.baseDir = v
	}

	return ac
}

func (ac *AppConfig) BaseDir() string {
	return ac.baseDir
}
func (ac *AppConfig) DailymemoDir() string {
	return filepath.Join(ac.baseDir, FOLDER_NAME_DAILYMEMO) // {basedir}/dailymemo
}
func (ac *AppConfig) DailymemoTemplateFile() string {
	return filepath.Join(ac.baseDir, FOLDER_NAME_DAILYMEMO, FILE_NAME_DAILYMEMO_TEMPLATE) // {basedir}/dailymemo/template.md
}
func (ac *AppConfig) WeeklyReportFile() string {
	return filepath.Join(ac.baseDir, FOLDER_NAME_DAILYMEMO, FILE_NAME_WEEKLY_REPORT) // {basedir}/dailymemo/weekly_report.md
}
func (ac *AppConfig) TipsDir() string {
	return filepath.Join(ac.baseDir, FOLDER_NAME_TIPS) // {basedir}/tips
}
func (ac *AppConfig) TipsTemplateFile() string {
	return filepath.Join(ac.baseDir, FOLDER_NAME_TIPS, FILE_NAME_TIPS_TEMPLATE) // {basedir}/tips/template.md
}
func (ac *AppConfig) TipsIndexFile() string {
	return filepath.Join(ac.baseDir, FOLDER_NAME_TIPS, FILE_NAME_TIPS_INDEX) // {basedir}/tips/index.md
}
