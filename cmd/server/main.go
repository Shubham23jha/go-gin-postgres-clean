package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/Shubham23jha/go-gin-postgres-clean/config"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/bootstrap"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/routes"
	"github.com/Shubham23jha/go-gin-postgres-clean/pkg/database"
)

func main() {
	config.LoadEnv()

	// DB
	database.Connect()

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.GetEnv("DB_USER"),
		config.GetEnv("DB_PASSWORD"),
		config.GetEnv("DB_HOST"),
		config.GetEnv("DB_PORT"),
		config.GetEnv("DB_NAME"),
		config.GetEnv("DB_SSLMODE"),
	)

	database.RunMigrations(dbURL)

	app := bootstrap.NewApp(database.DB)

	// Server
	r := gin.Default()

	routes.Register(r, app)

	r.Run(":" + config.GetEnv("PORT"))
}
