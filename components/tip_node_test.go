package components

import (
	"bytes"
	"testing"

	"github.com/hirotoni/memo/models"
	"github.com/stretchr/testify/assert"
)

func TestPrintMemoArchiveNode(t *testing.T) {
	type args struct {
		b  *bytes.Buffer
		tn *models.MemoArchiveNode
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "DIR",
			args: args{
				b: &bytes.Buffer{},
				tn: &models.MemoArchiveNode{
					Kind:  models.MEMOARCHIVENODEKIND_DIR,
					Depth: 1,
					Text:  "text",
					MemoArchive: models.MemoArchive{
						Text:        "text",
						Destination: "destination",
						Checked:     false,
					},
				},
			},
			want: "  - text\n",
		},
		{
			name: "TITLE",
			args: args{
				b: &bytes.Buffer{},
				tn: &models.MemoArchiveNode{
					Kind:  models.MEMOARCHIVENODEKIND_TITLE,
					Depth: 1,
					Text:  "text",
					MemoArchive: models.MemoArchive{
						Text:        "text",
						Destination: "destination",
						Checked:     false,
					},
				},
			},
			want: "  - text\n",
		},
		{
			name: "MEMOARCHIVE unchecked",
			args: args{
				b: &bytes.Buffer{},
				tn: &models.MemoArchiveNode{
					Kind:  models.MEMOARCHIVENODEKIND_MEMO,
					Depth: 1,
					Text:  "text",
					MemoArchive: models.MemoArchive{
						Text:        "text",
						Destination: "destination",
						Checked:     false,
					},
				},
			},
			want: "  - [ ] [text](destination)\n",
		},
		{
			name: "MEMOARCHIVE checked",
			args: args{
				b: &bytes.Buffer{},
				tn: &models.MemoArchiveNode{
					Kind:  models.MEMOARCHIVENODEKIND_MEMO,
					Depth: 1,
					Text:  "text",
					MemoArchive: models.MemoArchive{
						Text:        "text",
						Destination: "destination",
						Checked:     true,
					},
				},
			},
			want: "  - [x] [text](destination)\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			PrintMemoArchiveNode(tt.args.b, tt.args.tn)
			assert.Equal(tt.want, tt.args.b.String())
		})
	}
}

func TestPrintMemoArchiveNodeHeadingStyleSlice(t *testing.T) {
	type args struct {
		b   *bytes.Buffer
		tns []*models.MemoArchiveNode
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "DIR",
			args: args{
				b: &bytes.Buffer{},
				tns: []*models.MemoArchiveNode{
					{
						Kind:  models.MEMOARCHIVENODEKIND_DIR,
						Depth: 1,
						Text:  "text",
						MemoArchive: models.MemoArchive{
							Text:        "text",
							Destination: "destination",
							Checked:     false,
						},
					},
				},
			},
			want: "### text\n\n",
		},
		{
			name: "TITLE",
			args: args{
				b: &bytes.Buffer{},
				tns: []*models.MemoArchiveNode{
					{
						Kind:  models.MEMOARCHIVENODEKIND_TITLE,
						Depth: 1,
						Text:  "text",
						MemoArchive: models.MemoArchive{
							Text:        "text",
							Destination: "destination",
							Checked:     false,
						},
					},
				},
			},
			want: "### text\n\n",
		},
		{
			name: "MEMOARCHIVE unchecked",
			args: args{
				b: &bytes.Buffer{},
				tns: []*models.MemoArchiveNode{
					{
						Kind:  models.MEMOARCHIVENODEKIND_MEMO,
						Depth: 1,
						Text:  "text",
						MemoArchive: models.MemoArchive{
							Text:        "text",
							Destination: "destination",
							Checked:     false,
						},
					},
				},
			},
			want: "- [ ] [text](destination)\n",
		},
		{
			name: "MEMOARCHIVE checked",
			args: args{
				b: &bytes.Buffer{},
				tns: []*models.MemoArchiveNode{
					{
						Kind:  models.MEMOARCHIVENODEKIND_MEMO,
						Depth: 1,
						Text:  "text",
						MemoArchive: models.MemoArchive{
							Text:        "text",
							Destination: "destination",
							Checked:     true,
						},
					},
				},
			},
			want: "- [x] [text](destination)\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			PrintMemoArchiveNodeHeadingStyle(tt.args.b, tt.args.tns)
			assert.Equal(tt.want, tt.args.b.String())
		})
	}
}
