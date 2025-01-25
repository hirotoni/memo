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
