package renderer

import (
	"bufio"
	"reflect"
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

func TestMarkdownRenderer_renderEmphasis(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering true",
			args:  args{node: &ast.Emphasis{Level: 1}, entering: true},
			wants: wants{status: ast.WalkContinue, str: "*", err: false},
		},
		{
			name:  "entering false",
			args:  args{node: &ast.Emphasis{Level: 1}, entering: false},
			wants: wants{status: ast.WalkContinue, str: "*", err: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewMarkdownRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)
			dummySource := []byte("dummy source")

			got, err := r.renderEmphasis(w, dummySource, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderEmphasis() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderEmphasis() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(t, w.Flush())
			assert.Equal(t, tt.wants.str, sb.String())
		})
	}
}
