package repos

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/hirotoni/memo/config"
	"github.com/hirotoni/memo/models"
)

type DailymemoRepo struct {
	config *config.TomlConfig
}

func NewDailymemoRepo(config *config.TomlConfig) *DailymemoRepo {
	return &DailymemoRepo{
		config: config,
	}
}

var FILENAME_REGEX = `\d{4}-\d{2}-\d{2}-\S{3}\.md`

func (repo *DailymemoRepo) Entries() []models.Dailymemo {
	entries, err := os.ReadDir(repo.config.DailymemoDir()) // sorted by filename(=date)
	if err != nil {
		log.Fatal(err)
	}

	wantfiles := make([]string, 0, len(entries))
	reg := regexp.MustCompile(FILENAME_REGEX)
	for _, file := range entries {
		if reg.MatchString(file.Name()) {
			wantfiles = append(wantfiles, filepath.Join(repo.config.DailymemoDir(), file.Name()))
		}
	}

	dms := make([]models.Dailymemo, 0, len(wantfiles))
	for _, file := range wantfiles {
		dm := models.NewDailymemoFromFilepath(file)
		dms = append(dms, dm)
	}

	return dms
}

var FULL_LAYOUT = "2006-01-02-Mon"

func (repo *DailymemoRepo) FindByDate(date string) models.Dailymemo {
	_, err := time.Parse(FULL_LAYOUT, date)
	if err != nil {
		log.Fatal(err)
	}
	filepath := filepath.Join(repo.config.DailymemoDir(), date+".md")
	return models.NewDailymemoFromFilepath(filepath)
}
