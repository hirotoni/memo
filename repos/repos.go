package repos

import (
	"github.com/hirotoni/memo/configs"
)

type Repos struct {
	DailymemoRepo       *DailymemoRepo
	MemoArchiveRepo     *MemoArchiveRepo
	MemoArchiveNodeRepo *MemoArchiveNodeRepo
}

func NewRepos(config *configs.TomlConfig) *Repos {
	return &Repos{
		DailymemoRepo:       NewDailymemoRepo(config),
		MemoArchiveRepo:     NewMemoArchiveRepo(config),
		MemoArchiveNodeRepo: NewMemoArchiveNodeRepo(config),
	}
}
