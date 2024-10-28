package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	USER       string
	PASSWORD   string
	HOST       string
	PORT       string
	DBNAME     string
	JWTExpired int64
	JWTSecret  string
}

var Envs = initconfig()

func initconfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	return Config{
		USER:       getEnv("HOST", "root"),
		PASSWORD:   getEnv("PASSWORD", "mypassword"),
		HOST:       getEnv("HOST", "localhost"),
		PORT:       getEnv("PORT", "5432"),
		DBNAME:     getEnv("DBNAME", "bookstore"),
		JWTExpired: getEnvAsInt("JWTExp", 3600*24*7),
		JWTSecret:  getEnv("JWTSecret", "most_secret_value"),
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
