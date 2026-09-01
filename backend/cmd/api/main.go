package main

import (
	"log"

	"infrapilot/backend/internal/app"
)

func main() {
	if err := app.RunAPI(); err != nil {
		log.Fatal(err)
	}
}
