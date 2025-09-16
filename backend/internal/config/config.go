package config

import (
    "os"
)

type Config struct {
    Port        string
    DatabaseURL string
    JWTSecret   string
}

func Load() *Config {
    return &Config{
        Port:        getEnvOr("PORT", "8080"),
        DatabaseURL: getEnvOr("DATABASE_URL", "postgres://postgres:cupd@localhost:5432/cupcake_delivery?sslmode=disable"),
        JWTSecret:   getEnvOr("JWT_SECRET", "seu_jwt_secret_aqui"),
    }
}

func getEnvOr(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
