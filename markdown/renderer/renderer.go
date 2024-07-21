package renderer

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type MarkdownRenderer struct {
	MarkdownRendererConfig
}

func NewMarkdownRenderer() *MarkdownRenderer {
	return &MarkdownRenderer{
		MarkdownRendererConfig: *NewMarkdownRendererConfig(),
	}
}

func (r *MarkdownRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks
	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindHeading, r.renderHeading)
	// reg.Register(ast.KindBlockquote, r.renderBlockquote)
	// reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	// reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	// reg.Register(ast.KindHTMLBlock, r.renderHTMLBlock)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindTextBlock, r.renderTextBlock)
	// reg.Register(ast.KindThematicBreak, r.renderThematicBreak)

	// // inlines
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
	// reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	// reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindLink, r.renderLink)
	// reg.Register(ast.KindRawHTML, r.renderRawHTML)
	reg.Register(ast.KindText, r.renderText)
	// reg.Register(ast.KindString, r.renderString)
	reg.Register(extast.KindTaskCheckBox, r.renderTaskCheckBox)
}

// MARK: blocks

func (r *MarkdownRenderer) renderDocument(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// nothing to do
	return ast.WalkContinue, nil
}

func (r *MarkdownRenderer) renderHeading(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		if n.HasBlankPreviousLines() {
			_, _ = w.WriteString("\n\n")
		}
		_, _ = w.WriteString(strings.Repeat("#", n.Level) + " ")
	}
	return ast.WalkContinue, nil
}

func (r *MarkdownRenderer) renderParagraph(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Paragraph)
	if entering {
		if n.HasBlankPreviousLines() {
			_, _ = w.WriteString("\n\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *MarkdownRenderer) renderText(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Text)
	if entering {
		p := n.Parent()
		if p == nil {
			return ast.WalkContinue, nil
		}
		if p.Kind() == ast.KindLink {
			// r.renderLink() renders text in advance. no rendering needed here.
			return ast.WalkContinue, nil
		}

		w.WriteString(string(n.Text(source)))

		if n.SoftLineBreak() {
			w.WriteString("\n")

			pp := p.Parent()
			if pp, ok := pp.(*ast.ListItem); ok {
				// ListItem - TextBlock - Text(SoftLineBreak)
				w.WriteString(strings.Repeat(" ", pp.Offset))
			}
		}
	}
	return ast.WalkContinue, nil
}

func (r *MarkdownRenderer) renderList(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.List)
	if entering {
		if n.HasBlankPreviousLines() {
			_, _ = w.WriteString("\n\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *MarkdownRenderer) renderListItem(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.ListItem)

	// at least there must be a parent and a grandparent
	// e.g. Document - List - ListItem
	p := n.Parent()
	if p == nil {
		return ast.WalkContinue, nil
	}
	pp := p.Parent()
	if pp == nil {
		return ast.WalkContinue, nil
	}

	if entering {
		// If it is not the first element of the list or it is a nested listitems, add a line break
		if n.PreviousSibling() != nil || pp.Kind() == ast.KindListItem {
			w.WriteString("\n")
		}

		// Increase the indent for nested lists
		curpp := pp
		for curpp != nil {
			li, ok := curpp.(*ast.ListItem)
			if ok {
				w.WriteString(strings.Repeat(" ", li.Offset))
			}
			curpp = safeGroundParent(curpp)
		}

		if p, ok := p.(*ast.List); ok {
			if p.IsOrdered() {
				order := p.Start
				for node.PreviousSibling() != nil {
					order++
					node = node.PreviousSibling()
				}
				w.WriteString(fmt.Sprintf("%d", order))
			}

			w.WriteString(string(p.Marker) + " ")
		}
	}
	return ast.WalkContinue, nil
}

func (r *MarkdownRenderer) renderTextBlock(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// n := node.(*ast.TextBlock)
	// nothing to do
	return ast.WalkContinue, nil
}

// MARK: inlines

func (r *MarkdownRenderer) renderEmphasis(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)
	if entering {
		w.WriteString(strings.Repeat("*", n.Level))
	} else {
		w.WriteString(strings.Repeat("*", n.Level))
	}
	return ast.WalkContinue, nil
}

func (r *MarkdownRenderer) renderTaskCheckBox(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*extast.TaskCheckBox)
	if entering {
		if n.IsChecked {
			_, _ = w.WriteString("[x] ")
		} else {
			_, _ = w.WriteString("[ ] ")
		}
	}
	return ast.WalkContinue, nil
}

func (r *MarkdownRenderer) renderLink(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		// NOTE As of goldmark v1.7.1, ast.Link.Title is not set by default markdown parser, so that n.Text(source) is used here instead.
		// n.Text(source) retrieves text from n's child node (ast.Text) in advance to node-walking operation.
		w.WriteString(fmt.Sprintf("[%s](%s)", n.Text(source), n.Destination))
	}
	return ast.WalkContinue, nil
}

func (r *MarkdownRenderer) renderAutoLink(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.AutoLink)
	if entering {
		w.WriteString(fmt.Sprint(string(n.URL(source))))
	}
	return ast.WalkContinue, nil
}
