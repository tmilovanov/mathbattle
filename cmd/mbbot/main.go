package main

import (
	"mathbattle/config"
	"mathbattle/infrastructure"
	"mathbattle/interfaces/bot"
)

func main() {
	bot.Start(infrastructure.NewContainer(config.LoadConfig("config.yaml")))
}
