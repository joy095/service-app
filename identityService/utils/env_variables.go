package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func GetJWTSecret() []byte {

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		fmt.Println("WARNING: JWT_SECRET environment variable not set.")
		return []byte("default-insecure-secret-only-for-development")
	}
	return []byte(secret)
}

func GetJWTRefreshSecret() []byte {

	secret := os.Getenv("JWT_SECRET_REFRESH")
	if secret == "" {
		fmt.Println("WARNING: JWT_SECRET_REFRESH environment variable not set.")
		return []byte("default-insecure--refresh-secret-only-for-development")
	}
	return []byte(secret)
}

func GetWordFilterService() []byte {

	secret := os.Getenv("WORD_FILTER_SERVICE_URL")
	if secret == "" {
		fmt.Println("WARNING: WORD_FILTER_SERVICE_URL environment variable not set.")
		return []byte("http://localhost:8082")
	}
	return []byte(secret)
}
