package main

import (
	"mathbattle/config"
	"mathbattle/infrastructure"
	"mathbattle/interfaces/server"
)

func main() {
	server.Start(infrastructure.NewServerContainer(config.LoadConfig("config.yaml")))
}
