package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Host      string
	DBName    string
	SSlmode   string
	Password  string
	User      string
	PORT      string
	JWTSecret string
}

func Load() (*Config, error) {
	var err error = godotenv.Load()
	if err != nil {
		log.Println("Watning: .env file not found, using enviromental variables")
	}

	var config *Config = &Config{
		Port:      os.Getenv("DB_PORT"),
		Host:      os.Getenv("DB_HOST"),
		DBName:    os.Getenv("DB_NAME"),
		Password:  os.Getenv("DB_PASSWORD"),
		SSlmode:   os.Getenv("DB_SSLMODE"),
		User:      os.Getenv("DB_USER"),
		PORT:      os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	return config, nil

}
