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

type TipNodeRepo struct {
	config *config.TomlConfig
	gmw    *markdown.GoldmarkWrapper
}

func NewTipNodeRepo(config *config.TomlConfig, gmw *markdown.GoldmarkWrapper) *TipNodeRepo {
	return &TipNodeRepo{
		config: config,
		gmw:    gmw,
	}
}

func (repo *TipNodeRepo) TipNodesFromTipsDir(shown []*models.Tip) []*models.TipNode {
	var tns []*models.TipNode

	err := filepath.WalkDir(repo.config.TipsDir(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == repo.config.TipsTemplateFile() || path == repo.config.TipsIndexFile() {
			return nil
		}

		relpath, err := filepath.Rel(repo.config.DailymemoDir(), path)
		if err != nil {
			log.Fatal(err)
		}
		pathlist := strings.Split(relpath, string(filepath.Separator))
		depth := len(pathlist) - 3 // TODO avoid using magic number

		if d.IsDir() {
			if path == repo.config.TipsDir() {
				return nil
			}

			tmp := models.TipNode{
				Kind:  models.TIPNODEKIND_DIR,
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

				h1, h2s := repo.getTipsHeadings(b)
				if h1 == nil || h2s == nil {
					return nil
				}

				tmp := models.TipNode{
					Kind:  models.TIPNODEKIND_TITLE,
					Text:  string(h1.Text(b)),
					Depth: depth,
					Tip: models.Tip{
						Text:        string(h1.Text(b)),
						Destination: relpath,
					},
				}
				tns = append(tns, &tmp)

				for _, h2 := range h2s {
					destination := relpath + "#" + markdown.Text2tag(string(h2.Text(b)))
					checked := slices.ContainsFunc(shown, func(t *models.Tip) bool {
						return t.Destination == destination
					})

					tmp := models.TipNode{
						Kind:  models.TIPNODEKIND_TIP,
						Text:  string(h2.Text(b)),
						Depth: depth + 1,
						Tip: models.Tip{
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
func (repo *TipNodeRepo) getTipsHeadings(b []byte) (ast.Node, []ast.Node) {
	_, headings := repo.gmw.GetHeadingNodesByLevel(b, 1)
	if len(headings) == 0 {
		return nil, nil
	}
	heading := headings[0]
	heading1, nodes := repo.gmw.FindHeadingAndGetHangingNodes(b, models.Heading{Level: 1, Text: string(heading.Text(b))})

	heading2s := filter(nodes, func(n ast.Node) bool {
		h, ok := n.(*ast.Heading)
		return ok && h.Level == 2
	})

	return heading1, heading2s

}
