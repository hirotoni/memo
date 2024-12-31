package repos

import (
	"github.com/hirotoni/memo/config"
	"github.com/hirotoni/memo/markdown"
)

type Repos struct {
	DailymemoRepo *DailymemoRepo
	TipRepo       *TipRepo
	TipNodeRepo   *TipNodeRepo
}

func NewRepos(config *config.TomlConfig, gmw *markdown.GoldmarkWrapper) *Repos {
	return &Repos{
		DailymemoRepo: NewDailymemoRepo(config, gmw),
		TipRepo:       NewTipRepo(config, gmw),
		TipNodeRepo:   NewTipNodeRepo(config, gmw),
	}
}
