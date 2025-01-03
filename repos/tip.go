package repos

import (
	"log"
	"os"

	"github.com/hirotoni/memo/config"
	"github.com/hirotoni/memo/models"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

type TipRepo struct {
	config *config.TomlConfig
}

func NewTipRepo(config *config.TomlConfig) *TipRepo {
	return &TipRepo{
		config: config,
	}
}

func (repo *TipRepo) TipsFromIndex() []*models.Tip {
	b, err := os.ReadFile(repo.config.TipsIndexFile())
	if err != nil {
		log.Fatal(err)
	}

	doc := repo.config.Gmw.Parse(b)
	// doc.Dump(b, 1)

	var tips []*models.Tip
	var mywalker = func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if n.Kind() == ast.KindTextBlock && n.Parent().Kind() == ast.KindListItem {
				var t = &models.Tip{}
				for c := n.FirstChild(); c != nil; c = c.NextSibling() {
					if c, ok := c.(*ast.Link); ok {
						t.Text = string(c.Text(b))
						t.Destination = string(c.Destination)
					}
					if c, ok := c.(*extast.TaskCheckBox); ok {
						t.Checked = c.IsChecked
					}
				}
				if t.Text == "" && t.Destination == "" {
					return ast.WalkContinue, nil
				}
				tips = append(tips, t)
			}
		}
		return ast.WalkContinue, nil
	}
	err = ast.Walk(doc, mywalker)
	if err != nil {
		log.Fatal(err)
	}

	return tips
}

func (repo *TipRepo) TipsFromIndexChecked() []*models.Tip {
	tips := repo.TipsFromIndex()
	return filter(tips, func(t *models.Tip) bool { return t.Checked })
}

func filter[T any](ts []T, test func(T) bool) (ret []T) {
	for _, s := range ts {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
