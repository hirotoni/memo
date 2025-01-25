package repos

import (
	"testing"

	"github.com/hirotoni/memo/configs"
	"github.com/hirotoni/memo/markdown"
	"github.com/stretchr/testify/assert"
)

func TestDailymemoRepo_Entry(t *testing.T) {
	assert := assert.New(t)

	config := configs.NewTomlConfig("testdata", 7, markdown.NewGoldmarkWrapper())
	repo := NewDailymemoRepo(config)
	tests := []struct {
		name string
	}{
		{name: "Test Entry"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := repo.Entry("testdata/dailymemo/2024-12-30-Mon.md")
			assert.Nil(err)
			assert.Equal("testdata/dailymemo/2024-12-30-Mon.md", e.Filepath)
			assert.Equal("2024-12-30-Mon.md", e.BaseName)
			assert.Equal("2024-12-30", e.Date.Format("2006-01-02"))
			assert.Equal("testdata/dailymemo/2024-12-30-Mon.md", e.Filepath)
		})
	}

}

func TestDailymemoRepo_Entries(t *testing.T) {
	assert := assert.New(t)

	config := configs.NewTomlConfig("testdata", 7, markdown.NewGoldmarkWrapper())
	repo := NewDailymemoRepo(config)
	tests := []struct {
		name string
	}{
		{name: "Test Entries"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := repo.Entries()
			assert.Nil(err)
			assert.Len(e, 3)
		})
	}
}

func TestDailymemoRepo_FindByDate(t *testing.T) {
	assert := assert.New(t)

	config := configs.NewTomlConfig("testdata", 7, markdown.NewGoldmarkWrapper())
	repo := NewDailymemoRepo(config)
	tests := []struct {
		name string
	}{
		{name: "Test FindByDate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md, err := repo.FindByDate("2024-12-30-Mon")
			assert.Nil(err)
			assert.Equal("testdata/dailymemo/2024-12-30-Mon.md", md.Filepath)
		})
	}
}

func TestDailymemoRepo_MemosFromDailymemo(t *testing.T) {
	// assert := assert.New(t)

	config := configs.NewTomlConfig("testdata", 7, markdown.NewGoldmarkWrapper())
	repo := NewDailymemoRepo(config)
	tests := []struct {
		name string
	}{
		{name: "Test FindByDate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dm, _ := repo.Entry("testdata/dailymemo/2025-01-01-Wed.md")
			mm := repo.MemosFromDailymemo(dm)
			for _, m := range mm {
				m.Print()
			}
		})
	}
}
