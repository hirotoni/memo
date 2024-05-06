package markdown

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			f, err := os.ReadFile("./testdata/sample.md")
			assert.NoError(err)

			gmw := NewGoldmarkWrapper()
			doc := gmw.Parse(f)
			doc.Dump(f, 1)

			writer := &bytes.Buffer{}
			gmw.Render(writer, f, doc)

			fmt.Println(writer.String())
		})
	}
}
