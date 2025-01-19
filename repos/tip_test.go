package repos

import (
	"reflect"
	"testing"

	"github.com/hirotoni/memo/configs"
	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
)

func TestMemoArchiveRepo_MemoArchivesFromIndex(t *testing.T) {
	type fields struct {
		config *configs.TomlConfig
		gmw    *markdown.GoldmarkWrapper
	}

	testConfig := configs.LoadTomlConfig()
	testConfig.BaseDir = "testdata"

	tests := []struct {
		name   string
		fields fields
		want   []*models.MemoArchive
	}{
		{
			name: "Test MemoArchivesFromIndex",
			fields: fields{
				config: testConfig,
				gmw:    markdown.NewGoldmarkWrapper(),
			},
			want: []*models.MemoArchive{
				{Text: "some memoarchive", Destination: "somewhere", Checked: false},
				{Text: "another memoarchive", Destination: "anywhere", Checked: true},
				{Text: "yet another memoarchive", Destination: "everywhere", Checked: false},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MemoArchiveRepo{
				config: tt.fields.config,
			}
			if got := repo.MemoArchivesFromIndex(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoArchiveRepo.MemoArchivesFromIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoArchiveRepo_MemoArchivesFromIndexChecked(t *testing.T) {
	type fields struct {
		config *configs.TomlConfig
		gmw    *markdown.GoldmarkWrapper
	}

	testConfig := configs.LoadTomlConfig()
	testConfig.BaseDir = "testdata"

	tests := []struct {
		name   string
		fields fields
		want   []*models.MemoArchive
	}{
		{
			name: "Test MemoArchivesFromIndexChecked",
			fields: fields{
				config: testConfig,
				gmw:    markdown.NewGoldmarkWrapper(),
			},
			want: []*models.MemoArchive{
				{Text: "another memoarchive", Destination: "anywhere", Checked: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MemoArchiveRepo{
				config: tt.fields.config,
			}
			if got := repo.MemoArchivesFromIndexChecked(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoArchiveRepo.MemoArchivesFromIndexChecked() = %v, want %v", got, tt.want)
			}
		})
	}
}
