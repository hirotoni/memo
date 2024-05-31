package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	md "github.com/hirotoni/memo/markdown"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

type App struct {
	gmw    *md.GoldmarkWrapper
	config AppConfig
}

func NewApp() App {
	return App{
		gmw:    md.NewGoldmarkWrapper(),
		config: NewAppConfig(),
	}
}

// Initialize initializes dirs and files
func (app *App) Initialize() {
	// dailymemo dir
	_, err := os.Stat(app.config.DailymemoDir())
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(app.config.DailymemoDir(), 0750); err != nil {
			log.Fatal(err)
		}
		log.Printf("memo directory initialized: %s", app.config.BaseDir())
	}
	// dailymemo template file
	_, err = os.Stat(app.config.DailymemoTemplateFile())
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(app.config.DailymemoTemplateFile())
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.WriteString(TemplateDailymemo.String())

		log.Printf("dailymemo template file initialized: %s", app.config.DailymemoTemplateFile())
	}

	// tips dir
	_, err = os.Stat(app.config.TipsDir())
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(app.config.TipsDir(), 0750); err != nil {
			log.Fatal(err)
		}
		log.Printf("tips directory initialized: %s", app.config.TipsDir())
	}
	// tips template file
	_, err = os.Stat(app.config.TipsTemplateFile())
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(app.config.TipsTemplateFile())
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.WriteString(TemplateTips.String())

		log.Printf("tips template file initialized: %s", app.config.TipsTemplateFile())
	}
	// tips index file
	_, err = os.Stat(app.config.TipsIndexFile())
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(app.config.TipsIndexFile())
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.WriteString(TemplateTipsIndex.String())

		log.Printf("tips index file initialized: %s", app.config.TipsIndexFile())
	}
}

// OpenTodaysMemo opens today's memo
func (app *App) OpenTodaysMemo(truncate bool) {
	today := time.Now().Format(LAYOUT)
	targetFile := filepath.Join(app.config.DailymemoDir(), today+".md")

	log.Default().Printf("truncate: %v", truncate)

	_, err := os.Stat(targetFile)
	if errors.Is(err, os.ErrNotExist) || truncate {
		f, err := os.Create(targetFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		b, err := os.ReadFile(app.config.DailymemoTemplateFile())
		if err != nil {
			log.Fatal(err)
		}

		// inherit todos from previous memo
		b = app.InheritHeading(b, HEADING_NAME_TODOS)
		b = app.InheritHeading(b, HEADING_NAME_WANTTODOS)
		b = app.AppendTips(b)

		b = app.gmw.InsertTextAfter(b, HEADING_NAME_TITLE, today)

		f.Write(b)
	}

	// open memo dir with editor
	cmd := exec.Command("code", targetFile, "--folder-uri", app.config.BaseDir())
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// InheritHeading inherits todos from previous day's memo
func (app *App) InheritHeading(tb []byte, heading md.Heading) []byte {
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

		nodesToInsert := app.gmw.FindHeadingAndGetHangingNodes(pb, heading)
		tb = app.gmw.InsertNodesAfter(tb, heading, pb, nodesToInsert)
		break
	}

	return tb
}

// AppendTips appends tips
func (app *App) AppendTips(tb []byte) []byte {
	// tips are the things that you want to remember periodically such as
	// - ER diagrams, component diagrams, constants of application you are in charge
	// - product management, development process knowledge
	// - bookmarks, web links
	// - life sayings, someone's sayings

	var targetTipFiles []string
	var allTipsNotShown []Tip
	var allTipsShown []Tip
	var tipsToIndex []string

	var indexTipsShown = filter(app.getTipsFromIndex(), func(t Tip) bool { return t.Checked })

	var fn = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || path == app.config.TipsTemplateFile() {
			return nil
		}

		targetTipFiles = append(targetTipFiles, path)
		return nil
	}

	if err := filepath.Walk(app.config.TipsDir(), fn); err != nil {
		log.Fatal(err)
	}

	for _, v := range targetTipFiles {
		b, err := os.ReadFile(v)
		if err != nil {
			log.Fatal(err)
		}
		_, headings := app.gmw.GetHeadingNodes(b, 2)

		relpath, err := filepath.Rel(app.config.DailymemoDir(), v)
		if err != nil {
			log.Fatal(err)
		}
		for _, vv := range headings {
			title := string(vv.Text(b))
			destination := relpath + "#" + md.Text2tag(title)
			checked := slices.ContainsFunc(indexTipsShown, func(t Tip) bool { return t.Destination == destination })

			tip := Tip{
				Text:        title,
				Destination: destination,
				Checked:     checked,
			}

			if tip.Checked {
				allTipsShown = append(allTipsShown, tip)
			} else {
				allTipsNotShown = append(allTipsNotShown, tip)
			}
		}
	}

	// log.Default().Println(allTipsShown)
	// log.Default().Println(allTipsNotShown)
	// log.Default().Println(indexTipsShown)

	// if all tips have been shown, then reset
	if len(allTipsNotShown) == 0 {
		var tmp []Tip
		for _, v := range allTipsShown {
			v.Checked = false
			tmp = append(tmp, v)
		}
		allTipsNotShown = tmp
		allTipsShown = []Tip{}
	}

	// pick one
	chosen := rand.Intn(len(allTipsNotShown))

	// insert todays tip
	chosenTip := md.BuildList(md.BuildLink(allTipsNotShown[chosen].Text, allTipsNotShown[chosen].Destination))
	tb = app.gmw.InsertTextAfter(tb, HEADING_NAME_TITLE, chosenTip)

	// groom tips to index
	var groom []Tip
	for i, v := range allTipsNotShown {
		if i == chosen {
			v.Checked = true
		}
		groom = append(groom, v)
	}

	for _, v := range allTipsShown {
		v.Checked = true
		groom = append(groom, v)
	}

	slices.SortFunc(groom, func(a, b Tip) int {
		if a.Destination < b.Destination {
			return -1
		} else {
			return 1
		}
	})

	for _, v := range groom {
		index := md.BuildCheckbox(md.BuildLink(v.Text, v.Destination), v.Checked) + "\n"
		tipsToIndex = append(tipsToIndex, index)
	}

	// write tips to index
	f, err := os.Create(app.config.TipsIndexFile())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	tipsb := []byte(TemplateTipsIndex.String())
	tipsb = app.gmw.InsertTextAfter(tipsb, HEADING_NAME_TIPSINDEX, strings.Join(tipsToIndex, ""))

	f.Write(tipsb)

	return tb
}

func (app *App) WeeklyReport() {
	e, err := os.ReadDir(app.config.DailymemoDir())
	if err != nil {
		log.Fatal(err)
	}

	wantfiles := []string{}
	reg := regexp.MustCompile(LAYOUT_REGEX)
	for _, file := range e {
		if reg.MatchString(file.Name()) {
			wantfiles = append(wantfiles, filepath.Join(app.config.DailymemoDir(), file.Name()))
		}
	}

	f, err := os.Create(filepath.Join(app.config.WeeklyReportFile()))
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

		hangingNodes := app.gmw.FindHeadingAndGetHangingNodes(b, HEADING_NAME_MEMOS)

		var order = 0
		for _, node := range hangingNodes {
			if n, ok := node.(*ast.Heading); ok {
				relpath, err := filepath.Rel(app.config.DailymemoDir(), fpath)
				if err != nil {
					log.Fatal(err)
				}

				order++
				title := strings.Repeat("#", n.Level-2) + " " + string(node.Text(b))
				tag := md.Text2tag(string(node.Text(b)))
				s := md.BuildOrderedList(order, md.BuildLink(title, relpath+"#"+tag)) + "\n"
				f.WriteString(s)
			}
		}

		if order > 0 {
			f.WriteString("\n")
		}
	}
}

type Tip struct {
	Text        string
	Destination string
	Checked     bool
}

func (app *App) getTipsFromIndex() []Tip {
	b, err := os.ReadFile(app.config.TipsIndexFile())
	if err != nil {
		log.Fatal(err)
	}

	doc := app.gmw.Parse(b)
	// doc.Dump(b, 1)

	var tips []Tip
	var mywalker = func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if n.Kind() == ast.KindTextBlock && n.Parent().Kind() == ast.KindListItem {
				var t = Tip{}
				for c := n.FirstChild(); c != nil; c = c.NextSibling() {
					if c, ok := c.(*ast.Link); ok {
						t.Text = string(c.Text(b))
						t.Destination = string(c.Destination)
					}
					if c, ok := c.(*extast.TaskCheckBox); ok {
						t.Checked = c.IsChecked
					}
				}
				tips = append(tips, t)
			}
		}
		return ast.WalkContinue, nil
	}
	err = ast.Walk(doc, mywalker)
	if err != nil {
		log.Fatal(err)
	}

	return tips
}

func filter[T any](ts []T, test func(T) bool) (ret []T) {
	for _, s := range ts {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
