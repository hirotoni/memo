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

type App struct {
	gmw    *markdown.GoldmarkWrapper
	config AppConfig
}

func NewApp() App {
	return App{
		gmw:    markdown.NewGoldmarkWrapper(),
		config: NewAppConfig(),
	}
}

// Initialize initializes dirs and files
func (app *App) Initialize() {
	// dailymemo dir
	_, err := os.Stat(app.config.DailymemoDir)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(app.config.DailymemoDir, 0750); err != nil {
			log.Fatal(err)
		}
		log.Printf("memo directory initialized: %s", app.config.BaseDir)
	}
	// dailymemo template file
	_, err = os.Stat(app.config.DailymemoTemplateFile)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(app.config.DailymemoTemplateFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.WriteString(dailymemoTemplate)

		log.Printf("dailymemo template file initialized: %s", app.config.DailymemoTemplateFile)
	}

	// tips dir
	_, err = os.Stat(app.config.TipsDir)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(app.config.TipsDir, 0750); err != nil {
			log.Fatal(err)
		}
		log.Printf("tips directory initialized: %s", app.config.TipsDir)
	}
	// tips template file
	_, err = os.Stat(app.config.TipsTemplateFile)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(app.config.TipsTemplateFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.WriteString(tipsTemplate)

		log.Printf("tips template file initialized: %s", app.config.TipsTemplateFile)
	}
	// tips index file
	_, err = os.Stat(app.config.TipsIndexFile)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(app.config.TipsIndexFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.WriteString(tipsIndexTemplate)

		log.Printf("tips index file initialized: %s", app.config.TipsIndexFile)
	}
}

// OpenTodaysMemo opens today's memo
func (app *App) OpenTodaysMemo(trancate bool) {
	today := time.Now().Format(LAYOUT)
	targetFile := filepath.Join(app.config.DailymemoDir, today+".md")

	log.Default().Printf("trancate: %v", trancate)

	_, err := os.Stat(targetFile)
	if errors.Is(err, os.ErrNotExist) || trancate {
		f, err := os.Create(targetFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		b, err := os.ReadFile(app.config.DailymemoTemplateFile)
		if err != nil {
			log.Fatal(err)
		}

		// inherit todos from previous memo
		b = app.InheritHeading(b, HEADING_NAME_TODOS)
		b = app.InheritHeading(b, HEADING_NAME_WANTTODOS)
		b = app.AppendTips(b)

		doc := app.gmw.Parse(b)
		targetHeader := app.gmw.GetHeadingNode(doc, b, HEADING_NAME_TITLE, 1)
		b = app.gmw.InsertTextAfter(doc, targetHeader, today, b)

		f.Write(b)
	}

	// open memo dir with editor
	cmd := exec.Command("code", targetFile, "--folder-uri", app.config.BaseDir)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// InheritHeading inherits todos from previous day's memo
func (app *App) InheritHeading(tb []byte, text string) []byte {
	// today
	tDoc := app.gmw.Parse(tb)
	targetHeader := app.gmw.GetHeadingNode(tDoc, tb, text, 2)

	// previous days
	today := time.Now()
	for i := range make([]int, DAYS_TO_SEEK) {
		previousDay := today.AddDate(0, 0, -1*(i+1)).Format(LAYOUT)
		pb, err := os.ReadFile(filepath.Join(app.config.DailymemoDir, previousDay+".md"))
		if errors.Is(err, os.ErrNotExist) {
			if i+1 == DAYS_TO_SEEK {
				log.Printf("previous memos were not found in previous %d days.", DAYS_TO_SEEK)
			}
			continue
		} else if err != nil {
			log.Fatal(err)
		}

		pDoc := app.gmw.Parse(pb)

		nodesToInsert := app.gmw.FindHeadingAndGetHangingNodes(pDoc, pb, text, 2)
		tb = app.gmw.InsertAfter(tDoc, targetHeader, nodesToInsert, tb, pb)
		break
	}

	return tb
}

// AppendTips appends tips
func (app *App) AppendTips(tb []byte) []byte {
	// not yet fully implemented

	// tips are the things that you want to remember periodically such as
	// - ER diagrams, component diagrams, constants of application you are in charge
	// - product management, development process knowledge
	// - bookmarks, web links
	// - life sayings, someone's sayings

	var poolTipsToShow []string
	// var poolTipsAlreadyShown []string
	var targetTipFiles []string
	var fn = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || path == app.config.TipsTemplateFile {
			return nil
		}

		targetTipFiles = append(targetTipFiles, path)
		return nil
	}

	if err := filepath.Walk(app.config.TipsDir, fn); err != nil {
		log.Fatal(err)
	}

	for _, v := range targetTipFiles {
		b, err := os.ReadFile(v)
		if err != nil {
			log.Fatal(err)
		}
		doc := app.gmw.Parse(b)
		headings := app.gmw.GetHeadingNodes(doc, b, 2)

		relpath, err := filepath.Rel(app.config.DailymemoDir, v)
		if err != nil {
			log.Fatal(err)
		}
		for _, vv := range headings {
			title := string(vv.Text(b))
			tag := text2tag(title)
			poolTipsToShow = append(poolTipsToShow, fmt.Sprintf("- [%s](%s#%s)\n", title, relpath, tag))
		}
	}

	fmt.Println(strings.Join(poolTipsToShow, ""))

	doc := app.gmw.Parse(tb)
	targetHeader := app.gmw.GetHeadingNode(doc, tb, HEADING_NAME_TITLE, 1)
	tb = app.gmw.InsertTextAfter(doc, targetHeader, strings.Join(poolTipsToShow, ""), tb)

	return tb
}

func (app *App) WeeklyReport() {
	e, err := os.ReadDir(app.config.DailymemoDir)
	if err != nil {
		log.Fatal(err)
	}

	wantfiles := []string{}
	for _, file := range e {
		format := `\d{4}-\d{2}-\d{2}-\S{3}\.md`
		reg := regexp.MustCompile(format)
		if reg.MatchString(file.Name()) {
			wantfiles = append(wantfiles, filepath.Join(app.config.DailymemoDir, file.Name()))
		}
	}

	f, err := os.Create(filepath.Join(app.config.WeeklyReportFile))
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

		doc := app.gmw.Parse(b)
		hangingNodes := app.gmw.FindHeadingAndGetHangingNodes(doc, b, HEADING_NAME_MEMOS, 2)

		var order = 0
		for _, node := range hangingNodes {
			if n, ok := node.(*ast.Heading); ok {
				relpath, err := filepath.Rel(app.config.DailymemoDir, fpath)
				if err != nil {
					log.Fatal(err)
				}

				var format = "%d. [%s](%s#%s)\n"
				title := strings.Repeat("#", n.Level-2) + " " + string(node.Text(b))
				tag := text2tag(string(node.Text(b)))
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

func text2tag(text string) string {
	var tag = text
	tag = strings.ReplaceAll(tag, " ", "-")
	tag = strings.ReplaceAll(tag, "#", "")
	fullwidthchars := strings.Split("　！＠＃＄％＾＆＊（）＋｜〜＝￥｀「」｛｝；’：”、。・＜＞？【】『』《》〔〕［］‹›«»〘〙〚〛", "")
	for _, c := range fullwidthchars {
		tag = strings.ReplaceAll(tag, c, "")
	}
	return tag
}
