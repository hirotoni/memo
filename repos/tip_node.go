package repos

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hirotoni/memo/config"
	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
	"github.com/yuin/goldmark/ast"
)

type MemoArchiveNodeRepo struct {
	config *config.TomlConfig
}

func NewMemoArchiveNodeRepo(config *config.TomlConfig) *MemoArchiveNodeRepo {
	return &MemoArchiveNodeRepo{
		config: config,
	}
}

func (repo *MemoArchiveNodeRepo) MemoArchiveNodesFromMemoArchivesDir(shown []*models.MemoArchive) []*models.MemoArchiveNode {
	var tns []*models.MemoArchiveNode

	err := filepath.WalkDir(repo.config.MemoArchivesDir(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == repo.config.MemoArchivesTemplateFile() || path == repo.config.MemoArchivesIndexFile() {
			return nil
		}

		relpath, err := filepath.Rel(repo.config.DailymemoDir(), path)
		if err != nil {
			log.Fatal(err)
		}
		pathlist := strings.Split(relpath, string(filepath.Separator))
		depth := len(pathlist) - 3 // TODO avoid using magic number

		if d.IsDir() {
			if path == repo.config.MemoArchivesDir() {
				return nil
			}

			tmp := models.MemoArchiveNode{
				Kind:  models.MEMOARCHIVENODEKIND_DIR,
				Text:  d.Name(),
				Depth: depth,
			}
			tns = append(tns, &tmp)

		} else {
			if filepath.Ext(d.Name()) == ".md" {
				b, err := os.ReadFile(path)
				if err != nil {
					log.Fatal(err)
				}

				h1, h2s := repo.getMemoArchivesHeadings(b)
				if h1 == nil || h2s == nil {
					return nil
				}

				tmp := models.MemoArchiveNode{
					Kind:  models.MEMOARCHIVENODEKIND_TITLE,
					Text:  string(h1.Text(b)),
					Depth: depth,
					MemoArchive: models.MemoArchive{
						Text:        string(h1.Text(b)),
						Destination: relpath,
					},
				}
				tns = append(tns, &tmp)

				for _, h2 := range h2s {
					destination := relpath + "#" + markdown.Text2tag(string(h2.Text(b)))
					checked := slices.ContainsFunc(shown, func(t *models.MemoArchive) bool {
						return t.Destination == destination
					})

					tmp := models.MemoArchiveNode{
						Kind:  models.MEMOARCHIVENODEKIND_MEMO,
						Text:  string(h2.Text(b)),
						Depth: depth + 1,
						MemoArchive: models.MemoArchive{
							Text:        string(h2.Text(b)),
							Destination: destination,
							Checked:     checked,
						},
					}
					tns = append(tns, &tmp)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return tns
}
func (repo *MemoArchiveNodeRepo) getMemoArchivesHeadings(b []byte) (ast.Node, []ast.Node) {
	_, headings := repo.config.Gmw.GetHeadingNodesByLevel(b, 1)
	if len(headings) == 0 {
		return nil, nil
	}
	heading := headings[0]
	heading1, nodes := repo.config.Gmw.FindHeadingAndGetHangingNodes(b, models.Heading{Level: 1, Text: string(heading.Text(b))})

	heading2s := filter(nodes, func(n ast.Node) bool {
		h, ok := n.(*ast.Heading)
		return ok && h.Level == 2
	})

	return heading1, heading2s

}
