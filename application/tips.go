package application

import (
	"bytes"
	"log"
	"math/rand"
	"os"

	"github.com/hirotoni/memo/models"
	"github.com/hirotoni/memo/usecases"
)

// SaveTips generates tips index file
func (app *App) SaveTips() {
	app.saveTips(false)
}

func (app *App) saveTips(pickTip bool) *models.Tip {
	checkedTips := app.repos.TipRepo.TipsFromIndexChecked()           // TODO handle error
	allTips := app.repos.TipNodeRepo.TipNodesFromTipsDir(checkedTips) // TODO handle error
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
	usecases.PrintTipNodeHeadingStyle(buf, allTips)

	// write tips to index
	f, err := os.Create(app.Config.TipsIndexFile())
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
