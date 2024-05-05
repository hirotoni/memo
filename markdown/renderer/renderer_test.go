package renderer

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/ast"
)

func TestMarkdownRenderer_renderHeading(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		node     ast.Node
		entering bool
	}

	tests := []struct {
		name       string
		args       args
		wantStatus ast.WalkStatus
		wantStr    string
		wantErr    bool
	}{
		{
			name: "entering, level 1",
			args: args{
				node:     &ast.Heading{Level: 1},
				entering: true,
			},
			wantStatus: ast.WalkContinue,
			wantStr:    "# ",
			wantErr:    false,
		},
		{
			name: "entering, level 6",
			args: args{
				node:     &ast.Heading{Level: 6},
				entering: true,
			},
			wantStatus: ast.WalkContinue,
			wantStr:    "###### ",
			wantErr:    false,
		},
		{
			name: "exiting, level 1",
			args: args{
				node:     &ast.Heading{Level: 1},
				entering: false,
			},
			wantStatus: ast.WalkContinue,
			wantStr:    "",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewMarkdownRenderer()
			source := []byte("dummy source")

			buf := new(strings.Builder)
			w := bufio.NewWriter(buf)

			got, err := r.renderHeading(w, source, tt.args.node, tt.args.entering)

			if tt.wantErr {
				assert.NotNil(err)
			}

			assert.Nil(err)
			assert.Equal(tt.wantStatus, got)

			assert.NoError(w.Flush())
			assert.Equal(tt.wantStr, buf.String())
		})
	}
}
