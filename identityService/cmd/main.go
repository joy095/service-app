package main

import (
	"fmt"
	"identity/config/db"
	"identity/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	db.Connect()

	err := db.CreateUsersTable()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database migration completed!")
}

func main() {

	port := os.Getenv("PORT")

	r := gin.Default()

	routes.RegisterRoutes(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":" + port)
}
