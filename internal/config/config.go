package config

import (
	"os"
	"strconv"
)

type Config struct {
	App       AppConfig
	DB        DBConfig
	JWT       JWTConfig
	ShortCode ShortCodeConfig
	RateLimit RateLimitConfig
}

type AppConfig struct {
	Host   string
	Port   string
	Domain string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret      string
	ExpiryHours int
}

type ShortCodeConfig struct {
	Length   int
	Alphabet string
}

type RateLimitConfig struct {
	Requests int
	Window   int
}

func Load() *Config {
	return &Config{
		App: AppConfig{
			Host:   getEnv("APP_HOST", "0.0.0.0"),
			Port:   getEnv("APP_PORT", "8080"),
			Domain: getEnv("APP_DOMAIN", "localhost:8080"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "shorten_url"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
			ExpiryHours: getEnvInt("JWT_EXPIRY_HOURS", 24),
		},
		ShortCode: ShortCodeConfig{
			Length:   getEnvInt("SHORT_CODE_LENGTH", 6),
			Alphabet: getEnv("SHORT_CODE_ALPHABET", "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"),
		},
		RateLimit: RateLimitConfig{
			Requests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
			Window:   getEnvInt("RATE_LIMIT_WINDOW", 60),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
