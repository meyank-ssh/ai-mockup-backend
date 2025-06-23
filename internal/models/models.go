package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"unique;not null"`
	Name      string    `gorm:"size:255"`
	Avatar    string    `gorm:"type:text"`
	Provider  string    `gorm:"size:50;not null"` // e.g., "google"
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Subscription struct {
	ID             uint      `gorm:"primaryKey"`
	UserID         uint      `gorm:"not null"`
	DodoCustomerID string    `gorm:"size:255;not null"` // Customer ID from Dodo Payments
	PlanName       string    `gorm:"size:50;not null"`  // e.g., "Free", "Standard", "Pro"
	Status         string    `gorm:"size:50;not null"`  // e.g., "active", "canceled"
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}

type Image struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Prompt    string    `gorm:"type:text;not null"`
	ImageURL  string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Usage struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null"`
	ImageCount int       `gorm:"default:0"` // Tracks the number of images generated
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
