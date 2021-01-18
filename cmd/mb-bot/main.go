package main

import (
	"os"

	"mathbattle/config"
	"mathbattle/infrastructure"
	"mathbattle/interfaces/bot"
)

func main() {
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	bot.Start(infrastructure.NewBotContainer(config.LoadConfig(configPath)))
}
