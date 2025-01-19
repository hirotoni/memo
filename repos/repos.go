package repos

import (
	"github.com/hirotoni/memo/config"
)

type Repos struct {
	DailymemoRepo       *DailymemoRepo
	MemoArchiveRepo     *MemoArchiveRepo
	MemoArchiveNodeRepo *MemoArchiveNodeRepo
}

func NewRepos(config *config.TomlConfig) *Repos {
	return &Repos{
		DailymemoRepo:       NewDailymemoRepo(config),
		MemoArchiveRepo:     NewMemoArchiveRepo(config),
		MemoArchiveNodeRepo: NewMemoArchiveNodeRepo(config),
	}
}
