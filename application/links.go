package application

import (
	"fmt"
	"log"
	"strings"

	"github.com/hirotoni/memo/models"
)

func (app *App) Links() {
	// retrieve keys
	var memos []*models.Memo
	dms, err := app.repos.DailymemoRepo.Entries()
	if err != nil {
		log.Fatal("error")
	}
	for _, dm := range dms {
		mm := app.repos.DailymemoRepo.MemosFromDailymemo(dm)
		memos = append(memos, mm...)
	}

	// search links
	var links map[string][]string = make(map[string][]string)
	for _, m := range memos {
		for _, v := range memos {
			if strings.Contains(m.Content, v.SearchKey()) {
				key := v.Link()
				value := m.Link()
				links[key] = append(links[key], value)

				fmt.Println(m.Link(), " -> ", v.Link())

				// b, err := os.ReadFile(filepath.Join(app.Config.BaseDir, v.Filepath))
				// if err != nil {
				// 	log.Fatal(err)
				// }

				// h := markdown.Heading{
				// 	Level: 3,
				// 	Text:  v.Title,
				// }
				// inserted := app.Config.Gmw.InsertTextAfter(b, h, fmt.Sprintf("linked mentions\n\n- [%s](%s)", m.Link(), m.Link()))
				// os.WriteFile(filepath.Join(app.Config.BaseDir, v.Filepath), inserted, 0644)
			}
		}
	}
}
