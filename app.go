package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hirotoni/memo/config"
	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
	"github.com/hirotoni/memo/repos"
	"github.com/hirotoni/memo/usecases"
	"github.com/yuin/goldmark/ast"
)

const (
	TIMEZONE        = "Asia/Tokyo"
	LAYOUT          = "2006-01-02-Mon"
	FILENAME_REGEX  = `\d{4}-\d{2}-\d{2}-\S{3}\.md`
	FILENAME_FORMAT = "%s.md"

	// number of dates to seek back when inheriting todos from previous days
	DAYS_TO_SEEK = 10
)

type App struct {
	gmw    *markdown.GoldmarkWrapper
	config *config.AppConfig
}

func NewApp() App {
	return App{
		gmw:    markdown.NewGoldmarkWrapper(),
		config: config.NewAppConfig(),
	}
}

// Initialize initializes dirs and files
func (app *App) Initialize() {
	// dailymemo
	initializeDir(app.config.DailymemoDir())
	initializeFile(app.config.DailymemoTemplateFile(), usecases.TemplateDailymemo)
	// tips
	initializeDir(app.config.TipsDir())
	initializeFile(app.config.TipsTemplateFile(), usecases.TemplateTips)
	initializeFile(app.config.TipsIndexFile(), usecases.TemplateTipsIndex)
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

// OpenTodaysMemo opens today's memo
func (app *App) OpenTodaysMemo(truncate bool) {
	today := time.Now().Format(LAYOUT)
	filename := fmt.Sprintf(FILENAME_FORMAT, today)
	targetFile := filepath.Join(app.config.DailymemoDir(), filename)

	log.Default().Printf("truncate: %v", truncate)

	_, err := os.Stat(targetFile)
	if errors.Is(err, os.ErrNotExist) || truncate {
		f, err := os.Create(targetFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.Write(app.generateTodaysMemo(today))
	}

	// open memo dir with editor
	cmd := exec.Command("code", targetFile, "--folder-uri", app.config.BaseDir())
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func (app *App) generateTodaysMemo(today string) []byte {
	b, err := os.ReadFile(app.config.DailymemoTemplateFile())
	if err != nil {
		log.Fatal(err)
	}

	// inherit todos and wanttodos from previous memo
	b = app.inheritHeading(b, usecases.HEADING_NAME_TODOS)
	b = app.inheritHeading(b, usecases.HEADING_NAME_WANTTODOS)
	b = app.appendTips(b)

	b = app.gmw.InsertTextAfter(b, usecases.HEADING_NAME_TITLE, today)

	return b
}

// inheritHeading inherits todos from previous day's memo
func (app *App) inheritHeading(tb []byte, heading markdown.Heading) []byte {
	// previous days
	today := time.Now()
	for i := range make([]int, DAYS_TO_SEEK) {
		previousDay := today.AddDate(0, 0, -1*(i+1)).Format(LAYOUT)
		pb, err := os.ReadFile(filepath.Join(app.config.DailymemoDir(), previousDay+".md"))
		if errors.Is(err, os.ErrNotExist) {
			if i+1 == DAYS_TO_SEEK {
				log.Printf("previous memos were not found in previous %d days.", DAYS_TO_SEEK)
			}
			continue
		} else if err != nil {
			log.Fatal(err)
		}

		_, nodesToInsert := app.gmw.FindHeadingAndGetHangingNodes(pb, heading)
		tb = app.gmw.InsertNodesAfter(tb, heading, pb, nodesToInsert)
		break
	}

	return tb
}

// appendTips appends tips
func (app *App) appendTips(tb []byte) []byte {
	// tips are the things that you want to remember periodically such as
	// - ER diagrams, component diagrams, constants of application you are in charge
	// - product management, development process knowledge
	// - bookmarks, web links
	// - life sayings, someone's sayings

	picked := app.saveTips(true)

	// insert todays tip
	if picked.Destination != "" {
		chosenTip := markdown.BuildList(markdown.BuildLink(
			picked.Text,
			picked.Destination,
		))
		tb = app.gmw.InsertTextAfter(tb, usecases.HEADING_NAME_TITLE, chosenTip)
	}

	return tb
}

// WeeklyReport generates weekly report file
func (app *App) WeeklyReport(openEditor bool) {
	entries, err := os.ReadDir(app.config.DailymemoDir())
	if err != nil {
		log.Fatal(err)
	}

	wantfiles := []string{}
	reg := regexp.MustCompile(FILENAME_REGEX)
	for _, file := range entries {
		if reg.MatchString(file.Name()) {
			wantfiles = append(wantfiles, filepath.Join(app.config.DailymemoDir(), file.Name()))
		}
	}

	sb := strings.Builder{}
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
			sb.WriteString("## " + fmt.Sprint(year) + " | Week " + fmt.Sprint(week) + "\n\n")
			curWeekNum = week
		}

		sb.WriteString("### " + filepath.Base(fpath) + "\n\n")

		b, err := os.ReadFile(fpath)
		if err != nil {
			log.Fatal(err)
		}

		_, hangingNodes := app.gmw.FindHeadingAndGetHangingNodes(b, usecases.HEADING_NAME_MEMOS)

		var order = 0
		for _, node := range hangingNodes {
			if n, ok := node.(*ast.Heading); ok {
				relpath, err := filepath.Rel(app.config.DailymemoDir(), fpath)
				if err != nil {
					log.Fatal(err)
				}

				order++
				title := strings.Repeat("#", n.Level-2) + " " + string(node.Text(b))
				tag := markdown.Text2tag(string(node.Text(b)))
				s := markdown.BuildOrderedList(order, markdown.BuildLink(title, relpath+"#"+tag)) + "\n"
				sb.WriteString(s)
			}
		}

		if order > 0 {
			sb.WriteString("\n")
		}
	}

	f, err := os.Create(app.config.WeeklyReportFile())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(usecases.GenerateTemplateString(usecases.TemplateWeeklyReport) + "\n")
	f.WriteString(sb.String())

	if openEditor {
		// open memo dir with editor
		cmd := exec.Command("code", app.config.WeeklyReportFile(), "--folder-uri", app.config.BaseDir())
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

// SaveTips generates tips index file
func (app *App) SaveTips(openEditor bool) {
	app.saveTips(false)

	if openEditor {
		// open memo dir with editor
		cmd := exec.Command("code", app.config.TipsIndexFile(), "--folder-uri", app.config.BaseDir())
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func (app *App) saveTips(pickTip bool) *models.Tip {
	tRepo := repos.NewTipRepo(app.config, app.gmw)
	tnRepo := repos.NewTipNodeRepo(app.config, app.gmw)

	checkedTips := tRepo.TipsFromIndexChecked()        // TODO handle error
	allTips := tnRepo.TipNodesFromTipsDir(checkedTips) // TODO handle error
	if len(allTips) == 0 {
		return nil
	}

	var picked *models.Tip
	if pickTip {
		notShown := filter(allTips, func(tn *models.TipNode) bool { return tn.Kind == models.TIPNODEKIND_TIP && !tn.Tip.Checked })

		if len(notShown) == 0 {
			// reset all tips
			for _, v := range allTips {
				v.Tip.Checked = false
			}
			notShown = allTips
		}

		p, _ := randomPick(notShown)
		picked = &p.Tip

		for _, v := range allTips {
			if v.Tip.Destination == picked.Destination {
				v.Tip.Checked = true
			}
		}
	}

	var buf = &bytes.Buffer{}
	for _, v := range allTips {
		v.Print(buf)
	}

	// write tips to index
	f, err := os.Create(app.config.TipsIndexFile())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	tipsb := []byte(usecases.GenerateTemplateString(usecases.TemplateTipsIndex))
	tipsb = app.gmw.InsertTextAfter(tipsb, usecases.HEADING_NAME_TIPSINDEX, buf.String())
	f.Write(tipsb)

	return picked
}

func filter[T any](ts []T, test func(T) bool) (ret []T) {
	for _, s := range ts {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func randomPick[T any](s []T) (T, []T) {
	i := rand.Intn(len(s))
	picked := s[i]
	return picked, append(s[:i], s[i+1:]...)
}
