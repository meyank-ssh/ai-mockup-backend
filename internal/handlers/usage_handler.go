package handlers

import (
	"net/http"

	"ai-mockup/repository"

	"github.com/gin-gonic/gin"
)

type UsageResponse struct {
	ImageCount int `json:"image_count"`
}

func GetUsageHandler(c *gin.Context) {
	// Get the user ID from the context (set by AuthMiddleware)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	// Use repository to get or create the usage record
	usage, err := repository.GetOrCreateUsage(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get or create usage record"})
		return
	}

	// Return the usage data
	c.JSON(http.StatusOK, UsageResponse{ImageCount: usage.ImageCount})
}
