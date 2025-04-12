package relations

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joy095/identity/config/db"
)

// RelationController handles user relationship operations
type RelationController struct{}

func NewRelationController() *RelationController {
	return &RelationController{}
}

func (r *RelationController) SendRequest(c *gin.Context) {
	var payload struct {
		ToUserID string `json:"addressee_id"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fromUserID := c.GetString("user_id")

	// Check for empty user IDs
	if fromUserID == "" || payload.ToUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID(s)"})
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(fromUserID); err != nil {
		log.Println("Invalid fromUserID:", fromUserID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender UUID format"})
		return
	}
	if _, err := uuid.Parse(payload.ToUserID); err != nil {
		log.Println("Invalid toUserID:", payload.ToUserID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid addressee UUID format"})
		return
	}

	// Prevent sending a request to oneself
	if fromUserID == payload.ToUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send request to yourself"})
		return
	}

	// Check if relationship already exists
	var exists bool
	err := db.DB.QueryRow(context.Background(), `
		SELECT EXISTS(
			SELECT 1 FROM user_connections 
			WHERE (requester_id = $1 AND addressee_id = $2) 
			OR (requester_id = $2 AND addressee_id = $1)
		)
	`, fromUserID, payload.ToUserID).Scan(&exists)

	if err != nil {
		log.Printf("Error checking existing relationship: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing relationship"})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Relationship already exists"})
		return
	}

	// Insert the connection request
	_, err = db.DB.Exec(context.Background(), `
		INSERT INTO user_connections (requester_id, addressee_id, status)
		VALUES ($1, $2, 'pending')
	`, fromUserID, payload.ToUserID)

	if err != nil {
		log.Printf("Error inserting connection request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connection request sent"})
}

func (r *RelationController) AcceptRequest(c *gin.Context) {
	var payload struct {
		FromUserID string `json:"requester_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	toUserID := c.GetString("user_id")
	_, err := db.DB.Exec(context.Background(), `
		UPDATE user_connections SET status = 'accepted'
		WHERE requester_id = $1 AND addressee_id = $2 AND status = 'pending'
	`, payload.FromUserID, toUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connection request accepted"})
}

func (r *RelationController) RejectRequest(c *gin.Context) {
	var payload struct {
		FromUserID string `json:"requester_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	toUserID := c.GetString("user_id")
	_, err := db.DB.Exec(context.Background(), `
		DELETE FROM user_connections
		WHERE requester_id = $1 AND addressee_id = $2 AND status = 'pending'
	`, payload.FromUserID, toUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connection request rejected"})
}

func (r *RelationController) ListPendingRequests(c *gin.Context) {
	userID := c.GetString("user_id")
	rows, err := db.DB.Query(context.Background(), `
		SELECT requester_id FROM user_connections WHERE addressee_id = $1 AND status = 'pending'
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch requests"})
		return
	}
	defer rows.Close()

	var pending []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		pending = append(pending, id)
	}
	c.JSON(http.StatusOK, gin.H{"pending_requests": pending})
}

func (r *RelationController) ListConnections(c *gin.Context) {
	userID := c.GetString("user_id")
	rows, err := db.DB.Query(context.Background(), `
		SELECT requester_id, addressee_id FROM user_connections 
		WHERE (requester_id = $1 OR addressee_id = $1) AND status = 'accepted'
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch connections"})
		return
	}
	defer rows.Close()

	connections := []string{}
	for rows.Next() {
		var fromID, toID string
		rows.Scan(&fromID, &toID)
		if fromID == userID {
			connections = append(connections, toID)
		} else {
			connections = append(connections, fromID)
		}
	}
	c.JSON(http.StatusOK, gin.H{"connections": connections})
}

func (r *RelationController) CheckConnectionStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	targetID := c.Param("user_id")

	var status string
	err := db.DB.QueryRow(context.Background(), `
		SELECT status FROM user_connections 
		WHERE (requester_id = $1 AND addressee_id = $2)
		   OR (requester_id = $2 AND addressee_id = $1)
	`, userID, targetID).Scan(&status)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{"status": "none"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}
