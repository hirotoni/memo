package application

import (
	"errors"
	"log"
	"os"
	"os/exec"

	"github.com/hirotoni/memo/components"
	"github.com/hirotoni/memo/configs"
	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
	"github.com/hirotoni/memo/repos"
)

const (
	TIMEZONE        = "Asia/Tokyo"
	FULL_LAYOUT     = "2006-01-02-Mon"
	SHORT_LAYOUT    = "2006-01-02"
	FILENAME_FORMAT = "%s.md"
)

type App struct {
	Config *configs.TomlConfig
	gmw    *markdown.GoldmarkWrapper
	repos  *repos.Repos
}

func NewApp() App {
	gmw := markdown.NewGoldmarkWrapper()
	conf := configs.LoadTomlConfig()
	return App{
		gmw:    gmw,
		Config: conf,
		repos:  repos.NewRepos(conf),
	}
}

func (app *App) WithCustomConfig(conf configs.TomlConfig) {
	app.Config = &conf
	app.repos = repos.NewRepos(&conf)
}

// Initialize initializes dirs and files
func (app *App) Initialize() {
	// dailymemo
	initializeDir(app.Config.DailymemoDir())
	initializeFile(app.Config.DailymemoTemplateFile(), components.TemplateDailymemo)
	// memoarchives
	initializeDir(app.Config.MemoArchivesDir())
	initializeFile(app.Config.MemoArchivesTemplateFile(), components.TemplateMemoArchives)
	initializeFile(app.Config.MemoArchivesIndexFile(), components.TemplateMemoArchivesIndex)
}

func initializeFile(filepath string, template models.Template) {
	_, err := os.Stat(filepath)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.WriteString(components.GenerateTemplateString(template))

		log.Printf("file initialized: %s", filepath)
	}
}

func initializeDir(dirpath string) {
	_, err := os.Stat(dirpath)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(dirpath, 0750); err != nil {
			log.Fatal(err)
		}
		log.Printf("directory initialized: %s", dirpath)
	}
}

// OpenEditor opens editor
func (app *App) OpenEditor(path string) {
	cmd := exec.Command("code", path, "--folder-uri", app.Config.BaseDir)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
