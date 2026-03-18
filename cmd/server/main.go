package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Shubham23jha/digital-post-office/config"
	"github.com/Shubham23jha/digital-post-office/internal/bootstrap"
	"github.com/Shubham23jha/digital-post-office/internal/routes"
	"github.com/Shubham23jha/digital-post-office/pkg/database"
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

	app, err := bootstrap.InitializeApp(database.DB)
	if err != nil {
		panic(err)
	}

	// Server
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	routes.Register(r, app)

	// Serve Static UI
	r.Static("/static", "./web-server")
	r.StaticFile("/", "./web-server/index.html")

	// In Distributed Mode (K8s), these will run as separate pods.
	// For local development, we can still start them if needed.
	if config.GetEnv("RUN_BACKGROUND_SERVICES") == "true" {
		log.Println("⏳ Starting background services (Monolith Mode)...")
		go app.Publisher.Start(context.TODO())
		go app.WorkerPool.Start(context.TODO())
	}

	r.Run(":" + config.GetEnv("PORT"))
}
