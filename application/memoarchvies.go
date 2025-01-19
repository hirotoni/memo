package application

import (
	"bytes"
	"log"
	"math/rand"
	"os"

	"github.com/hirotoni/memo/components"
	"github.com/hirotoni/memo/models"
)

// SaveMemoArchives generates memo archives index file
func (app *App) SaveMemoArchives() {
	app.saveMemoArchives(false)
}

func (app *App) saveMemoArchives(pickMemoArchive bool) *models.MemoArchive {
	checkedMemoArchives := app.repos.MemoArchiveRepo.MemoArchivesFromIndexChecked()                           // TODO handle error
	allMemoArchives := app.repos.MemoArchiveNodeRepo.MemoArchiveNodesFromMemoArchivesDir(checkedMemoArchives) // TODO handle error
	if len(allMemoArchives) == 0 {
		return nil
	}

	var picked *models.MemoArchive
	if pickMemoArchive {
		picked = pickRandomMemoArchive(allMemoArchives)
	}

	var buf = &bytes.Buffer{}
	components.PrintMemoArchiveNodeHeadingStyle(buf, allMemoArchives)

	// write memo archives to index
	f, err := os.Create(app.Config.MemoArchivesIndexFile())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	masb := []byte(components.GenerateTemplateString(components.TemplateMemoArchivesIndex))
	masb = app.gmw.InsertTextAfter(masb, components.HEADING_NAME_MEMOARCHIVES_INDEX, buf.String())
	f.Write(masb)

	return picked
}

func pickRandomMemoArchive(allMemoArchives []*models.MemoArchiveNode) *models.MemoArchive {
	notShown := filter(allMemoArchives, func(tn *models.MemoArchiveNode) bool {
		return tn.Kind == models.MEMOARCHIVENODEKIND_MEMO && !tn.MemoArchive.Checked
	})

	if len(notShown) == 0 {
		// reset all memo archives
		for _, v := range allMemoArchives {
			v.MemoArchive.Checked = false
		}
		notShown = allMemoArchives
	}

	picked, _ := randomPick(notShown)

	for _, v := range allMemoArchives {
		if v.MemoArchive.Destination == picked.MemoArchive.Destination {
			v.MemoArchive.Checked = true
		}
	}
	return &picked.MemoArchive
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
