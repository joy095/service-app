package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joy095/message-service/db"
	"github.com/joy095/message-service/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	middleware "github.com/joy095/message-service/middlewares/cors"
)

var clients = make(map[*websocket.Conn]struct{})

type Message struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func init() {
	logger.InitLoggers()

	db.Connect()
	// db.RunMigrations("db/schema.sql")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	router := gin.Default()

	router.Use(middleware.CorsMiddleware())

	router.GET("/ws", serveWs)

	log.Println("Server starting on: " + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		allowed := os.Getenv("ALLOWED_ORIGINS")
// 		origins := strings.Split(allowed, ",")
// 		origin := r.Header.Get("Origin")

// 		for _, o := range origins {
// 			if strings.TrimSpace(o) == origin {
// 				return true
// 			}
// 		}
// 		return false
// 	},
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // ⚠️ Accept all origins — only for dev/testing!
	},
}

func serveWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	clients[conn] = struct{}{}
	log.Println("Client connected")

	go handleClient(conn)
}

func handleClient(conn *websocket.Conn) {
	defer func() {
		delete(clients, conn)
		conn.Close()
		log.Println("Client disconnected")
	}()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		log.Printf("Received message: %+v", msg)
		broadcast(msg)
	}
}

func broadcast(msg Message) {
	for conn := range clients {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Broadcast error: %v", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}
