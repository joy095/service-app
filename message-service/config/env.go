// internal/config/env.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := os.Getenv("GO_ENV")

	var envFile string
	if env == "development" {
		envFile = ".env.local"
	} else {
		envFile = ".env"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("No %s file found or failed to load: %v", envFile, err)
	} else {
		log.Printf("Loaded environment variables from %s", envFile)
	}
}
