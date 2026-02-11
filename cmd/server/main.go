package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/Shubham23jha/go-gin-postgres-clean/config"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/handler"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/repository"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/routes"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/service"
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

	// Layers
	userRepo := repository.NewUserRepository(database.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Server
	r := gin.Default()
	routes.Register(r, userHandler)

	r.Run(":" + config.GetEnv("PORT"))
}


