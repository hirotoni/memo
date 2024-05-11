package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hirotoni/memo/markdown"
	"github.com/yuin/goldmark/ast"
)

const (
	ENV_NAME_DEFAULT_BASE_DIR = "MEMOAPP_BASE_DIR"

	FOLDER_NAME_CONFIG    = ".config/memoapp/"
	FOLDER_NAME_DAILYMEMO = "dailymemo/"
	FOLDER_NAME_TIPS      = "tips"

	FILE_NAME_DAILYMEMO_TEMPLATE = "template.md"
	FILE_NAME_TIPS_TEMPLATE      = "template.md"
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

var (
	HOME_DIR                = os.Getenv("HOME")
	DEFAULT_BASE_DIR        = filepath.Join(HOME_DIR, FOLDER_NAME_CONFIG)                                          // .config/memoapp/
	DAILYMEMO_DIR           = filepath.Join(DEFAULT_BASE_DIR, FOLDER_NAME_DAILYMEMO)                               // .config/memoapp/dailymemo/
	DAILYMEMO_TEMPLATE_FILE = filepath.Join(DEFAULT_BASE_DIR, FOLDER_NAME_DAILYMEMO, FILE_NAME_DAILYMEMO_TEMPLATE) // .config/memoapp/dailymemo/template.md
	TIPS_DIR                = filepath.Join(DEFAULT_BASE_DIR, FOLDER_NAME_TIPS)                                    // .config/memoapp/tips/
	TIPS_TEMPLATE_FILE      = filepath.Join(DEFAULT_BASE_DIR, FOLDER_NAME_TIPS, FILE_NAME_TIPS_TEMPLATE)           // .config/memoapp/tips/template.md
)

type AppConfig struct {
	BaseDir               string
	DailymemoDir          string
	DailymemoTemplateFile string
	TipsDir               string
	TipsTemplateFile      string
}

func NewAppConfig() AppConfig {
	ac := AppConfig{
		BaseDir:               DEFAULT_BASE_DIR,
		DailymemoDir:          DAILYMEMO_DIR,
		DailymemoTemplateFile: DAILYMEMO_TEMPLATE_FILE,
		TipsDir:               TIPS_DIR,
		TipsTemplateFile:      TIPS_TEMPLATE_FILE,
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
		ac.TipsDir = filepath.Join(ac.BaseDir, FOLDER_NAME_TIPS)
		ac.TipsTemplateFile = filepath.Join(ac.BaseDir, FOLDER_NAME_TIPS, FILE_NAME_TIPS_TEMPLATE)
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
	// dailymemo dir
	_, err := os.Stat(c.config.DailymemoDir)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(c.config.DailymemoDir, 0750); err != nil {
			log.Fatal(err)
		}
		log.Printf("memo directory initialized: %s", c.config.BaseDir)
	}
	// dailymemo template file
	_, err = os.Stat(c.config.DailymemoTemplateFile)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(c.config.DailymemoTemplateFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.WriteString(fmt.Sprintf("# %s\n\n", HEADING_NAME_TITLE))
		f.WriteString(fmt.Sprintf("## %s\n\n", HEADING_NAME_TODOS))
		f.WriteString(fmt.Sprintf("## %s\n\n", HEADING_NAME_WANTTODOS))
		f.WriteString(fmt.Sprintf("## %s\n\n", HEADING_NAME_MEMOS))

		log.Printf("dailymemo template file initialized: %s", c.config.DailymemoTemplateFile)
	}

	// tips dir
	_, err = os.Stat(c.config.TipsDir)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(c.config.TipsDir, 0750); err != nil {
			log.Fatal(err)
		}
		log.Printf("tips directory initialized: %s", c.config.TipsDir)
	}
	// tips template file
	_, err = os.Stat(c.config.TipsTemplateFile)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(c.config.TipsTemplateFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.WriteString(fmt.Sprintf("# %s\n\n", "template (<- FILENAME HERE)"))
		f.WriteString(fmt.Sprintf("## %s\n\n", "how to eat sushi (<- YOUR TIPS HERE)"))

		log.Printf("tips template file initialized: %s", c.config.TipsTemplateFile)
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

		b, err := os.ReadFile(c.config.DailymemoTemplateFile)
		if err != nil {
			log.Fatal(err)
		}

		gmw := markdown.NewGoldmarkWrapper()

		// inherit todos from previous memo
		b = c.InheritHeading(b, HEADING_NAME_TODOS)
		b = c.InheritHeading(b, HEADING_NAME_WANTTODOS)
		b = c.AppendTips(b)

		doc := gmw.Parse(b)
		targetHeader := gmw.GetHeadingNode(doc, b, HEADING_NAME_TITLE, 1)
		b = gmw.InsertTextAfter(doc, targetHeader, today, b)

		f.Write(b)
	}

	// open memo dir with editor
	cmd := exec.Command("code", targetFile, "--folder-uri", c.config.BaseDir)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// InheritHeading inherits todos from previous day's memo
func (c *App) InheritHeading(tb []byte, text string) []byte {
	gmw := markdown.NewGoldmarkWrapper()

	// today
	tDoc := gmw.Parse(tb)
	targetHeader := gmw.GetHeadingNode(tDoc, tb, text, 2)

	// previous days
	tz, err := time.LoadLocation(TIMEZONE)
	if err != nil {
		log.Fatal(err)
	}
	today := time.Now().In(tz)
	for i := range make([]int, DAYS_TO_SEEK) {
		previousDay := today.AddDate(0, 0, -1*(i+1)).Format(LAYOUT)
		pb, err := os.ReadFile(filepath.Join(c.config.DailymemoDir, previousDay+".md"))
		if errors.Is(err, os.ErrNotExist) {
			if i+1 == DAYS_TO_SEEK {
				log.Printf("previous memos were not found in previous %d days.", DAYS_TO_SEEK)
			}
			continue
		} else if err != nil {
			log.Fatal(err)
		}

		pDoc := gmw.Parse(pb)

		nodesToInsert := gmw.FindHeadingAndGetHangingNodes(pDoc, pb, text, 2)
		tb = gmw.InsertAfter(tDoc, targetHeader, nodesToInsert, tb, pb)
		break
	}

	return tb
}

// AppendTips appends tips
func (c *App) AppendTips(tb []byte) []byte {
	// not yet fully implemented

	// tips are the things that you want to remember periodically such as
	// - ER diagrams, component diagrams, constants of application you are in charge
	// - product management, development process knowledge
	// - bookmarks, web links
	// - life sayings, someone's sayings

	gmw := markdown.NewGoldmarkWrapper()

	var poolTipsToShow []string
	// var poolTipsAlreadyShown []string
	var targetTipFiles []string
	var fn = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || path == c.config.TipsTemplateFile {
			return nil
		}

		targetTipFiles = append(targetTipFiles, path)
		return nil
	}

	if err := filepath.Walk(c.config.TipsDir, fn); err != nil {
		log.Fatal(err)
	}

	for _, v := range targetTipFiles {
		b, err := os.ReadFile(v)
		if err != nil {
			log.Fatal(err)
		}
		doc := gmw.Parse(b)
		headings := gmw.GetHeadingNodes(doc, b, 2)

		relpath, err := filepath.Rel(c.config.DailymemoDir, v)
		if err != nil {
			log.Fatal(err)
		}
		for _, vv := range headings {
			poolTipsToShow = append(poolTipsToShow, fmt.Sprintf("- [%s](%s#%s)\n", string(vv.Text(b)), relpath, string(vv.Text(b))))
		}
	}

	fmt.Println(strings.Join(poolTipsToShow, ""))

	doc := gmw.Parse(tb)
	targetHeader := gmw.GetHeadingNode(doc, tb, HEADING_NAME_TITLE, 1)
	tb = gmw.InsertTextAfter(doc, targetHeader, strings.Join(poolTipsToShow, ""), tb)

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

		datestring, found := strings.CutSuffix(filepath.Base(fpath), ".md")
		if !found {
			log.Fatal("failed to cut suffix.")
		}

		date, err := time.Parse(LAYOUT, datestring)
		if err != nil {
			log.Fatal(err)
		}

		year, week := date.ISOWeek()
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
			if n, ok := node.(*ast.Heading); ok {
				relpath, err := filepath.Rel(c.config.DailymemoDir, fpath)
				if err != nil {
					log.Fatal(err)
				}

				var format = "%d. [%s](%s#%s)\n"
				title := strings.Repeat("#", n.Level-2) + " " + string(node.Text(b))
				tag := strings.ReplaceAll(string(node.Text(b)), " ", "-")
				tag = strings.ReplaceAll(tag, "#", "")
				fullwidthchars := strings.Split("　！＠＃＄％＾＆＊（）＋｜〜＝￥｀「」｛｝；’：”、。・＜＞？【】『』《》〔〕［］‹›«»〘〙〚〛", "")
				for _, c := range fullwidthchars {
					tag = strings.ReplaceAll(tag, c, "")
				}

				order++
				s := fmt.Sprintf(format, order, title, relpath, tag)
				f.WriteString(s)
			}
		}

		if order > 0 {
			f.WriteString("\n")
		}
	}
}
