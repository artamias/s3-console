package main

import (
	"context"
	"fmt"
	"log"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	mc         *minio.Client
	bucketName string
)

// func env(key, fallback string) string {
// 	if v := os.Getenv(key); v != "" {
// 		return v
// 	}
// 	return fallback
// }

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("[config] .env not found, using environment variables")
	} else {
		log.Println("[config] .env loaded")
	}

	switch os.Getenv("APP_ENV") {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

func initMinio() {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"
	region := os.Getenv("MINIO_REGION")
	bucketName = os.Getenv("MINIO_BUCKET")
	var err error
	mc, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		log.Fatalf("[minio] init error: %v", err)
	}

	ctx := context.Background()
	exists, err := mc.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("[minio] BucketExists error: %v", err)
	}
	if !exists {
		if err = mc.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			log.Fatalf("[minio] MakeBucket error: %v", err)
		}
		log.Printf("[minio] bucket '%s' created", bucketName)
	} else {
		log.Printf("[minio] connected to bucket '%s' at %s", bucketName, endpoint)
	}

}

func presignExpiry() time.Duration {
	secs, err := strconv.Atoi(os.Getenv("MINIO_PRESIGN_EXPIRY"))
	if err != nil || secs <= 0 {
		return time.Hour
	}
	return time.Duration(secs) * time.Second
}

func corsOrigins() []string {
	// MINIO_CORS_ALLOW_ORIGINS is retained for existing deployments.
	raw := os.Getenv("CORS_ALLOW_ORIGINS")
	if raw == "*" {
		return []string{"*"}
	}
	origins := []string{}
	for _, o := range strings.Split(raw, ",") {
		if s := strings.TrimSpace(o); s != "" {
			origins = append(origins, s)
		}
	}
	return origins
}

// FileInfo adalah response shape untuk satu object.
type FileInfo struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag"`
}

// GET /api/v1/files
func listFiles(c *gin.Context) {
	ctx := context.Background()
	objects := mc.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true})

	result := make([]FileInfo, 0)
	for obj := range objects {
		if obj.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": obj.Err.Error()})
			return
		}
		result = append(result, FileInfo{
			Name:         obj.Key,
			Size:         obj.Size,
			LastModified: obj.LastModified,
			ETag:         obj.ETag,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": result, "total": len(result)})
}

// POST /api/v1/files/presigned-upload
// Mengembalikan presigned PUT URL — client upload langsung ke MinIO tanpa lewat backend.
func presignedUploadURL(c *gin.Context) {
	var req struct {
		Filename    string `json:"filename"     binding:"required"`
		ContentType string `json:"content_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.ContentType == "" {
		req.ContentType = "application/octet-stream"
	}

	ctx := context.Background()
	expiry := presignExpiry()

	u, err := mc.PresignedPutObject(ctx, bucketName, req.Filename, expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upload_url":   u.String(),
		"method":       "PUT",
		"filename":     req.Filename,
		"content_type": req.ContentType,
		"expires_in":   expiry.String(),
	})
}

// POST /api/v1/files  — fallback upload via backend (untuk non-browser client)
func uploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field 'file' is required"})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	ctx := context.Background()
	info, err := mc.PutObject(
		ctx, bucketName, header.Filename, file, header.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("'%s' uploaded successfully", info.Key),
		"data": gin.H{
			"name": info.Key,
			"size": info.Size,
			"etag": info.ETag,
		},
	})
}

// GET /api/v1/files/:name/presigned returns a direct MinIO URL for view or download.
func presignedURL(c *gin.Context) {
	name := c.Param("name")
	ctx := context.Background()

	if _, err := mc.StatObject(ctx, bucketName, name, minio.StatObjectOptions{}); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	expiry := presignExpiry()
	responseOverrides := url.Values{}
	if c.Query("download") == "true" {
		filename := path.Base(name)
		if filename == "." || filename == "/" || filename == "" {
			filename = "download"
		}
		responseOverrides.Set("response-content-disposition", mime.FormatMediaType("attachment", map[string]string{"filename": filename}))
		responseOverrides.Set("response-content-type", "application/octet-stream")
	}
	u, err := mc.PresignedGetObject(ctx, bucketName, name, expiry, responseOverrides)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":        u.String(),
		"filename":   name,
		"expires_in": expiry.String(),
	})
}

// GET /api/v1/files/:name/download streams an object with attachment headers.
func downloadFile(c *gin.Context) {
	name := c.Param("name")
	ctx := context.Background()

	obj, err := mc.GetObject(ctx, bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer obj.Close()

	stat, err := obj.Stat()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	filename := path.Base(name)
	if filename == "." || filename == "/" || filename == "" {
		filename = "download"
	}
	c.Header("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": filename}))
	c.Header("X-Content-Type-Options", "nosniff")
	c.DataFromReader(http.StatusOK, stat.Size, "application/octet-stream", obj, nil)
}

// DELETE /api/v1/files/:name
func deleteFile(c *gin.Context) {
	name := c.Param("name")
	ctx := context.Background()

	if _, err := mc.StatObject(ctx, bucketName, name, minio.StatObjectOptions{}); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	if err := mc.RemoveObject(ctx, bucketName, name, minio.RemoveObjectOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("'%s' deleted successfully", name)})
}

// GET /api/v1/info
func storageInfo(c *gin.Context) {
	ctx := context.Background()

	objects := mc.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true})
	count := 0
	var totalSize int64
	for obj := range objects {
		if obj.Err == nil {
			count++
			totalSize += obj.Size
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"bucket":         bucketName,
		"endpoint":       os.Getenv("MINIO_ENDPOINT"),
		"files":          count,
		"total_size":     totalSize,
		"presign_expiry": presignExpiry().String(),
	})
}

func main() {
	loadEnv()
	initMinio()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins(),
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Disposition", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "env": os.Getenv("APP_ENV")})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/info", storageInfo)

		// File list
		v1.GET("/files", listFiles)

		// Upload: presigned
		v1.POST("/files/presigned-upload", presignedUploadURL)
		v1.POST("/files", uploadFile)

		// Access: presigned
		v1.GET("/files/:name/presigned", presignedURL)
		v1.GET("/files/:name/download", downloadFile)

		// Delete
		v1.DELETE("/files/:name", deleteFile)
	}

	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")
	address := net.JoinHostPort(host, port)
	log.Printf("[server] listening on %s (env: %s)", address, os.Getenv("APP_ENV"))
	if err := r.Run(address); err != nil {
		log.Fatal(err)
	}
}
