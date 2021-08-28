package main

import (
	"log"

	initModule "gdrive-telegram-bot/src/init"
)

func main() {
	// Initialize
	err := initModule.Start()
	if err != nil {
		log.Fatal("failed to initialize app - ", err)
	}
}
