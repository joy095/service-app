package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := os.Getenv("GO_ENV")

	if env == "development" {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Printf("No .env.local file found or failed to load: %v", err)
		} else {
			log.Println("Loaded environment variables from .env.local")
		}
	} else {
		log.Println("Production mode: using system environment variables only")
	}
}
