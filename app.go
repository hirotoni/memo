package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/hirotoni/memo/markdown"
	"github.com/yuin/goldmark/ast"
)

const (
	ENV_DEFAULT_BASE_DIR = "MEMOAPP_BASE_PATH"

	FOLDER_NAME_CONFIG    = ".config/memoapp/"
	FOLDER_NAME_DAILYMEMO = "dailymemo/"

	FILE_NAME_TEMPLATE      = "template.md"
	FILE_NAME_WEEKLY_REPORT = "weekly_report.md"

	LAYOUT   = "2006-01-02-Mon"
	TIMEZONE = "Asia/Tokyo"

	HEADING_NAME_TITLE     = "daily memo"
	HEADING_NAME_TODOS     = "todos"
	HEADING_NAME_WANTTODOS = "wanttodos"
	HEADING_NAME_MEMOS     = "memos"
)

var (
	HOME_DIR         = os.Getenv("HOME")
	DEFAULT_BASE_DIR = filepath.Join(HOME_DIR, FOLDER_NAME_CONFIG)            // .config/memoapp/
	DAILYMEMO_DIR    = filepath.Join(DEFAULT_BASE_DIR, FOLDER_NAME_DAILYMEMO) // .config/memoapp/dailymemo/
	TEMPLATE_FILE    = filepath.Join(DAILYMEMO_DIR, FILE_NAME_TEMPLATE)       // .config/memoapp/dailymemo/template.md
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

	v, found := os.LookupEnv(ENV_DEFAULT_BASE_DIR)
	if found {
		if _, err := os.Stat(v); errors.Is(err, os.ErrNotExist) {
			log.Printf("the directory specified by $%s(%s) does not exist. Using default value(%s).", ENV_DEFAULT_BASE_DIR, v, ac.BaseDir)
			return ac
		}

		ac.BaseDir = v
		ac.DailymemoDir = filepath.Join(v, FOLDER_NAME_DAILYMEMO)
		ac.TemplateFile = filepath.Join(ac.DailymemoDir, FILE_NAME_TEMPLATE)
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

		f.WriteString(fmt.Sprintf("# %s\n\n", HEADING_NAME_TITLE))
		f.WriteString(fmt.Sprintf("## %s\n\n", HEADING_NAME_TODOS))
		f.WriteString(fmt.Sprintf("## %s\n\n", HEADING_NAME_WANTTODOS))
		f.WriteString(fmt.Sprintf("## %s\n\n", HEADING_NAME_MEMOS))

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

		gmw := markdown.NewGoldmarkWrapper()
		doc := gmw.Parse(b)
		targetHeader := gmw.GetHeadingNode(doc, b, HEADING_NAME_TITLE, 1)
		s := targetHeader.Lines().At(0)

		buf := []byte{}
		buf = append(buf, b[:s.Stop]...)
		buf = append(buf, []byte("\n\n"+today)...)
		buf = append(buf, b[s.Stop:]...)
		b = buf

		// inherit todos from previous memo
		b = c.InheritHeading(f, b, HEADING_NAME_TODOS)
		b = c.InheritHeading(f, b, HEADING_NAME_WANTTODOS)
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

func (c *App) WeeklyReport() {
	e, err := os.ReadDir(c.config.DailymemoDir)
	if err != nil {
		log.Fatal(err)
	}

	wantfiles := []string{}
	for _, file := range e {
		format := `\d{4}-\d{2}-\d{2}-\S{3}\.md`
		reg := regexp.MustCompile(format)
		if reg.MatchString(file.Name()) {
			wantfiles = append(wantfiles, filepath.Join(c.config.DailymemoDir, file.Name()))
		}
	}

	f, err := os.Create(filepath.Join(c.config.DailymemoDir, FILE_NAME_WEEKLY_REPORT))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString("# Weekly Report\n\n")

	var curWeekNum int
	for _, fpath := range wantfiles {

		tz, err := time.LoadLocation(TIMEZONE)
		if err != nil {
			log.Fatal(err)
		}
		year, week := time.Now().In(tz).ISOWeek()
		if curWeekNum != week {
			f.WriteString("## " + fmt.Sprint(year) + " | Week " + fmt.Sprint(week) + "\n\n")
			curWeekNum = week
		}

		f.WriteString("### " + filepath.Base(fpath) + "\n\n")

		b, err := os.ReadFile(fpath)
		if err != nil {
			log.Fatal(err)
		}
		gmw := markdown.NewGoldmarkWrapper()
		doc := gmw.Parse(b)
		hangingNodes := gmw.FindHeadingAndGetHangingNodes(doc, b, HEADING_NAME_MEMOS, 2)

		var order = 0
		for _, node := range hangingNodes {
			if node.Kind() == ast.KindHeading {
				relpath, err := filepath.Rel(c.config.DailymemoDir, fpath)
				if err != nil {
					log.Fatal(err)
				}

				format := "%d. [%s](%s#%s)\n"
				title := string(node.Text(b))
				order++
				s := fmt.Sprintf(format, order, title, relpath, title)
				f.WriteString(s)
			}
		}

		if order > 0 {
			f.WriteString("\n")
		}
	}
}
