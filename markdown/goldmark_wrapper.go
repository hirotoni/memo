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

type Heading struct {
	Level int
	Text  string
}

func NewHeading(level int, text string) Heading {
	return Heading{Level: level, Text: text}
}

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

func (gmw *GoldmarkWrapper) RenderSlice(writer io.Writer, input []byte, nodes []ast.Node) error {
	for _, n := range nodes {
		err := gmw.Render(writer, input, n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gmw *GoldmarkWrapper) GetHeadingNodesByLevel(source []byte, level int) (ast.Node, []ast.Node) {
	doc := gmw.Parse(source)
	var foundNodes []ast.Node
	for c := doc.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Kind() == ast.KindHeading {
			levelMatched := c.(*ast.Heading).Level == level
			if levelMatched {
				foundNodes = append(foundNodes, c)
			}
		}
	}
	return doc, foundNodes
}

func (gmw *GoldmarkWrapper) GetHeadingNode(source []byte, heading Heading) (ast.Node, ast.Node) {
	doc := gmw.Parse(source)

	var foundNode ast.Node
	for c := doc.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Kind() == ast.KindHeading {
			levelMatched := c.(*ast.Heading).Level == heading.Level
			textMatched := strings.Contains(string(c.Text(source)), heading.Text)
			if levelMatched && textMatched {
				foundNode = c
				break
			}
		}
	}
	return doc, foundNode
}

// FindHeadingAndGetHangingNodes finds a heading that matches given text and level, then returns the found heading and hanging nodes of the heading
func (gmw *GoldmarkWrapper) FindHeadingAndGetHangingNodes(source []byte, heading Heading) (ast.Node, []ast.Node) {
	doc := gmw.Parse(source)

	const (
		modeSearching = iota
		modeExiting
	)

	var foundHeading ast.Node
	var hangingNodes []ast.Node
	var mode = modeSearching

loop:
	for c := doc.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Kind() == ast.KindHeading {
			switch mode {
			case modeSearching:
				levelMatched := c.(*ast.Heading).Level == heading.Level
				textMatched := strings.Contains(string(c.Text(source)), heading.Text)
				if levelMatched && textMatched {
					foundHeading = c
					mode = modeExiting
				}
			case modeExiting:
				if c.(*ast.Heading).Level <= heading.Level {
					break loop
				} else {
					hangingNodes = append(hangingNodes, c)
				}
			}
		} else {
			switch mode {
			case modeSearching:
				continue
			case modeExiting:
				hangingNodes = append(hangingNodes, c)
			}
		}
	}

	return foundHeading, hangingNodes
}

// InsertNodesAfter inserts nodes to document at target position, and returns updated byte array of document as the result of the insert operation
func (gmw *GoldmarkWrapper) InsertNodesAfter(sourceSelf []byte, targetHeading Heading, sourceNodesToInsert []byte, nodesToInsert []ast.Node) []byte {
	// insert from tail nodes
	slices.Reverse(nodesToInsert)

	doc, targetNode := gmw.GetHeadingNode(sourceSelf, targetHeading)

	// for compaitibility
	if targetNode == nil {
		// TODO return info that target heading is not found
		return sourceSelf
	}

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

// InsertTextAfter inserts text to document at target position, and returns updated byte array of document as the result of the insert operation
func (gmw *GoldmarkWrapper) InsertTextAtHeadingStart(sourceSelf []byte, targetHeading Heading, text string) []byte {
	_, targetHeadingNode := gmw.GetHeadingNode(sourceSelf, targetHeading)

	// for compaitibility
	if targetHeadingNode == nil {
		// TODO return info that target heading is not found
		return sourceSelf
	}

	s := targetHeadingNode.Lines().At(0)

	buf := []byte{}
	buf = append(buf, sourceSelf[:s.Stop]...)
	buf = append(buf, []byte("\n\n"+text)...)
	buf = append(buf, sourceSelf[s.Stop:]...)
	sourceSelf = buf

	return sourceSelf
}

// InsertTextAfter inserts text to document at target position, and returns updated byte array of document as the result of the insert operation
func (gmw *GoldmarkWrapper) InsertTextAtHeadingEnd(sourceSelf []byte, targetHeading Heading, text string) []byte {
	foundHeading, children := gmw.FindHeadingAndGetHangingNodes(sourceSelf, targetHeading)
	if foundHeading == nil {
		// TODO return info that target heading is not found
		return sourceSelf
	}

	// determine the last position of the found children
	var last int
	if len(children) > 0 {
		last = children[len(children)-1].Lines().At(0).Stop
	} else {
		last = foundHeading.Lines().At(0).Stop
	}

	// append text after the last position
	buf := []byte{}
	buf = append(buf, sourceSelf[:last]...)
	buf = append(buf, []byte("\n\n"+text)...)
	sourceSelf = buf

	return sourceSelf
}
