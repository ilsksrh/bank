package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	ServerPort    string

	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env не найден, продолжаем с переменными окружения")
	}

	AppConfig = &Config{
		AccessSecret:  getEnv("ACCESS_SECRET", "access-secret"),
		RefreshSecret: getEnv("REFRESH_SECRET", "refresh-secret"),
		AccessTTL:     parseDuration("ACCESS_TTL", 5*time.Minute),
		RefreshTTL:    parseDuration("REFRESH_TTL", 30*time.Minute),
		ServerPort:    getEnv("SERVER_PORT", ":8080"),

		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "jusan"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	dur, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return dur
}

func (c *Config) GetDBConn() string {
	return "user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" host=" + c.DBHost +
		" port=" + c.DBPort +
		" sslmode=disable"
}
