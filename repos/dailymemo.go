package repos

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hirotoni/memo/configs"
	"github.com/hirotoni/memo/models"
)

type DailymemoRepo struct {
	config *configs.TomlConfig
}

func NewDailymemoRepo(config *configs.TomlConfig) *DailymemoRepo {
	return &DailymemoRepo{
		config: config,
	}
}

var FILENAME_REGEX = `\d{4}-\d{2}-\d{2}-\S{3}\.md`

func (repo *DailymemoRepo) Entry(fpath string) models.Dailymemo {
	basename := filepath.Base(fpath)
	datestring, found := strings.CutSuffix(basename, ".md")
	if !found {
		log.Fatal("failed to cut suffix.")
	}
	date, err := time.Parse(FULL_LAYOUT, datestring)
	if err != nil {
		log.Fatal(err)
	}
	b, err := os.ReadFile(fpath)
	if err != nil {
		log.Fatal(err)
	}
	return models.Dailymemo{
		Filepath: fpath,
		BaseName: basename,
		Date:     date,
		Content:  b,
	}
}

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
	for _, fpath := range wantfiles {
		dm := repo.Entry(fpath)
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
	return repo.Entry(filepath)
}
