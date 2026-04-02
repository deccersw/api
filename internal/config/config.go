package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	Host            string
	DBName          string
	SSlmode         string
	Password        string
	User            string
	PORT            string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func parseDuration(raw string, fallback time.Duration) time.Duration {
	if raw == "" {
		return fallback
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		log.Printf("invalid duration %q, using default %s", raw, fallback)
		return fallback
	}
	return d
}

func Load() (*Config, error) {
	var err error = godotenv.Load()
	if err != nil {
		log.Println("Watning: .env file not found, using enviromental variables")
	}

	var config *Config = &Config{
		Port:            os.Getenv("DB_PORT"),
		Host:            os.Getenv("DB_HOST"),
		DBName:          os.Getenv("DB_NAME"),
		Password:        os.Getenv("DB_PASSWORD"),
		SSlmode:         os.Getenv("DB_SSLMODE"),
		User:            os.Getenv("DB_USER"),
		PORT:            os.Getenv("PORT"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		AccessTokenTTL:  parseDuration(os.Getenv("ACCESS_TOKEN_TTL"), 15*time.Minute),
		RefreshTokenTTL: parseDuration(os.Getenv("REFRESH_TOKEN_TTL"), 24*time.Hour),
	}

	return config, nil

}
