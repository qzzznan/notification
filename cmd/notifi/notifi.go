package main

import (
	"log"
	"notification/config"
	"notification/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}
	app.Run(cfg)
}
