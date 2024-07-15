package renderer

import (
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

func genHeaderNode(level int, isBlankSpacePrevious bool) ast.Node {
	h := ast.NewHeading(level)
	h.SetBlankPreviousLines(isBlankSpacePrevious)
	return h
}

func genLinkNode(text, destination []byte) ast.Node {
	nl := ast.NewLink()
	nl.Destination = destination

	// segment
	t := ast.NewText()
	t.Segment.Start = 0
	t.Segment.Stop = len(text)
	nl.AppendChild(nl, t)

	return nl
}

func genAutoLinkNode(text []byte) ast.Node {
	// segment
	t := ast.NewText()
	t.Segment.Start = 0
	t.Segment.Stop = len(text)

	al := ast.NewAutoLink(ast.AutoLinkURL, t)

	return al
}

func genTaskCheckBoxNode(checked bool) ast.Node {
	return extast.NewTaskCheckBox(checked)
}

func genEnphasisNode(level int) ast.Node {
	return ast.NewEmphasis(level)
}
