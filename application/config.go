package application

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hirotoni/memo/configs"
)

func (app *App) EditConfig() {
	configFile, err := configs.ConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("vim", configFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func (app *App) ShowConfig() {
	tomlConfig := configs.LoadTomlConfig()
	fmt.Printf("%+v", tomlConfig)
}
