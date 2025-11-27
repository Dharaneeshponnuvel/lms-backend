package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	AppEnv  string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	JWTSecret       string
	AccessTokenTTL  int64
	RefreshTokenTTL int64
}

func Load() *Config {
	cfg := &Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "auth_service_db"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		JWTSecret:     getEnv("JWT_SECRET", "secret"),
	}

	// Simple: keep RedisDB = 0
	cfg.RedisDB = 0

	cfg.AccessTokenTTL = mustParseInt64(getEnv("ACCESS_TOKEN_TTL", "3600"))
	cfg.RefreshTokenTTL = mustParseInt64(getEnv("REFRESH_TOKEN_TTL", "604800"))

	return cfg
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func mustParseInt64(s string) int64 {
	var v int64
	_, err := fmt.Sscan(s, &v)
	if err != nil {
		log.Fatalf("invalid int64 for config: %s", s)
	}
	return v
}
