package config

import (
	"os"
"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	MongoURI       string
	DBName         string
	JWTPublicKey   string
	ServiceKey     string
	AllowedOrigins []string
}

func parseOrigins(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func Load() *Config {
	if _, err := os.Stat(".env.local"); err == nil {
		_ = godotenv.Load(".env.local")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "tracer"
	}

	return &Config{
		Port:           os.Getenv("PORT"),
		MongoURI:       os.Getenv("MONGODB_URI"),
		DBName:         dbName,
		JWTPublicKey:   os.Getenv("JWT_PUBLIC_KEY"),
		ServiceKey:     os.Getenv("SERVICE_KEY"),
		AllowedOrigins: parseOrigins(os.Getenv("ALLOWED_ORIGINS")),
	}
}
