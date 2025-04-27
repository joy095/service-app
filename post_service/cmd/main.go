package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Post represents the post data structure
type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	ImageURL  string    `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostRequest represents the incoming request for creating a post
type PostRequest struct {
	Title   string `form:"title" binding:"required"`
	Content string `form:"content" binding:"required"`
}

// App holds all the application dependencies
type App struct {
	DB     *pgxpool.Pool
	S3     *s3.Client
	Bucket string
	R2URL  string
}

func main() {
	godotenv.Load()

	// Initialize PostgreSQL connection
	dbpool, err := initPostgres()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Initialize Cloudflare R2 client
	s3Client, bucket, r2URL, err := initR2()
	if err != nil {
		log.Fatalf("Unable to initialize R2: %v", err)
	}

	// Initialize app
	app := &App{
		DB:     dbpool,
		S3:     s3Client,
		Bucket: bucket,
		R2URL:  r2URL,
	}

	// Setup database
	if err := setupDB(dbpool); err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Routes
	router.POST("/posts", app.createPost)
	router.GET("/posts", app.getPosts)
	router.GET("/posts/:id", app.getPost)

	// Start server
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Initialize PostgreSQL connection
func initPostgres() (*pgxpool.Pool, error) {
	connString := "postgres://postgres:admin123@localhost:5432/posts_db"

	// Use environment variable if available
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		connString = dbURL
	}

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}

// Initialize Cloudflare R2 client
func initR2() (*s3.Client, string, string, error) {
	// R2 credentials - in production, use environment variables
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("R2_ACCESS_KEY_SECRET")
	bucket := os.Getenv("R2_BUCKET")

	// Correct R2 endpoint format
	r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	if accountID == "" || accessKeyID == "" || accessKeySecret == "" || bucket == "" {
		return nil, "", "", fmt.Errorf("missing R2 credentials")
	}

	// Create the AWS configuration with custom endpoint
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, "", "", fmt.Errorf("unable to load SDK config: %w", err)
	}

	// Create S3 client with custom endpoint options
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(r2Endpoint)
		o.UsePathStyle = true
	})

	return client, bucket, r2Endpoint, nil
}

// Setup database tables
func setupDB(db *pgxpool.Pool) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS posts (
		id UUID PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		image_url TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(context.Background(), createTableSQL)
	return err
}

// Create a new post with image upload
func (app *App) createPost(c *gin.Context) {
	// Parse form data
	var req PostRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate post ID
	postID := uuid.New().String()

	// Process image if provided
	var imageURL string
	file, header, err := c.Request.FormFile("image")
	if err == nil {
		// Close the file when done
		defer file.Close()

		// Get file extension
		fileExt := filepath.Ext(header.Filename)

		// Generate unique image name with UUID
		uniqueImageName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)

		// Use a common "images" folder for all uploads
		objectKey := fmt.Sprintf("images/%s", uniqueImageName)

		// Upload to R2
		imageURL, err = app.uploadToR2(file, objectKey, header.Header.Get("Content-Type"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to upload image: " + err.Error(),
			})
			return
		}
	}

	// Insert post into database
	now := time.Now()
	post := Post{
		ID:        postID,
		Title:     req.Title,
		Content:   req.Content,
		ImageURL:  imageURL,
		CreatedAt: now,
		UpdatedAt: now,
	}

	insertSQL := `
	INSERT INTO posts (id, title, content, image_url, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, title, content, image_url, created_at, updated_at;`

	err = app.DB.QueryRow(
		context.Background(),
		insertSQL,
		post.ID,
		post.Title,
		post.Content,
		post.ImageURL,
		post.CreatedAt,
		post.UpdatedAt,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.ImageURL,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create post: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// Upload file to Cloudflare R2
func (app *App) uploadToR2(file io.Reader, key, contentType string) (string, error) {
	ctx := context.Background()

	// Read the entire file into memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Upload the file to R2
	_, err = app.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(app.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileBytes), // Use bytes.NewReader which is seekable
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate public URL
	imageURL := fmt.Sprintf("%s/%s/%s", app.R2URL, app.Bucket, key)
	return imageURL, nil
}

// Get all posts
func (app *App) getPosts(c *gin.Context) {
	var posts []Post

	rows, err := app.DB.Query(context.Background(), `
		SELECT id, title, content, image_url, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT 100
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.ImageURL,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan post"})
			return
		}
		posts = append(posts, post)
	}

	c.JSON(http.StatusOK, posts)
}

// Get a single post by ID
func (app *App) getPost(c *gin.Context) {
	id := c.Param("id")
	var post Post

	err := app.DB.QueryRow(context.Background(), `
		SELECT id, title, content, image_url, created_at, updated_at
		FROM posts
		WHERE id = $1
	`, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.ImageURL,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}
