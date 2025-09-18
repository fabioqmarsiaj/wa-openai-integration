package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var verifyToken = os.Getenv("WEBHOOK_VERIFY_TOKEN")

func main() {
	router := gin.Default()

	router.GET("/webhooks", handleVerification)
	router.POST("/webhooks", handleMessageNotification)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting webhook server on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

// handleVerification processes the GET request for webhook setup
func handleVerification(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == verifyToken {
		log.Println("Webhook verified successfully!")
		c.Data(200, "text/plain", []byte(challenge))
	} else {
		log.Println("Error: Invalid verification token.")
		c.JSON(403, gin.H{"error": "Forbidden"})
	}
}

// handleMessageNotification processes the POST request with message data
func handleMessageNotification(c *gin.Context) {
	// Gin automatically parses the JSON body for us
	var payload interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("Error binding JSON: %v\n", err)
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// For demonstration, print the received payload
	log.Printf("Received webhook payload: %+v\n", payload)

	// Acknowledge receipt by sending a 200 OK
	c.JSON(200, gin.H{"status": "ok"})

	// Here you would process the message data
	// e.g., check the message type, sender, and content
	// and send a reply via the WhatsApp Cloud API.
}
