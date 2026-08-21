package config

import (
    "os"
    "strconv"
)

type Config struct {
    Port        string
    DatabaseURL string
    JWTSecret   string
    AutoMigrate bool
}

func Load() *Config {
    return &Config{
        Port:        getEnvOr("PORT", "8080"),
        DatabaseURL: getEnvOr("DATABASE_URL", "postgres://postgres:cupd@localhost:5432/cupcake_delivery?sslmode=disable"),
        JWTSecret:   getEnvOr("JWT_SECRET", "seu_jwt_secret_aqui"),
        AutoMigrate: getEnvBool("AUTO_MIGRATE", true),
    }
}

func getEnvBool(key string, defaultValue bool) bool {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    parsed, err := strconv.ParseBool(value)
    if err != nil {
        return defaultValue
    }
    return parsed
}

func getEnvOr(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
