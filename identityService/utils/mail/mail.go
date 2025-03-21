package mail

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joy095/identity/config/db"
	"github.com/joy095/identity/models"

	"github.com/joy095/identity/logger"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	mail "github.com/xhit/go-simple-mail/v2"
	"golang.org/x/crypto/argon2"
)

var smtpClient *mail.SMTPClient
var redisClient *redis.Client
var ctx = context.Background()
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Failed to load .env file")
	}

	// Initialize Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Println("Connected to Redis")
			return nil
		},
	})
}

// Create a new SMTP client connection
func newSMTPClient() (*mail.SMTPClient, error) {
	server := mail.NewSMTPClient()
	server.Host = os.Getenv("SMTP_HOST")
	server.Port, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	server.Username = os.Getenv("SMTP_USERNAME")
	server.Password = os.Getenv("SMTP_PASSWORD")
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	return server.Connect()
}

// Generate a secure OTP using crypto/rand
func GenerateSecureOTP() string {
	const otpChars = "0123456789"
	bytes := make([]byte, 6)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Println("Error generating secure OTP:", err)
		return "000000"
	}
	for i := range bytes {
		bytes[i] = otpChars[bytes[i]%byte(len(otpChars))]
	}
	return string(bytes)
}

// Hash OTP using Argon2 for security
func hashOTP(otp string) string {
	salt := []byte("some_random_salt")
	hashed := argon2.IDKey([]byte(otp), salt, 1, 64*1024, 4, 32)
	return fmt.Sprintf("%x", hashed)
}

// Store OTP hash in Redis with expiration
func storeOTP(email, otp string) error {
	hashedOTP := hashOTP(otp)
	return redisClient.Set(context.Background(), "otp:"+email, hashedOTP, 10*time.Minute).Err()
}

func SendOTP(emailAddress, otp string) error {
	logger.InfoLogger.Info("SendOTP called on mail")

	var user models.User
	query := `SELECT id, email FROM users WHERE email = $1`

	// query := `SELECT id FROM users WHERE id = $1`
	err := db.DB.QueryRow(context.Background(), query, emailAddress).Scan(&user.ID, &user.Email)
	if err != nil {
		return err
	}

	// Store OTP before sending email
	if err := storeOTP(emailAddress, otp); err != nil {
		return err
	}

	tmpl, err := template.ParseFiles("otp_template.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer
	data := struct {
		OTP  string
		Year int
	}{
		OTP:  otp,
		Year: time.Now().Year(),
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	// Create a new SMTP client for each email
	smtpClient, err := newSMTPClient()
	if err != nil {
		logger.ErrorLogger.Errorf("failed to connect to SMTP server: %v", err)
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer smtpClient.Close()

	email := mail.NewMSG()
	email.SetFrom(os.Getenv("FROM_EMAIL")).
		AddTo(user.Email).
		SetSubject("Your OTP Code").
		SetBody(mail.TextHTML, body.String())

	logger.InfoLogger.Info("Sending OTP email to: ", user.Email)

	return email.Send(smtpClient)
}

// Generate JWT token
func generateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(10 * time.Minute).Unix(), // Token expires in 10 minutes
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Request OTP API
func RequestOTP(c *gin.Context) {
	logger.InfoLogger.Info("RequestOTP called on mail")

	var request struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.ErrorLogger.Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if request.Email == "" {
		logger.ErrorLogger.Error("Email is required")

		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	// Check if email exists in database
	var count int
	err := db.DB.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE email = $1", request.Email).Scan(&count)
	if err != nil {
		logger.ErrorLogger.Error("Failed to process request")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	if count == 0 {
		logger.InfoLogger.Info("If the email exists, an OTP has been sent")
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, an OTP has been sent"})
		return
	}

	otp := GenerateSecureOTP()
	err = storeOTP(request.Email, otp)
	if err != nil {
		logger.ErrorLogger.Error("Failed to store OTP")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store OTP"})
		return
	}

	err = SendOTP(request.Email, otp)
	if err != nil {
		logger.ErrorLogger.Error("Failed to send OTP")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	logger.InfoLogger.Info("OTP send successfully")

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// Verify OTP and return JWT token
func VerifyOTP(c *gin.Context) {
	logger.InfoLogger.Info("VerifyOTP called on mail")

	var request struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.ErrorLogger.Error("Invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if request.Email == "" || request.OTP == "" {
		logger.ErrorLogger.Error("Email and OTP are required")

		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and OTP are required"})
		return
	}

	// Retrieve OTP hash from Redis
	storedHash, err := redisClient.Get(ctx, "otp:"+request.Email).Result()
	if err != nil {
		logger.ErrorLogger.Error("OTP expired or not found")

		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired or not found"})
		return
	}

	// Verify OTP
	if hashOTP(request.OTP) != storedHash {
		logger.ErrorLogger.Error("Incorrect OTP")

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect OTP"})
		return
	}

	// Generate JWT token
	token, err := generateJWT(request.Email)
	if err != nil {
		logger.ErrorLogger.Error("Failed to generate JWT token")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	// Store JWT in Redis for 10 minutes
	redisClient.Set(ctx, "jwt:"+request.Email, token, 10*time.Minute)

	// Delete OTP from Redis
	redisClient.Del(ctx, "otp:"+request.Email)

	logger.InfoLogger.Info("OTP verified successfully")

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})

	// Update user's email verification status
	_, err = db.DB.Exec(context.Background(), "UPDATE users SET is_verified_email = true WHERE email = $1", request.Email)
	if err != nil {
		logger.ErrorLogger.Error("Failed to update email verification status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email verification status"})
		return
	}

}
