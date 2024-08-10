package models

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
