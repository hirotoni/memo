package renderer

import (
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

func generateHeader(level int, isBlankSpacePrevious bool) ast.Node {
	h := ast.NewHeading(level)
	h.SetBlankPreviousLines(isBlankSpacePrevious)
	return h
}

func generateLink(text, destination []byte) ast.Node {
	nl := ast.NewLink()
	nl.Destination = destination

	// segment
	t := ast.NewText()
	t.Segment.Start = 0
	t.Segment.Stop = len(text)
	nl.AppendChild(nl, t)

	return nl
}

func generateAutoLink(text []byte) ast.Node {
	// segment
	t := ast.NewText()
	t.Segment.Start = 0
	t.Segment.Stop = len(text)

	al := ast.NewAutoLink(ast.AutoLinkURL, t)

	return al
}

func generateTaskCheckBox(checked bool) ast.Node {
	return extast.NewTaskCheckBox(checked)
}

func generateEmphasis(level int) ast.Node {
	return ast.NewEmphasis(level)
}
