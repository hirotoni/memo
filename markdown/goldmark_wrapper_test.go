package markdown

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var updateGolden = false

func TestGoldmarkWrapper_Render(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "sample",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := "./testdata/sample.md"
			golden := "./testdata/sample.md.golden"

			f, err := os.ReadFile(filename)
			assert.NoError(err)

			gmw := NewGoldmarkWrapper()
			doc := gmw.Parse(f)
			// doc.Dump(f, 1)

			writer := &bytes.Buffer{}
			gmw.Render(writer, f, doc)

			// fmt.Println(writer.String())
			assert.Equal(string(f), writer.String())

			if updateGolden {
				os.WriteFile(golden, writer.Bytes(), 0644)
			}
		})
	}
}

func TestGoldmarkWrapper_FindHeadingAndGetHangingNodes(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "sample",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := "./testdata/sample.md"
			f, err := os.ReadFile(filename)
			assert.NoError(err)

			gmw := NewGoldmarkWrapper()
			h := NewHeading(2, "ordered list")
			_, hangingNodes := gmw.FindHeadingAndGetHangingNodes(f, h)

			assert.NotEmpty(hangingNodes)
		})
	}

}
func TestGoldmarkWrapper_AppendTextAfterHeadingBlock(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name          string
		inputMarkdown string
		targetHeading Heading
		textToAppend  string
		expected      string
	}{
		{
			name: "append text after heading with no children",
			inputMarkdown: `# Heading 1
## Heading 2
Content under heading 2.`,
			targetHeading: NewHeading(2, "Heading 2"),
			textToAppend:  "Appended text.",
			expected: `# Heading 1
## Heading 2
Content under heading 2.

Appended text.`,
		},
		{
			name: "append text after heading with children",
			inputMarkdown: `# Heading 1
## Heading 2
Content under heading 2.
### Heading 3
Content under heading 3.`,
			targetHeading: NewHeading(2, "Heading 2"),
			textToAppend:  "Appended text.",
			expected: `# Heading 1
## Heading 2
Content under heading 2.
### Heading 3
Content under heading 3.

Appended text.`,
		},
		{
			name: "append text after heading with no matching heading",
			inputMarkdown: `# Heading 1
## Heading 2
Content under heading 2.`,
			targetHeading: NewHeading(3, "Non-existent Heading"),
			textToAppend:  "Appended text.",
			expected: `# Heading 1
## Heading 2
Content under heading 2.`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gmw := NewGoldmarkWrapper()
			result := gmw.InsertTextAfterHeadingBlock([]byte(tt.inputMarkdown), tt.targetHeading, tt.textToAppend)
			assert.Equal(tt.expected, string(result))
		})
	}
}
