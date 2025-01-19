package usecases

import (
	"errors"
	"os"
	"testing"

	"github.com/hirotoni/memo/models"
	"github.com/stretchr/testify/assert"
)

var updateGolden = false

func TestGenerateTemplateString(t *testing.T) {
	type args struct {
		t models.Template
	}

	dailyMemoGolden, err := os.ReadFile("./testdata/templatedailymemo.md.golden")
	if errors.Is(err, os.ErrNotExist) {
		dailyMemoGolden = []byte{}
	}
	weeklyReportGolden, err := os.ReadFile("./testdata/templateweeklyreport.md.golden")
	if errors.Is(err, os.ErrNotExist) {
		weeklyReportGolden = []byte{}
	}
	memoArchivesGolden, err := os.ReadFile("./testdata/templatememoarchives.md.golden")
	if errors.Is(err, os.ErrNotExist) {
		memoArchivesGolden = []byte{}
	}
	memoArchivesIndexGolden, err := os.ReadFile("./testdata/templatememoarchivesindex.md.golden")
	if errors.Is(err, os.ErrNotExist) {
		memoArchivesIndexGolden = []byte{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "templatedailymemo",
			args: args{t: TemplateDailymemo},
			want: string(dailyMemoGolden),
		},
		{
			name: "templateweeklyreport",
			args: args{t: TemplateWeeklyReport},
			want: string(weeklyReportGolden),
		},
		{
			name: "templatememoarchives",
			args: args{t: TemplateMemoArchives},
			want: string(memoArchivesGolden),
		},
		{
			name: "templatememoarchivesindex",
			args: args{t: TemplateMemoArchivesIndex},
			want: string(memoArchivesIndexGolden),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			got := GenerateTemplateString(tt.args.t)
			assert.Equal(tt.want, got)

			if updateGolden {
				os.WriteFile("./testdata/"+tt.name+".md.golden", []byte(got), 0644)
			}
		})
	}
}
