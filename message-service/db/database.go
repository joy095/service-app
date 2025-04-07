package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joy095/message-service/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")

	if _, err := pgxpool.ParseConfig(dsn); err != nil {
		logger.ErrorLogger.Error("Unable to parse database URL:", err)
		fmt.Println("Unable to parse database URL:", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.ErrorLogger.Error("Unable to connect to database:", err)

		fmt.Println("Unable to connect to database:", err)
		os.Exit(1)
	}

	DB = pool
	logger.InfoLogger.Info("Connected to PostgreSQL!")

	fmt.Println("Connected to PostgreSQL!")
}

func RunMigrations(schemaPath string) {
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		logger.ErrorLogger.Error("Failed to read schema file:", err)
		fmt.Println("Failed to read schema file:", err)
		os.Exit(1)
	}

	queries := strings.Split(string(content), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		_, err := DB.Exec(context.Background(), query)
		if err != nil {
			logger.ErrorLogger.Error("Migration failed on query:", query, "Error:", err)
			fmt.Println("Migration failed on query:", query, "\nError:", err)
			os.Exit(1)
		}
	}

	logger.InfoLogger.Info("Database migration completed successfully!")
	fmt.Println("Database migration completed successfully!")
}
