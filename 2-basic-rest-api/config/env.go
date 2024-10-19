package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	USER     string
	PASSWORD string
	HOST     string
	PORT     string
	DBNAME   string
}

var Envs = initconfig()

func initconfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	return Config{
		USER:     getEnv("HOST", "root"),
		PASSWORD: getEnv("PASSWORD", "mypassword"),
		HOST:     getEnv("HOST", "localhost"),
		PORT:     getEnv("PORT", "5432"),
		DBNAME:   getEnv("DBNAME", "restapi"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
