package main

import (
	"log"
	"net/http"

	"gangsaur.com/billing-exercise/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	apiServer := server.NewApiServer()

	log.Printf("Listening at %v...", apiServer.GetServerAddress())
	if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server err: %s", err.Error())
	}
	log.Printf("Exiting...")
}
