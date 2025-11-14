package main

import (
	"log"

	"avitotest/cmd/server/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
