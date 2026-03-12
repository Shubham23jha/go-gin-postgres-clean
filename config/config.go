package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

func LoadEnv() {
	once.Do(func() {
		// Try root first (for go run cmd/server/main.go)
		err := godotenv.Load(".env")
		if err != nil {
			// Try two levels up (for tests)
			err = godotenv.Load("../../.env")
		}

		if err != nil {
			log.Println("⚠️ No .env file found, using system env")
		} else {
			log.Println("✅ .env loaded")
		}
	})
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
