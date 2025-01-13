package repos

import (
	"testing"

	"github.com/hirotoni/memo/config"
	"github.com/hirotoni/memo/markdown"
	"github.com/stretchr/testify/assert"
)

func TestDailymemoRepo_Entry(t *testing.T) {
	assert := assert.New(t)

	config := config.NewTomlConfig("testdata", 7, markdown.NewGoldmarkWrapper())
	repo := NewDailymemoRepo(config)
	tests := []struct {
		name string
	}{
		{name: "Test Entry"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := repo.Entry("testdata/dailymemo/2024-12-30-Mon.md")
			assert.Equal("testdata/dailymemo/2024-12-30-Mon.md", e.Filepath)
			assert.Equal("2024-12-30-Mon.md", e.BaseName)
			assert.Equal("2024-12-30", e.Date.Format("2006-01-02"))
			assert.Equal("testdata/dailymemo/2024-12-30-Mon.md", e.Filepath)
		})
	}

}

func TestDailymemoRepo_Entries(t *testing.T) {
	assert := assert.New(t)

	config := config.NewTomlConfig("testdata", 7, markdown.NewGoldmarkWrapper())
	repo := NewDailymemoRepo(config)
	tests := []struct {
		name string
	}{
		{name: "Test Entries"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := repo.Entries()
			assert.Len(e, 3)
		})
	}
}
