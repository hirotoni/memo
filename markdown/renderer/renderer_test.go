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
			name: "entering, level 1",
			args: args{
				node:     genHeaderNode(1, false),
				entering: true,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    "# ",
				err:    false,
			},
		},
		{
			name: "exiting, level 1",
			args: args{
				node:     genHeaderNode(1, false),
				entering: false,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    "",
				err:    false,
			},
		},
		{
			name: "entering, level 6",
			args: args{
				node:     genHeaderNode(6, false),
				entering: true,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    "###### ",
				err:    false,
			},
		},
		{
			name: "entering, blank previous lines",
			args: args{
				node:     genHeaderNode(1, true),
				entering: true,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    "\n\n# ",
				err:    false,
			},
		},
		{
			name: "exiting, blank previous lines",
			args: args{
				node:     genHeaderNode(1, true),
				entering: false,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    "",
				err:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewMarkdownRenderer()
			source := []byte("dummy source")

			buf := new(strings.Builder)
			w := bufio.NewWriter(buf)

			got, err := r.renderHeading(w, source, tt.args.node, tt.args.entering)

			if tt.wants.err {
				assert.NotNil(err)
			}

			assert.Nil(err)
			assert.Equal(tt.wants.status, got)

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, buf.String())
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
			args:  args{node: genEnphasisNode(1), entering: true},
			wants: wants{status: ast.WalkContinue, str: "*", err: false},
		},
		{
			name:  "entering false",
			args:  args{node: genEnphasisNode(1), entering: false},
			wants: wants{status: ast.WalkContinue, str: "*", err: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewMarkdownRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)
			source := []byte("dummy source")

			got, err := r.renderEmphasis(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderEmphasis() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderEmphasis() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

// MARK: inlines

func TestMarkdownRenderer_renderTaskCheckBox(t *testing.T) {
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
			name:  "isChecked true, entering true",
			args:  args{node: genTaskCheckBoxNode(true), entering: true},
			wants: wants{status: ast.WalkContinue, str: "[x] ", err: false},
		},
		{
			name:  "isChecked false, entering true",
			args:  args{node: genTaskCheckBoxNode(false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "[ ] ", err: false},
		},
		{
			name:  "isChecked true, entering false",
			args:  args{node: genTaskCheckBoxNode(true), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
		{
			name:  "isChecked false, entering false",
			args:  args{node: genTaskCheckBoxNode(false), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewMarkdownRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)
			dummySource := []byte("dummy source")

			got, err := r.renderTaskCheckBox(w, dummySource, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderTaskCheckBox() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderTaskCheckBox() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderLink(t *testing.T) {
	type args struct {
		source   []byte
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
			name: "entering true",
			args: args{
				source:   []byte("text"),
				node:     genLinkNode([]byte("text"), []byte("destination")),
				entering: true,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    "[text](destination)",
				err:    false,
			},
		},
		{
			name: "entering false",
			args: args{
				source:   []byte("text"),
				node:     genLinkNode([]byte("text"), []byte("destination")),
				entering: false,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    "",
				err:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewMarkdownRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderLink(w, tt.args.source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderLink() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderLink() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderAutoLink(t *testing.T) {
	type args struct {
		source   []byte
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("https://example.com")

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "entering true",
			args: args{
				source:   source,
				node:     genAutoLinkNode(source),
				entering: true,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    "https://example.com",
				err:    false,
			},
		},
		{
			name: "entering false",
			args: args{
				source:   source,
				node:     genAutoLinkNode(source),
				entering: false,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    "",
				err:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewMarkdownRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderAutoLink(w, tt.args.source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderAutoLink() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderAutoLink() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderText(t *testing.T) {
	type args struct {
		source   []byte
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("test text")

	offset := 2
	li := ast.NewListItem(offset)
	tb := ast.NewTextBlock()
	tb.AppendChild(li, tb)

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "entering",
			args: args{
				source:   source,
				node:     genTextNode(source, false, ast.NewTextBlock()),
				entering: true,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    string(source),
				err:    false,
			},
		},
		{
			name: "entering, soft line break",
			args: args{
				source:   source,
				node:     genTextNode(source, true, ast.NewTextBlock()),
				entering: true,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    string(source) + "\n",
				err:    false,
			},
		},
		{
			name: "entering, parent link",
			args: args{
				source:   source,
				node:     genTextNode(source, true, ast.NewLink()),
				entering: true,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    "",
				err:    false,
			},
		},
		{
			name: "entering, parent parent listitem",
			args: args{
				source:   source,
				node:     genTextNode(source, true, tb),
				entering: true,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    string(source) + "\n" + strings.Repeat(" ", offset),
				err:    false,
			},
		},
		{
			name: "entering, parent nil",
			args: args{
				source:   source,
				node:     genTextNode(source, true, nil),
				entering: true,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    "",
				err:    false,
			},
		},
		{
			name: "entering false",
			args: args{
				source:   source,
				node:     genTextNode(source, false, ast.NewTextBlock()),
				entering: false,
			},
			wants: wants{
				status: ast.WalkContinue,
				str:    "",
				err:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewMarkdownRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderText(w, tt.args.source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderText() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderText() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}
