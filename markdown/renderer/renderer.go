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
	reg.Register(extast.KindTaskCheckBox, r.renderTaskCheckBox)

	// // inlines
	// reg.Register(ast.KindAutoLink, r.renderAutoLink)
	// reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	// reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindLink, r.renderLink)
	// reg.Register(ast.KindRawHTML, r.renderRawHTML)
	reg.Register(ast.KindText, r.renderText)
	// reg.Register(ast.KindString, r.renderString)
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
		if n.Parent().Kind() == ast.KindLink { // TODO needs nil check?
			// r.renderLink() renders text in advance. no rendering needed here.
			return ast.WalkContinue, nil
		}

		value := n.Text(source)
		w.WriteString(string(value))

		sibling := node.NextSibling()
		if n.SoftLineBreak() {
			if sibling != nil && sibling.Kind() == ast.KindText {
				if siblingText := sibling.(*ast.Text).Text(source); len(siblingText) != 0 {
					pp := n.Parent().Parent()
					switch pp := pp.(type) {
					case *ast.ListItem:
						w.WriteString("\n")
						w.WriteString(strings.Repeat(" ", pp.Offset))
					default:
						w.WriteString("\n")
					}
				}
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

	if entering {
		// リストの最初の要素でない場合、あるいは入れ子のリスト要素である場合は改行する
		if n.PreviousSibling() != nil || n.Parent().Parent().Kind() == ast.KindListItem {
			w.WriteString("\n")
		}

		// 入れ子になっている分、インデントを増やす
		var cur = n.Parent().Parent()
		for {
			if li, ok := cur.(*ast.ListItem); ok {
				w.WriteString(strings.Repeat(" ", li.Offset))
				cur = cur.Parent().Parent()
			} else {
				break
			}
		}

		p := n.Parent().(*ast.List)
		if p.IsOrdered() {
			order := p.Start
			cur := node
			for {
				if cur.PreviousSibling() != nil {
					order++
					cur = cur.PreviousSibling()
				} else {
					break
				}
			}
			w.WriteString(fmt.Sprintf("%d", order))
			w.WriteString(string(p.Marker) + " ")
		} else {
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
		// NOTE As of goldmark v1.7.1, ast.Link.Title is not set by default markdown parser, so that n.Text(source) is used here.
		// n.Text(source) retrieves text from n's child node (ast.Text) in advance to node-walking operation.
		w.WriteString(fmt.Sprintf("[%s](%s)", n.Text(source), n.Destination))
	}
	return ast.WalkContinue, nil
}
