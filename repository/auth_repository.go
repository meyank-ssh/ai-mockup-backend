package repository

import (
	"errors"

	"ai-mockup/internal/models"
	"ai-mockup/internal/service"
)

// CreateUser creates a new user in the database
func CreateUser(name, email, avatar, provider string) (*models.User, error) {
	user := models.User{
		Email:    email,
		Name:     name,
		Avatar:   avatar,
		Provider: provider,
	}
	if err := service.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail fetches a user by email
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := service.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// GetUserByID fetches a user by ID
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := service.DB.First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
