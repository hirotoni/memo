package repos

import (
	"testing"

	"github.com/hirotoni/memo/config"
	"github.com/hirotoni/memo/markdown"
	"github.com/stretchr/testify/assert"
)

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
