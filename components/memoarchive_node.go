package components

import (
	"bytes"
	"strings"

	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
)

func PrintMemoArchiveNode(b *bytes.Buffer, tn *models.MemoArchiveNode) {
	var out string
	switch tn.Kind {
	case models.MEMOARCHIVENODEKIND_DIR:
		out = strings.Repeat("  ", tn.Depth) + markdown.BuildList(tn.Text)
	case models.MEMOARCHIVENODEKIND_TITLE:
		out = strings.Repeat("  ", tn.Depth) + markdown.BuildList(tn.Text)
	case models.MEMOARCHIVENODEKIND_MEMO:
		out = strings.Repeat("  ", tn.Depth) + markdown.BuildCheckbox(markdown.BuildLink(tn.MemoArchive.Text, tn.MemoArchive.Destination), tn.MemoArchive.Checked)
	}
	b.WriteString(out + "\n")
}

func PrintMemoArchiveNodeHeadingStyle(b *bytes.Buffer, tns []*models.MemoArchiveNode) {
	var out string
	for i, tn := range tns {
		switch tn.Kind {
		case models.MEMOARCHIVENODEKIND_DIR:
			out = strings.Repeat("#", tn.Depth+2) + " " + tn.Text + "\n"
		case models.MEMOARCHIVENODEKIND_TITLE:
			out = strings.Repeat("#", tn.Depth+2) + " " + tn.Text + "\n"
		case models.MEMOARCHIVENODEKIND_MEMO:
			out = markdown.BuildCheckbox(markdown.BuildLink(tn.MemoArchive.Text, tn.MemoArchive.Destination), tn.MemoArchive.Checked)
			if i < len(tns)-1 && tns[i+1].Kind != models.MEMOARCHIVENODEKIND_MEMO {
				out += "\n"
			}
		}
		b.WriteString(out + "\n")
	}
}
