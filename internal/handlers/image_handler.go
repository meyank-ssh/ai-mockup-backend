package handlers

import (
	"bytes"
	"net/http"
	"time"

	"ai-mockup/internal/models"
	"ai-mockup/repository"

	"github.com/gin-gonic/gin"
)

type GenerateResponse struct {
	ImageURL string `json:"image_url"`
}

func GenerateImageHandler(c *gin.Context) {
	// Parse the prompt from the form-data
	prompt := c.PostForm("prompt")
	if prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt is required"})
		return
	}

	// Parse the file from the form-data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Open the file and read its content
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer fileContent.Close()

	// Read the file content into a byte array
	var fileBuffer bytes.Buffer
	if _, err := fileBuffer.ReadFrom(fileContent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}

	// Call OpenAI API (or other AI service) to generate the image
	imageURL, err := callOpenAIWithFile(fileBuffer.Bytes(), prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate image"})
		return
	}

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

	// Save the image metadata to the database
	image := models.Image{
		UserID:    userID,
		Prompt:    prompt,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
	}
	if err := repository.SaveImage(image); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Update usage count
	if err := repository.IncrementUsage(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update usage"})
		return
	}

	// Return the generated image URL
	c.JSON(http.StatusOK, GenerateResponse{ImageURL: imageURL})
}

// Mock function to call OpenAI API with file content and prompt (replace with actual implementation)
func callOpenAIWithFile(fileContent []byte, prompt string) (string, error) {
	// Simulate OpenAI API call
	return "https://example.com/generated-image.png", nil
}
