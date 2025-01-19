package repos

import (
	"testing"

	"github.com/hirotoni/memo/configs"
	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
	"github.com/stretchr/testify/assert"
)

func TestMemoArchiveNodeRepo_MemoArchiveNodesFromMemoArchivesDir(t *testing.T) {
	type fields struct {
		config *configs.TomlConfig
		gmw    *markdown.GoldmarkWrapper
	}
	type args struct {
		shown []*models.MemoArchive
	}

	testConfig := configs.LoadTomlConfig()
	testConfig.BaseDir = "testdata"

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*models.MemoArchiveNode
	}{
		{
			name: "Test MemoArchiveNodesFromMemoArchviesDir",
			fields: fields{
				config: testConfig,
				gmw:    markdown.NewGoldmarkWrapper(),
			},
			args: args{
				shown: []*models.MemoArchive{},
			},
			want: []*models.MemoArchiveNode{
				{
					Kind: models.MEMOARCHIVENODEKIND_DIR, Text: "testmemoarchives", Depth: 0,
				},
				{
					Kind: models.MEMOARCHIVENODEKIND_TITLE, Text: "testmemoarchive", Depth: 1,
					MemoArchive: models.MemoArchive{Text: "testmemoarchive", Destination: "../memoarchives/testmemoarchives/testmemoarchive-1.md", Checked: false},
				},
				{
					Kind: models.MEMOARCHIVENODEKIND_MEMO, Text: "super duper pepper", Depth: 2,
					MemoArchive: models.MemoArchive{Text: "super duper pepper", Destination: "../memoarchives/testmemoarchives/testmemoarchive-1.md#super-duper-pepper", Checked: false},
				},
				{
					Kind: models.MEMOARCHIVENODEKIND_TITLE, Text: "testmemoarchive", Depth: 1,
					MemoArchive: models.MemoArchive{Text: "testmemoarchive", Destination: "../memoarchives/testmemoarchives/testmemoarchive-2.md", Checked: false},
				},
				{
					Kind: models.MEMOARCHIVENODEKIND_MEMO, Text: "super duper pepper", Depth: 2,
					MemoArchive: models.MemoArchive{Text: "super duper pepper", Destination: "../memoarchives/testmemoarchives/testmemoarchive-2.md#super-duper-pepper", Checked: false},
				},
			},
		},
		{
			name: "Test MemoArchiveNodesFromMemoArchivesDir checked",
			fields: fields{
				config: testConfig,
				gmw:    markdown.NewGoldmarkWrapper(),
			},
			args: args{
				shown: []*models.MemoArchive{
					{Text: "testmemoarchive", Destination: "../memoarchives/testmemoarchives/testmemoarchive-1.md#super-duper-pepper", Checked: false},
				},
			},
			want: []*models.MemoArchiveNode{
				{
					Kind: models.MEMOARCHIVENODEKIND_DIR, Text: "testmemoarchives", Depth: 0,
				},
				{
					Kind: models.MEMOARCHIVENODEKIND_TITLE, Text: "testmemoarchive", Depth: 1,
					MemoArchive: models.MemoArchive{Text: "testmemoarchive", Destination: "../memoarchives/testmemoarchives/testmemoarchive-1.md", Checked: false},
				},
				{
					Kind: models.MEMOARCHIVENODEKIND_MEMO, Text: "super duper pepper", Depth: 2,
					MemoArchive: models.MemoArchive{Text: "super duper pepper", Destination: "../memoarchives/testmemoarchives/testmemoarchive-1.md#super-duper-pepper", Checked: true},
				},
				{
					Kind: models.MEMOARCHIVENODEKIND_TITLE, Text: "testmemoarchive", Depth: 1,
					MemoArchive: models.MemoArchive{Text: "testmemoarchive", Destination: "../memoarchives/testmemoarchives/testmemoarchive-2.md", Checked: false},
				},
				{
					Kind: models.MEMOARCHIVENODEKIND_MEMO, Text: "super duper pepper", Depth: 2,
					MemoArchive: models.MemoArchive{Text: "super duper pepper", Destination: "../memoarchives/testmemoarchives/testmemoarchive-2.md#super-duper-pepper", Checked: false},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMemoArchiveNodeRepo(tt.fields.config)
			assert := assert.New(t)
			assert.Equal(tt.want, repo.MemoArchiveNodesFromMemoArchivesDir(tt.args.shown))
		})
	}
}
