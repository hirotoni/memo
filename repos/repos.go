package repos

import (
	"github.com/hirotoni/memo/config"
)

type Repos struct {
	DailymemoRepo *DailymemoRepo
	TipRepo       *TipRepo
	TipNodeRepo   *TipNodeRepo
}

func NewRepos(config *config.TomlConfig) *Repos {
	return &Repos{
		DailymemoRepo: NewDailymemoRepo(config),
		TipRepo:       NewTipRepo(config),
		TipNodeRepo:   NewTipNodeRepo(config),
	}
}
