package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerPort     string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	LogLevel       string
	MaxConnections int
}

// Load загружает конфигурацию из .env
func Load() (*Config, error) {
	// Загружаем .env файл, если есть
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "subscription"),
		DBSSLMode:      getEnv("DB_SSL_MODE", "disable"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 10),
	}, nil
}

// GetDBConnectionString - возвращает строку DSN для подключения к PostgreSQL
func (c *Config) GetDBConnectionString() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword +
		"@" + c.DBHost + ":" + c.DBPort +
		"/" + c.DBName +
		"?sslmode=" + c.DBSSLMode
}

// getEnv — читает переменную, если есть. Если нет — дефолт
func getEnv(key string, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultValue
}

// getEnvAsInt — читает переменную как int. Если нет или ошибка — дефолт.
func getEnvAsInt(key string, defaultValue int) int {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}
