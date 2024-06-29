package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
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
		b = app.inheritHeading(b, HEADING_NAME_TODOS)
		b = app.inheritHeading(b, HEADING_NAME_WANTTODOS)
		b = app.appendTips(b)

		b = app.gmw.InsertTextAfter(b, HEADING_NAME_TITLE, today)

		f.Write(b)
	}

	// open memo dir with editor
	cmd := exec.Command("code", targetFile, "--folder-uri", app.config.BaseDir())
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// WeeklyReport generates weekly report file
func (app *App) WeeklyReport(openEditor bool) {
	entries, err := os.ReadDir(app.config.DailymemoDir())
	if err != nil {
		log.Fatal(err)
	}

	wantfiles := []string{}
	reg := regexp.MustCompile(LAYOUT_REGEX)
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

		_, hangingNodes := app.gmw.FindHeadingAndGetHangingNodes(b, HEADING_NAME_MEMOS)

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

	f.WriteString(TemplateWeeklyReport.String() + "\n")
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

func (app *App) saveTips(pickTip bool) Tip {
	var picked Tip

	checkedTips := app.getTipsCheckedFromIndex()
	allTips := app.getTipNodesFromDir(checkedTips)

	if pickTip {
		notShown := filter(allTips, func(tn TipNode) bool { return tn.kind == KIND_TIP && !tn.tip.Checked })
		p, _ := randomPick(notShown)
		picked = p.tip

		for i, v := range allTips {
			if v.tip.Destination == picked.Destination {
				allTips[i].tip.Checked = true
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
	tipsb := []byte(TemplateTipsIndex.String())
	tipsb = app.gmw.InsertTextAfter(tipsb, HEADING_NAME_TIPSINDEX, buf.String())
	f.Write(tipsb)

	return picked
}

// inheritHeading inherits todos from previous day's memo
func (app *App) inheritHeading(tb []byte, heading md.Heading) []byte {
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
	chosenTip := md.BuildList(md.BuildLink(
		picked.Text,
		picked.Destination,
	))
	tb = app.gmw.InsertTextAfter(tb, HEADING_NAME_TITLE, chosenTip)

	return tb
}

// getTipsFromIndex reads tips from index file
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
				if t.Text == "" && t.Destination == "" {
					return ast.WalkContinue, nil
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

func (app *App) getTipsCheckedFromIndex() []Tip {
	tips := app.getTipsFromIndex()
	return filter(tips, func(t Tip) bool { return t.Checked })
}

func (app *App) getTipNodesFromDir(shown []Tip) []TipNode {
	var tns []TipNode

	err := filepath.WalkDir(app.config.TipsDir(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == app.config.TipsTemplateFile() || path == app.config.TipsIndexFile() {
			return nil
		}

		relpath, err := filepath.Rel(app.config.DailymemoDir(), path)
		if err != nil {
			log.Fatal(err)
		}
		pathlist := strings.Split(relpath, string(filepath.Separator))
		depth := len(pathlist) - 3 // TODO avoid using magic number

		if d.IsDir() {
			if path == app.config.TipsDir() {
				return nil
			}

			tmp := TipNode{
				kind:  KIND_DIR,
				text:  d.Name(),
				depth: depth,
			}
			tns = append(tns, tmp)

		} else {
			if filepath.Ext(d.Name()) == ".md" {
				b, err := os.ReadFile(path)
				if err != nil {
					log.Fatal(err)
				}

				h1, h2s := app.getTipsHeadings(b)
				if h1 == nil || h2s == nil {
					return nil
				}

				tmp := TipNode{
					kind:  KIND_TITLE,
					text:  string(h1.Text(b)),
					depth: depth,
					tip: Tip{
						Text:        string(h1.Text(b)),
						Destination: relpath,
					},
				}
				tns = append(tns, tmp)

				for _, h2 := range h2s {
					destination := relpath + "#" + md.Text2tag(string(h2.Text(b)))
					checked := slices.ContainsFunc(shown, func(t Tip) bool {
						return t.Destination == destination
					})

					tmp := TipNode{
						kind:  KIND_TIP,
						text:  string(h2.Text(b)),
						depth: depth + 1,
						tip: Tip{
							Text:        string(h2.Text(b)),
							Destination: destination,
							Checked:     checked,
						},
					}
					tns = append(tns, tmp)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return tns
}

func (app *App) getTipsHeadings(b []byte) (ast.Node, []ast.Node) {
	_, headings := app.gmw.GetHeadingNodesByLevel(b, 1)
	if len(headings) == 0 {
		return nil, nil
	}
	heading := headings[0]
	heading1, nodes := app.gmw.FindHeadingAndGetHangingNodes(b, md.Heading{Level: 1, Text: string(heading.Text(b))})

	heading2s := filter(nodes, func(n ast.Node) bool {
		h, ok := n.(*ast.Heading)
		return ok && h.Level == 2
	})

	return heading1, heading2s

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
