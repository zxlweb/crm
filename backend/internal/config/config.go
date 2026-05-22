package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Env        string
	GinMode    string
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret           string
	AccessTokenTTLMin   int
	RefreshTokenTTLDays int
	AutoMigrate         bool
}

func LoadConfig() *Config {
	cfg := &Config{
		Env:         getEnv("APP_ENV", "development"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		Port:        getEnv("PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "password"),
		DBName:      getEnv("DB_NAME", "crm"),
		JWTSecret:           getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		AccessTokenTTLMin:   getEnvInt("ACCESS_TOKEN_TTL_MIN", 60),
		RefreshTokenTTLDays: getEnvInt("REFRESH_TOKEN_TTL_DAYS", 7),
		AutoMigrate:         getEnv("AUTO_MIGRATE", "false") == "true",
	}
	gin.SetMode(cfg.GinMode)
	return cfg
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort,
	)
}

func (c *Config) AccessTokenTTL() time.Duration {
	return time.Duration(c.AccessTokenTTLMin) * time.Minute
}

func (c *Config) RefreshTokenTTL() time.Duration {
	return time.Duration(c.RefreshTokenTTLDays) * 24 * time.Hour
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if n, err := strconv.Atoi(value); err == nil {
			return n
		}
	}
	return defaultValue
}
