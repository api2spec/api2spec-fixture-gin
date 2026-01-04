package main

import (
	"log"
	"os"

	"github.com/api2spec/api2spec-fixture-gin/internal/router"
)

func main() {
	r := router.Setup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Tea API running at http://localhost:%s", port)
	log.Printf("TIF signature: http://localhost:%s/brew", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
