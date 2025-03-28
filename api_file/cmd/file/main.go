package main

import (
	"log"

	"github.com/reversersed/LitGO-backend/tree/main/api_file/internal/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatalf("error while creating app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("error while starting app: %v", err)
	}
}
