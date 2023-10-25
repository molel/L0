package main

import (
	"L0/config"
	"L0/internal/app"
)

func main() {
	cfg := config.NewConfig()

	app.Run(cfg)
}
