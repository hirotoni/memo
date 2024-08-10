package usecases

import (
	"bytes"
	"testing"

	"github.com/hirotoni/memo/models"
	"github.com/stretchr/testify/assert"
)

func TestPrintTipNode(t *testing.T) {
	type args struct {
		b  *bytes.Buffer
		tn *models.TipNode
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
				tn: &models.TipNode{
					Kind:  models.TIPNODEKIND_DIR,
					Depth: 1,
					Text:  "text",
					Tip: models.Tip{
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
				tn: &models.TipNode{
					Kind:  models.TIPNODEKIND_TITLE,
					Depth: 1,
					Text:  "text",
					Tip: models.Tip{
						Text:        "text",
						Destination: "destination",
						Checked:     false,
					},
				},
			},
			want: "  - text\n",
		},
		{
			name: "TIP unchecked",
			args: args{
				b: &bytes.Buffer{},
				tn: &models.TipNode{
					Kind:  models.TIPNODEKIND_TIP,
					Depth: 1,
					Text:  "text",
					Tip: models.Tip{
						Text:        "text",
						Destination: "destination",
						Checked:     false,
					},
				},
			},
			want: "  - [ ] [text](destination)\n",
		},
		{
			name: "TIP checked",
			args: args{
				b: &bytes.Buffer{},
				tn: &models.TipNode{
					Kind:  models.TIPNODEKIND_TIP,
					Depth: 1,
					Text:  "text",
					Tip: models.Tip{
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
			PrintTipNode(tt.args.b, tt.args.tn)
			assert.Equal(tt.want, tt.args.b.String())
		})
	}
}
