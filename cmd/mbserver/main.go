package main

import (
	"mathbattle/config"
	"mathbattle/infrastructure"
	"mathbattle/interfaces/server"
)

func main() {
	server.Start(infrastructure.NewContainer(config.LoadConfig("config.yaml")))
}
