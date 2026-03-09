package main

import (
	"log"

	"coffee-consortium/backend/internal/transport/httpapi"
)

func main() {
	app, err := httpapi.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

