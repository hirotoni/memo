package main

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
	kind  Kind
	tip   Tip
	text  string
	depth int
}

func (tn *TipNode) Print(b *bytes.Buffer) {
	var out string
	switch tn.kind {
	case KIND_DIR:
		out = strings.Repeat("  ", tn.depth) + markdown.BuildList(tn.text)
	case KIND_TITLE:
		out = strings.Repeat("  ", tn.depth) + markdown.BuildList(tn.text)
	case KIND_TIP:
		out = strings.Repeat("  ", tn.depth) + markdown.BuildCheckbox(markdown.BuildLink(tn.tip.Text, tn.tip.Destination), tn.tip.Checked)
	}
	b.WriteString(out + "\n")
}
