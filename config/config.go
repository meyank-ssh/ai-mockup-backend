package configs

import (
	"log"
	"os"
	"strconv"
)

type GoogleAuth struct {
	RedirectURL  string
	ClientID     string
	ClientSecret string
}

type Config struct {
	GoogleAuth   GoogleAuth
	DatabaseURL  string
	JWTSecret    string
	DashboardURL string
}

var AppConfig Config

func parseFloatEnv(key string) float64 {
	val := os.Getenv(key)
	if val == "" {
		return 0
	}
	parsedVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Fatalf("Error parsing float value for environment variable %s: %v", key, err)
	}
	return parsedVal
}

func LoadConfig() {
	required := []string{
		"REDIRECT_URL",
		"GOOGLE_CLIENT_ID",
		"GOOGLE_CLIENT_SECRET",
		"DATABASE_URL",
		"JWT_SECRET",
		"DASHBOARD_URL",
	}

	missing := []string{}
	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		log.Fatalf("Missing required environment variables: %v", missing)
	}

	AppConfig = Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		DashboardURL: os.Getenv("DASHBOARD_URL"),
		GoogleAuth: GoogleAuth{
			RedirectURL:  os.Getenv("REDIRECT_URL"),
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		},
	}
}
