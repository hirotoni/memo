package markdown

import (
	"bytes"
	"io"
	"slices"
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
	return &GoldmarkWrapper{
		Goldmark: goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithRendererOptions(
				renderer.WithNodeRenderers(
					util.Prioritized(myrenderer.NewMarkdownRenderer(), 1),
				),
			),
		),
	}
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

func (gmw *GoldmarkWrapper) GetHeadingNodes(doc ast.Node, source []byte, level int) []ast.Node {
	document := doc.OwnerDocument()
	if document == nil {
		return nil
	}
	var foundNodes []ast.Node
	for c := document.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Kind() == ast.KindHeading {
			levelMatched := c.(*ast.Heading).Level == level
			if levelMatched {
				foundNodes = append(foundNodes, c)
			}
		}
	}
	return foundNodes
}

func (gmw *GoldmarkWrapper) GetHeadingNode(doc ast.Node, source []byte, heading Heading) ast.Node {
	document := doc.OwnerDocument()
	if document == nil {
		return nil
	}

	var foundNode ast.Node
	for c := document.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Kind() == ast.KindHeading {
			levelMatched := c.(*ast.Heading).Level == heading.Level
			textMatched := strings.Contains(string(c.Text(source)), heading.Text)
			if levelMatched && textMatched {
				foundNode = c
				break
			}
		}
	}
	return foundNode
}

// FindHeadingAndGetHangingNodes finds a heading that matches given text and level, then returns the hanging nodes of the heading
func (gmw *GoldmarkWrapper) FindHeadingAndGetHangingNodes(doc ast.Node, source []byte, heading Heading) []ast.Node {
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
				levelMatched := c.(*ast.Heading).Level == heading.Level
				textMatched := strings.Contains(string(c.Text(source)), heading.Text)
				if levelMatched && textMatched {
					mode = modeExiting
				}
			case modeExiting:
				if c.(*ast.Heading).Level <= heading.Level {
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

// InsertNodesAfter inserts nodes to document at target position, and returns updated byte array of document as the result of the insert operation
func (gmw *GoldmarkWrapper) InsertNodesAfter(doc ast.Node, targetNode ast.Node, nodesToInsert []ast.Node, sourceSelf, sourceNodesToInsert []byte) []byte {
	// insert from tail nodes
	slices.Reverse(nodesToInsert)

	for _, n := range nodesToInsert {
		doc.InsertAfter(doc, targetNode, n)
		s := targetNode.Lines().At(0) // TODO error handling

		tmp := new(bytes.Buffer)
		gmw.Render(tmp, sourceNodesToInsert, n)

		buf := []byte{}
		buf = append(buf, sourceSelf[:s.Stop]...)
		buf = append(buf, tmp.Bytes()...)
		buf = append(buf, sourceSelf[s.Stop:]...)
		sourceSelf = buf
	}

	return sourceSelf
}

func (gmw *GoldmarkWrapper) InsertTextAfter(sourceSelf []byte, targetHeading Heading, text string) []byte {
	doc := gmw.Parse(sourceSelf)
	targetHeaderNode := gmw.GetHeadingNode(doc, sourceSelf, targetHeading)

	s := targetHeaderNode.Lines().At(0)

	buf := []byte{}
	buf = append(buf, sourceSelf[:s.Stop]...)
	buf = append(buf, []byte("\n\n"+text)...)
	buf = append(buf, sourceSelf[s.Stop:]...)
	sourceSelf = buf

	return sourceSelf
}
