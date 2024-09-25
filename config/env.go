package config

import (
	"fmt"
	"github.com/lpernett/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Host                   string
	Port                   string
	User                   string
	Password               string
	DBName                 string
	SSLMode                string
	JWTSecret              string
	JWTExpirationInSeconds int64
}

var Env = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found or unable to load")
	}
	return Config{
		Port:                   getEnv("DB_PORT", "5432"),
		User:                   getEnv("DB_USER", "postgres"),
		Password:               getEnv("DB_PASSWORD", "6083aee9a93aa5e52dd69323bead66c6"),
		Host:                   getEnv("DB_HOST", "tests-database.cluw2u2yc5me.us-east-1.rds.amazonaws.com"),
		DBName:                 getEnv("DB_NAME", "golangtest"),
		SSLMode:                getEnv("DB_SSL_MODE", "require"),
		JWTSecret:              getEnv("JWT_SECRET", "not-so-secret-now-is-it?"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
