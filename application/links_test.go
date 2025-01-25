package application

import (
	"testing"

	"github.com/hirotoni/memo/configs"
	"github.com/hirotoni/memo/markdown"
)

func TestLinks(t *testing.T) {
	app := NewApp()
	app.WithCustomConfig(
		*configs.NewTomlConfig(
			"testdata",
			10,
			markdown.NewGoldmarkWrapper(),
		),
	)
	app.Links()
}
