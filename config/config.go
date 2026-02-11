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
		err := godotenv.Load()
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
