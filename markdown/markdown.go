package markdown

import (
	"bytes"
	"io"
	"strings"

	myrenderer "github.com/hirotoni/memo/markdown/renderer"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type GoldmarkWrapper struct {
	Goldmark goldmark.Markdown
}

func NewGoldmarkWrapper() *GoldmarkWrapper {
	gm := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			renderer.WithNodeRenderers(
				util.Prioritized(myrenderer.NewMarkdownRenderer(), 1),
			),
		))

	wrapper := GoldmarkWrapper{
		Goldmark: gm,
	}

	return &wrapper
}

func (gmw *GoldmarkWrapper) Parse(source []byte) ast.Node {
	reader := text.NewReader(source)
	doc := gmw.Goldmark.Parser().Parse(reader)
	return doc
}

func (gmw *GoldmarkWrapper) Render(writer io.Writer, input []byte, doc ast.Node) error {
	err := gmw.Goldmark.Renderer().Render(writer, input, doc)
	if err != nil {
		return err
	}
	return nil
}

func (gmw *GoldmarkWrapper) GetHeadingNode(doc ast.Node, source []byte, text string, level int) ast.Node {
	// TODO define (text, level) type struct

	document := doc.OwnerDocument()
	if document == nil {
		return nil
	}

	var foundNode ast.Node
	for c := document.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Kind() == ast.KindHeading {
			levelMatched := c.(*ast.Heading).Level == level
			textMatched := strings.Contains(string(c.Text(source)), text)
			if levelMatched && textMatched {
				foundNode = c
				break
			}
		}
	}
	return foundNode
}

// FindHeadingAndGetHangingNodes finds a heading that matches given text and level, then returns the hanging nodes of the heading
func (gmw *GoldmarkWrapper) FindHeadingAndGetHangingNodes(doc ast.Node, source []byte, text string, level int) []ast.Node {
	// TODO define (text, level) type struct

	document := doc.OwnerDocument()
	if document == nil {
		return nil
	}

	const (
		modeSearching = iota
		modeExiting
	)

	mode := modeSearching
	resultNodes := []ast.Node{}

loop:
	for c := doc.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Kind() == ast.KindHeading {
			switch mode {
			case modeSearching:
				levelMatched := c.(*ast.Heading).Level == level
				textMatched := strings.Contains(string(c.Text(source)), text)
				if levelMatched && textMatched {
					mode = modeExiting
				}
			case modeExiting:
				if c.(*ast.Heading).Level <= level {
					break loop
				} else {
					resultNodes = append(resultNodes, c)
				}
			}
		} else {
			switch mode {
			case modeSearching:
				continue
			case modeExiting:
				resultNodes = append(resultNodes, c)
			}
		}
	}

	return resultNodes
}

// InsertAfter inserts insertees to self node at taregt position, and returns updated byte array of self node as the result of the insertion
func (gmw *GoldmarkWrapper) InsertAfter(self ast.Node, target ast.Node, insertees []ast.Node, selfSource, nodeSource []byte) []byte {
	// insert from tail nodes
	length := len(insertees)
	for i := range make([]int, length) {
		n := insertees[length-i-1]
		self.InsertAfter(self, target, n)
		s := target.Lines().At(0) // TODO error handling

		tmp := new(bytes.Buffer)
		gmw.Render(tmp, nodeSource, n)

		buf := []byte{}
		buf = append(buf, selfSource[:s.Stop]...)
		buf = append(buf, tmp.Bytes()...)
		buf = append(buf, selfSource[s.Stop:]...)
		selfSource = buf
	}

	return selfSource
}