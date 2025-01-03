package application

import (
	"errors"
	"log"
	"os"
	"os/exec"

	"github.com/hirotoni/memo/config"
	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
	"github.com/hirotoni/memo/repos"
	"github.com/hirotoni/memo/usecases"
)

const (
	TIMEZONE        = "Asia/Tokyo"
	FULL_LAYOUT     = "2006-01-02-Mon"
	SHORT_LAYOUT    = "2006-01-02"
	FILENAME_FORMAT = "%s.md"
)

type App struct {
	Config *config.TomlConfig
	gmw    *markdown.GoldmarkWrapper
	repos  *repos.Repos
}

func NewApp() App {
	gmw := markdown.NewGoldmarkWrapper()
	config := config.LoadTomlConfig()

	return App{
		gmw:    gmw,
		Config: config,
		repos:  repos.NewRepos(config, gmw),
	}
}

// Initialize initializes dirs and files
func (app *App) Initialize() {
	// dailymemo
	initializeDir(app.Config.DailymemoDir())
	initializeFile(app.Config.DailymemoTemplateFile(), usecases.TemplateDailymemo)
	// tips
	initializeDir(app.Config.TipsDir())
	initializeFile(app.Config.TipsTemplateFile(), usecases.TemplateTips)
	initializeFile(app.Config.TipsIndexFile(), usecases.TemplateTipsIndex)
}

func initializeFile(filepath string, template models.Template) {
	_, err := os.Stat(filepath)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.WriteString(usecases.GenerateTemplateString(template))

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
