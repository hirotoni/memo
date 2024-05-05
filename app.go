package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/hirotoni/memo/markdown"
)

const (
	ENV_DEFAULT_BASE_PATH = "MEMOCLI_BASE_PATH"

	CONFIG_FOLDER_NAME    = ".config/memoapp/"
	DAILYMEMO_FOLDER_NAME = "dailymemo/"

	TEMPLATE_FILE_NAME = "template.md"
	LAYOUT             = "2006-01-02-Mon"

	TODOS_HEADING     = "todos"
	WANTTODOS_HEADING = "wanttodos"
)

var (
	HOME_DIR         = os.Getenv("HOME")
	DEFAULT_BASE_DIR = filepath.Join(HOME_DIR, CONFIG_FOLDER_NAME)            // .config/memoapp/
	DAILYMEMO_DIR    = filepath.Join(DEFAULT_BASE_DIR, DAILYMEMO_FOLDER_NAME) // .config/memoapp/dailymemo/
	TEMPLATE_FILE    = filepath.Join(DAILYMEMO_DIR, TEMPLATE_FILE_NAME)       // .config/memoapp/dailymemo/template.md
)

type AppConfig struct {
	BaseDir      string
	DailymemoDir string
	TemplateFile string
}

func NewAppConfig() AppConfig {
	ac := AppConfig{
		BaseDir:      DEFAULT_BASE_DIR,
		DailymemoDir: DAILYMEMO_DIR,
		TemplateFile: TEMPLATE_FILE,
	}

	if v, found := os.LookupEnv("MEMOAPP_BASE_PATH"); found {
		ac.BaseDir = v
		ac.DailymemoDir = filepath.Join(v, DAILYMEMO_FOLDER_NAME)
		ac.TemplateFile = filepath.Join(ac.DailymemoDir, TEMPLATE_FILE_NAME)
	}

	return ac
}

type App struct {
	config AppConfig
}

func NewApp() App {
	return App{
		config: NewAppConfig(),
	}
}

// Initialize initializes dirs and files
func (c *App) Initialize() {
	// base dir
	_, err := os.Stat(c.config.DailymemoDir)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(c.config.DailymemoDir, 0750); err != nil {
			log.Fatal(err)
		}
		log.Println("memo base directory initialized.")
	}

	// template file
	_, err = os.Stat(c.config.TemplateFile)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(c.config.TemplateFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.WriteString(fmt.Sprintf("# todays memo\n\n## %s\n\n## %s", TODOS_HEADING, WANTTODOS_HEADING))
		log.Println("memo template file initialized.")
	}
}

// OpenTodaysMemo opens today's memo
func (c *App) OpenTodaysMemo() {
	today := time.Now().Format(LAYOUT)
	targetFile := filepath.Join(c.config.DailymemoDir, today+".md")

	_, err := os.Stat(targetFile)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(targetFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		b, err := os.ReadFile(c.config.TemplateFile)
		if err != nil {
			log.Fatal(err)
		}

		// inherit todos from previous memo
		b = c.InheritHeading(f, b, TODOS_HEADING)
		b = c.InheritHeading(f, b, WANTTODOS_HEADING)
		b = c.AppendTips(b)
		f.Write(b)
	}

	// open memo dir with editor
	cmd := exec.Command("code", targetFile, "--folder-uri", c.config.BaseDir)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// InheritHeading inherits todos from previous day's memo
func (c *App) InheritHeading(w *os.File, tb []byte, text string) []byte {
	gmw := markdown.NewGoldmarkWrapper()

	// today
	tDoc := gmw.Parse(tb)
	targetHeader := gmw.GetHeadingNode(tDoc, tb, text, 2)

	// previous day
	// TODO seek for another few days when yesterday memo does not exist
	// TODO handle when previous memos have not been created
	yesterday := time.Now().Add(time.Hour * -24).Format(LAYOUT)
	yb, err := os.ReadFile(filepath.Join(c.config.DailymemoDir, yesterday+".md"))
	if err != nil {
		log.Fatal(err)
	}
	yDoc := gmw.Parse(yb)

	nodesToInsert := gmw.FindHeadingAndGetHangingNodes(yDoc, yb, text, 2)
	tb = gmw.InsertAfter(tDoc, targetHeader, nodesToInsert, tb, yb)

	return tb
}

// AppendTips appends tips
func (c *App) AppendTips(tb []byte) []byte {
	// not yet implemented

	// tips are the things that you want to remember periodically such as
	// - ER diagrams, component diagrams, constants of application you are in charge
	// - product management, development process knowledge
	// - bookmarks, web links
	// - life sayings, someone's sayings

	return tb
}
