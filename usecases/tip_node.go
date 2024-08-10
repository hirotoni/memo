package usecases

import (
	"bytes"
	"strings"

	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
)

func PrintTipNode(b *bytes.Buffer, tn *models.TipNode) {
	var out string
	switch tn.Kind {
	case models.TIPNODEKIND_DIR:
		out = strings.Repeat("  ", tn.Depth) + markdown.BuildList(tn.Text)
	case models.TIPNODEKIND_TITLE:
		out = strings.Repeat("  ", tn.Depth) + markdown.BuildList(tn.Text)
	case models.TIPNODEKIND_TIP:
		out = strings.Repeat("  ", tn.Depth) + markdown.BuildCheckbox(markdown.BuildLink(tn.Tip.Text, tn.Tip.Destination), tn.Tip.Checked)
	}
	b.WriteString(out + "\n")
}
