package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/hirotoni/memo/markdown"
)

const (
	FOLDER_NAME_CONFIG    = ".config/memoapp/"
	FOLDER_NAME_DAILYMEMO = "dailymemo/"
	FOLDER_NAME_TIPS      = "tips"

	FILE_NAME_CONFIG             = "config.toml"
	FILE_NAME_DAILYMEMO_TEMPLATE = "template.md"
	FILE_NAME_TIPS_TEMPLATE      = "template.md"
	FILE_NAME_TIPS_INDEX         = "index.md"
	FILE_NAME_WEEKLY_REPORT      = "weekly_report.md"
)

type TomlConfig struct {
	BaseDir    string `toml:"basedir"`    // memoapp base directory
	DaysToSeek int    `toml:"daystoseek"` // days to seek back
	Gmw        *markdown.GoldmarkWrapper
}

func NewTomlConfig(baseDir string, daystoseek int, gmw *markdown.GoldmarkWrapper) *TomlConfig {
	return &TomlConfig{
		BaseDir:    baseDir,
		DaysToSeek: daystoseek,
		Gmw:        gmw,
	}
}

func LoadTomlConfig() *TomlConfig {
	var tomlConfig = &TomlConfig{}

	configDirPath, err := ConfigDirPath()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	if _, err := os.Stat(configDirPath); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(configDirPath, 0755); err != nil {
			log.Fatal(err)
			return nil
		}
	}

	configFilePath, err := ConfigFilePath()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	if _, err := os.Stat(configFilePath); err != nil {
		f, err := os.Create(configFilePath)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		defer f.Close()

		tomlConfig.BaseDir = configDirPath
		tomlConfig.DaysToSeek = 10

		enc := toml.NewEncoder(f)
		if err := enc.Encode(tomlConfig); err != nil {
			log.Fatal(err)
			return nil
		}
	}

	_, err = toml.DecodeFile(configFilePath, tomlConfig)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	tomlConfig.Gmw = markdown.NewGoldmarkWrapper()

	return tomlConfig
}

func ConfigDirPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, FOLDER_NAME_CONFIG), nil
}

func ConfigFilePath() (string, error) {
	configDir, err := ConfigDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, FILE_NAME_CONFIG), nil
}

func (tc *TomlConfig) DailymemoDir() string {
	return filepath.Join(tc.BaseDir, FOLDER_NAME_DAILYMEMO) // {basedir}/dailymemo
}
func (tc *TomlConfig) DailymemoTemplateFile() string {
	return filepath.Join(tc.BaseDir, FOLDER_NAME_DAILYMEMO, FILE_NAME_DAILYMEMO_TEMPLATE) // {basedir}/dailymemo/template.md
}
func (tc *TomlConfig) WeeklyReportFile() string {
	return filepath.Join(tc.BaseDir, FOLDER_NAME_DAILYMEMO, FILE_NAME_WEEKLY_REPORT) // {basedir}/dailymemo/weekly_report.md
}
func (tc *TomlConfig) TipsDir() string {
	return filepath.Join(tc.BaseDir, FOLDER_NAME_TIPS) // {basedir}/tips
}
func (tc *TomlConfig) TipsTemplateFile() string {
	return filepath.Join(tc.BaseDir, FOLDER_NAME_TIPS, FILE_NAME_TIPS_TEMPLATE) // {basedir}/tips/template.md
}
func (tc *TomlConfig) TipsIndexFile() string {
	return filepath.Join(tc.BaseDir, FOLDER_NAME_TIPS, FILE_NAME_TIPS_INDEX) // {basedir}/tips/index.md
}
