package repos

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hirotoni/memo/components"
	"github.com/hirotoni/memo/configs"
	"github.com/hirotoni/memo/markdown"

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

func (repo *DailymemoRepo) Entry(fpath string) (*models.Dailymemo, error) {
	basename := filepath.Base(fpath)
	datestring, found := strings.CutSuffix(basename, ".md")
	if !found {
		log.Fatal("failed to cut suffix.")
	}
	date, err := time.Parse(FULL_LAYOUT, datestring)
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	dm := &models.Dailymemo{
		Filepath: fpath,
		BaseName: basename,
		Date:     &date,
		Content:  b,
	}
	return dm, nil
}

func (repo *DailymemoRepo) Entries() ([]*models.Dailymemo, error) {
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

	dms := make([]*models.Dailymemo, 0, len(wantfiles))
	for _, fpath := range wantfiles {
		dm, err := repo.Entry(fpath)
		if err != nil {
			return nil, err
		}
		dms = append(dms, dm)
	}

	return dms, nil
}

var FULL_LAYOUT = "2006-01-02-Mon"

func (repo *DailymemoRepo) FindByDate(date string) (*models.Dailymemo, error) {
	_, err := time.Parse(FULL_LAYOUT, date)
	if err != nil {
		return nil, err
	}
	filepath := filepath.Join(repo.config.DailymemoDir(), date+".md")
	dm, err := repo.Entry(filepath)
	if err != nil {
		return nil, err
	}
	return dm, nil
}

func (repo *DailymemoRepo) Template() (*models.Dailymemo, error) {
	b, err := os.ReadFile(repo.config.DailymemoTemplateFile())
	if err != nil {
		return nil, err
	}
	dm := &models.Dailymemo{
		Filepath: repo.config.DailymemoTemplateFile(),
		BaseName: filepath.Base(repo.config.DailymemoTemplateFile()),
		Date:     nil,
		Content:  b,
	}
	return dm, nil
}

func (repo *DailymemoRepo) MemosFromDailymemo(dm *models.Dailymemo) []*models.Memo {
	// find memo block
	gmw := repo.config.Gmw
	_, b := gmw.FindHeadingAndGetHangingNodes(dm.Content, components.HEADING_NAME_MEMOS)
	sb := new(strings.Builder)
	gmw.RenderSlice(sb, dm.Content, b)

	memoBlock := []byte(sb.String())

	// extract titles
	_, memoHeadings := gmw.GetHeadingNodesByLevel(memoBlock, 3)
	var titles []string
	for _, heading := range memoHeadings {
		titles = append(titles, string(heading.Text(memoBlock)))
	}

	// extract each memo block
	var memos []*models.Memo
	for _, title := range titles {
		hhh := markdown.NewHeading(components.HEADING_NAME_MEMOS.Level+1, title)
		_, b := gmw.FindHeadingAndGetHangingNodes(memoBlock, hhh)
		sb := new(strings.Builder)
		gmw.RenderSlice(sb, memoBlock, b)
		mm := models.NewMemo(dm.Filepath, title, sb.String())

		memos = append(memos, mm)
	}
	return memos
}
