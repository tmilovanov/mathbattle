package main

import (
	"os"

	"mathbattle/config"
	"mathbattle/infrastructure"
	"mathbattle/interfaces/server"
)

func main() {
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	server.Start(infrastructure.NewServerContainer(config.LoadConfig(configPath)))
}
