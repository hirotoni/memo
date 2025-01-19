package models

type MemoArchiveNodeKind int

const (
	MEMOARCHIVENODEKIND_DIR   MemoArchiveNodeKind = iota // 0
	MEMOARCHIVENODEKIND_TITLE                            // 1
	MEMOARCHIVENODEKIND_MEMO                             // 2
)

type MemoArchiveNode struct {
	Kind        MemoArchiveNodeKind
	MemoArchive MemoArchive
	Text        string
	Depth       int
}
