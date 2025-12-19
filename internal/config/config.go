package config

import (
	"os"
	"strconv"
)

type Config struct {
	Env       string // development, staging, production
	App       AppConfig
	DB        DBConfig
	JWT       JWTConfig
	ShortCode ShortCodeConfig
	RateLimit RateLimitConfig
	Redis     RedisConfig
}

type AppConfig struct {
	Host   string
	Port   string
	Domain string
	Debug  bool
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Enabled  bool
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
	env := getEnv("APP_ENV", "development")

	return &Config{
		Env: env,
		App: AppConfig{
			Host:   getEnv("APP_HOST", "0.0.0.0"),
			Port:   getEnv("APP_PORT", "8080"),
			Domain: getEnv("APP_DOMAIN", "localhost:8080"),
			Debug:  env != "production",
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
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			Enabled:  getEnvBool("REDIS_ENABLED", false),
		},
	}
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Env == "production"
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

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
