package models

import (
	"bytes"
	"strings"

	"github.com/hirotoni/memo/markdown"
)

type TipNodeKind int

const (
	TIPNODEKIND_DIR   TipNodeKind = iota // 0
	TIPNODEKIND_TITLE                    // 1
	TIPNODEKIND_TIP                      // 2
)

type TipNode struct {
	Kind  TipNodeKind
	Tip   Tip
	Text  string
	Depth int
}

func (tn *TipNode) Print(b *bytes.Buffer) {
	var out string
	switch tn.Kind {
	case TIPNODEKIND_DIR:
		out = strings.Repeat("  ", tn.Depth) + markdown.BuildList(tn.Text)
	case TIPNODEKIND_TITLE:
		out = strings.Repeat("  ", tn.Depth) + markdown.BuildList(tn.Text)
	case TIPNODEKIND_TIP:
		out = strings.Repeat("  ", tn.Depth) + markdown.BuildCheckbox(markdown.BuildLink(tn.Tip.Text, tn.Tip.Destination), tn.Tip.Checked)
	}
	b.WriteString(out + "\n")
}
