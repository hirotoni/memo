package repos

import (
	"reflect"
	"testing"

	"github.com/hirotoni/memo/config"
	"github.com/hirotoni/memo/markdown"
	"github.com/hirotoni/memo/models"
)

func TestTipRepo_TipsFromIndex(t *testing.T) {
	type fields struct {
		config *config.TomlConfig
		gmw    *markdown.GoldmarkWrapper
	}

	testConfig := config.LoadTomlConfig()
	testConfig.BaseDir = "testdata"

	tests := []struct {
		name   string
		fields fields
		want   []*models.Tip
	}{
		{
			name: "Test TipsFromIndex",
			fields: fields{
				config: testConfig,
				gmw:    markdown.NewGoldmarkWrapper(),
			},
			want: []*models.Tip{
				{Text: "some tip", Destination: "somewhere", Checked: false},
				{Text: "another tip", Destination: "anywhere", Checked: true},
				{Text: "yet another tip", Destination: "everywhere", Checked: false},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TipRepo{
				config: tt.fields.config,
				gmw:    tt.fields.gmw,
			}
			if got := repo.TipsFromIndex(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TipRepo.TipsFromIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTipRepo_TipsFromIndexChecked(t *testing.T) {
	type fields struct {
		config *config.TomlConfig
		gmw    *markdown.GoldmarkWrapper
	}

	testConfig := config.LoadTomlConfig()
	testConfig.BaseDir = "testdata"

	tests := []struct {
		name   string
		fields fields
		want   []*models.Tip
	}{
		{
			name: "Test TipsFromIndexChecked",
			fields: fields{
				config: testConfig,
				gmw:    markdown.NewGoldmarkWrapper(),
			},
			want: []*models.Tip{
				{Text: "another tip", Destination: "anywhere", Checked: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TipRepo{
				config: tt.fields.config,
				gmw:    tt.fields.gmw,
			}
			if got := repo.TipsFromIndexChecked(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TipRepo.TipsFromIndexChecked() = %v, want %v", got, tt.want)
			}
		})
	}
}
