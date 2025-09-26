// main.go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, relying on system env variables")
	}

	apiKey := os.Getenv("TICKETMASTER_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ TICKETMASTER_API_KEY not set in environment")
	}
	log.Printf("🔑 Using Ticketmaster API Key: %s", apiKey)

	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Serve static files
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// API routes
	r.POST("/ingest/:marketplace", ingestHandler(apiKey))
	r.GET("/events/:marketplace", getEventsHandler)

	log.Println("🚀 Server listening on :8080")
	log.Fatal(r.Run(":8080"))
}
