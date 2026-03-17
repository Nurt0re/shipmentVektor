package pkg

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConn string
	Port   string
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	return &Config{
		DBConn: os.Getenv("DB_CONN"),
		Port:   os.Getenv("PORT"),
	}
}
