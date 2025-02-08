package application

import (
	"fmt"
	"log"
	"slices"
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
			}
		}
	}

	var keys []string
	for k := range links {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	for _, k := range keys {
		for _, v := range links[k] {
			fmt.Println(v, " -> ", k)
		}
	}
}
