package mail

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"html/template"
	"identity/config/db"
	"identity/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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
var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // Ensure this is set in your .env file

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

	// Setup SMTP server
	server := mail.NewSMTPClient()
	server.Host = os.Getenv("SMTP_HOST")
	server.Port, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	server.Username = os.Getenv("SMTP_USERNAME")
	server.Password = os.Getenv("SMTP_PASSWORD")
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err = server.Connect()
	if err != nil {
		log.Fatal("Failed to connect to SMTP server:", err)
	}
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

	email := mail.NewMSG()
	email.SetFrom(os.Getenv("FROM_EMAIL")).
		AddTo(user.Email).
		SetSubject("Your OTP Code").
		SetBody(mail.TextHTML, body.String())

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
	var request struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if request.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	// Check if email exists in database
	var count int
	err := db.DB.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE email = $1", request.Email).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, an OTP has been sent"})
		return
	}

	otp := GenerateSecureOTP()
	err = storeOTP(request.Email, otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store OTP"})
		return
	}

	err = SendOTP(request.Email, otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// Verify OTP and return JWT token
func VerifyOTP(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if request.Email == "" || request.OTP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and OTP are required"})
		return
	}

	// Retrieve OTP hash from Redis
	storedHash, err := redisClient.Get(ctx, "otp:"+request.Email).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired or not found"})
		return
	}

	// Verify OTP
	if hashOTP(request.OTP) != storedHash {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect OTP"})
		return
	}

	// Generate JWT token
	token, err := generateJWT(request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	// Store JWT in Redis for 10 minutes
	redisClient.Set(ctx, "jwt:"+request.Email, token, 10*time.Minute)

	// Delete OTP from Redis
	redisClient.Del(ctx, "otp:"+request.Email)

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully", "token": token})

	// Update user's email verification status
	_, err = db.DB.Exec(context.Background(), "UPDATE users SET is_verified_email = true WHERE email = $1", request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email verification status"})
		return
	}

}
