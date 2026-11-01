package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	// which Env mean which environment we operating
	Env    string
	DB_URL string
}

// Must pattern

func MustLoad() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Env File not found!")
	}

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is requried")
	}

	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is requried")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		panic("DATABASE_URL is requried")
	}

	return Config{
		Port:   port,
		Env:    env,
		DB_URL: dbUrl,
	}
}
