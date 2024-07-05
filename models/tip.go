package models

import (
	"bytes"
	"strings"

	"github.com/hirotoni/memo/markdown"
)

type Tip struct {
	Text        string
	Destination string
	Checked     bool
}

type Kind int

const (
	KIND_DIR   Kind = iota // 0
	KIND_TITLE             // 1
	KIND_TIP               // 2
)

type TipNode struct {
	Kind  Kind
	Tip   Tip
	Text  string
	Depth int
}

func (tn *TipNode) Print(b *bytes.Buffer) {
	var out string
	switch tn.Kind {
	case KIND_DIR:
		out = strings.Repeat("  ", tn.Depth) + markdown.BuildList(tn.Text)
	case KIND_TITLE:
		out = strings.Repeat("  ", tn.Depth) + markdown.BuildList(tn.Text)
	case KIND_TIP:
		out = strings.Repeat("  ", tn.Depth) + markdown.BuildCheckbox(markdown.BuildLink(tn.Tip.Text, tn.Tip.Destination), tn.Tip.Checked)
	}
	b.WriteString(out + "\n")
}
