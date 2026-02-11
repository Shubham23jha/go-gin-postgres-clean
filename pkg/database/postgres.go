package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/Shubham23jha/go-gin-postgres-clean/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func Connect() {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			config.GetEnv("DB_HOST"),
			config.GetEnv("DB_USER"),
			config.GetEnv("DB_PASSWORD"),
			config.GetEnv("DB_NAME"),
			config.GetEnv("DB_PORT"),
			config.GetEnv("DB_SSLMODE"),
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("❌ DB connection failed:", err)
		}

		DB = db
		log.Println("✅ PostgreSQL connected (initialized once)")
	})
}
