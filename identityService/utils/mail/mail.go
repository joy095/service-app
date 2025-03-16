package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	mail "github.com/xhit/go-simple-mail/v2"
)

var smtpClient *mail.SMTPClient
var otpStore = make(map[string]string) // Stores OTP temporarily
var otpMutex sync.Mutex                // Ensures thread safety

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Failed to load .env file")
	}

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

// Generate a random 6-digit OTP
func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// Store OTP in memory for 5 minutes
func storeOTP(email, otp string) {
	otpMutex.Lock()
	defer otpMutex.Unlock()
	otpStore[email] = otp

	// Automatically delete OTP after 5 minutes
	go func() {
		time.Sleep(5 * time.Minute)
		otpMutex.Lock()
		delete(otpStore, email)
		otpMutex.Unlock()
	}()
}

// Send OTP via email
func sendOTP(emailAddress, otp string) error {
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

	if err := tmpl.ExecuteTemplate(&body, "otp_template.html", data); err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(os.Getenv("FROM_EMAIL")).
		AddTo(emailAddress).
		SetSubject("Your OTP Code").
		SetBody(mail.TextHTML, body.String())

	return email.Send(smtpClient)
}

// Request OTP API
func RequestOTP(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}

	// Parse JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate email
	if request.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	// Generate OTP
	otp := generateOTP()
	storeOTP(request.Email, otp) // Store OTP temporarily

	// Send OTP email
	err := sendOTP(request.Email, otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP: " + err.Error()})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// Verify OTP API
func VerifyOTP(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	// Parse JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate email and OTP
	if request.Email == "" || request.OTP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and OTP are required"})
		return
	}

	// Check OTP in memory
	otpMutex.Lock()
	storedOTP, exists := otpStore[request.Email]
	otpMutex.Unlock()

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired or not found"})
		return
	}

	if storedOTP != request.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect OTP"})
		return
	}

	// OTP verified successfully
	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})

	// Remove OTP after successful verification
	otpMutex.Lock()
	delete(otpStore, request.Email)
	otpMutex.Unlock()
}
