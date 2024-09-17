package repos

import (
	"testing"

	"github.com/hirotoni/memo/config"
	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
	"github.com/stretchr/testify/assert"
)

func TestTipNodeRepo_TipNodesFromTipsDir(t *testing.T) {
	type fields struct {
		config *config.TomlConfig
		gmw    *markdown.GoldmarkWrapper
	}
	type args struct {
		shown []*models.Tip
	}

	testConfig := config.LoadTomlConfig()
	testConfig.BaseDir = "testdata"

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*models.TipNode
	}{
		{
			name: "Test TipNodesFromTipsDir",
			fields: fields{
				config: testConfig,
				gmw:    markdown.NewGoldmarkWrapper(),
			},
			args: args{
				shown: []*models.Tip{},
			},
			want: []*models.TipNode{
				{
					Kind: models.TIPNODEKIND_DIR, Text: "testtips", Depth: 0,
				},
				{
					Kind: models.TIPNODEKIND_TITLE, Text: "testtip", Depth: 1,
					Tip: models.Tip{Text: "testtip", Destination: "../tips/testtips/testtip-1.md", Checked: false},
				},
				{
					Kind: models.TIPNODEKIND_TIP, Text: "super duper tipper", Depth: 2,
					Tip: models.Tip{Text: "super duper tipper", Destination: "../tips/testtips/testtip-1.md#super-duper-tipper", Checked: false},
				},
				{
					Kind: models.TIPNODEKIND_TITLE, Text: "testtip", Depth: 1,
					Tip: models.Tip{Text: "testtip", Destination: "../tips/testtips/testtip-2.md", Checked: false},
				},
				{
					Kind: models.TIPNODEKIND_TIP, Text: "super duper tipper", Depth: 2,
					Tip: models.Tip{Text: "super duper tipper", Destination: "../tips/testtips/testtip-2.md#super-duper-tipper", Checked: false},
				},
			},
		},
		{
			name: "Test TipNodesFromTipsDir checked",
			fields: fields{
				config: testConfig,
				gmw:    markdown.NewGoldmarkWrapper(),
			},
			args: args{
				shown: []*models.Tip{
					{Text: "testtip", Destination: "../tips/testtips/testtip-1.md#super-duper-tipper", Checked: false},
				},
			},
			want: []*models.TipNode{
				{
					Kind: models.TIPNODEKIND_DIR, Text: "testtips", Depth: 0,
				},
				{
					Kind: models.TIPNODEKIND_TITLE, Text: "testtip", Depth: 1,
					Tip: models.Tip{Text: "testtip", Destination: "../tips/testtips/testtip-1.md", Checked: false},
				},
				{
					Kind: models.TIPNODEKIND_TIP, Text: "super duper tipper", Depth: 2,
					Tip: models.Tip{Text: "super duper tipper", Destination: "../tips/testtips/testtip-1.md#super-duper-tipper", Checked: true},
				},
				{
					Kind: models.TIPNODEKIND_TITLE, Text: "testtip", Depth: 1,
					Tip: models.Tip{Text: "testtip", Destination: "../tips/testtips/testtip-2.md", Checked: false},
				},
				{
					Kind: models.TIPNODEKIND_TIP, Text: "super duper tipper", Depth: 2,
					Tip: models.Tip{Text: "super duper tipper", Destination: "../tips/testtips/testtip-2.md#super-duper-tipper", Checked: false},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTipNodeRepo(tt.fields.config, tt.fields.gmw)
			assert := assert.New(t)
			assert.Equal(tt.want, repo.TipNodesFromTipsDir(tt.args.shown))
		})
	}
}
