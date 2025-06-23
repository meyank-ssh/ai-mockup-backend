package repository

import (
	"ai-mockup/internal/models"
	"ai-mockup/internal/service"
)

// SaveImage saves the image metadata to the database
func SaveImage(image models.Image) error {
	return service.DB.Create(&image).Error
}

// IncrementUsage increments the user's image generation count
func IncrementUsage(userID uint) error {
	var usage models.Usage
	if err := service.DB.Where("user_id = ?", userID).First(&usage).Error; err != nil {
		// If no usage record exists, create one
		if err := service.DB.Create(&models.Usage{UserID: userID, ImageCount: 1}).Error; err != nil {
			return err
		}
		return nil
	}
	// Increment the image count
	usage.ImageCount++
	return service.DB.Save(&usage).Error
}

// GetOrCreateUsage fetches or creates a usage record for a user
func GetOrCreateUsage(userID uint) (*models.Usage, error) {
	var usage models.Usage
	if err := service.DB.Where("user_id = ?", userID).First(&usage).Error; err != nil {
		// If not found, create a new usage record
		usage = models.Usage{UserID: userID, ImageCount: 0}
		if err := service.DB.Create(&usage).Error; err != nil {
			return nil, err
		}
	}
	return &usage, nil
}
