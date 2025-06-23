package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey = []byte("diddy_dickler") // Replace with your actual secret key

func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"token_type": "auth",
		"user_id":    userID,
		"exp":        time.Now().Add(24 * time.Hour * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GenerateTokenPayment(merchantID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"token_type":  "payment_only",
		"merchant_id": merchantID.String(),
		"exp":         time.Now().Add(24 * time.Hour * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return secretKey, nil
	})
}
