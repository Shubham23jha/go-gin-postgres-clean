package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shubham23jha/go-gin-postgres-clean/config"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/bootstrap"
	"github.com/Shubham23jha/go-gin-postgres-clean/pkg/database"
)

func main() {
	config.LoadEnv()
	database.Connect()

	app, err := bootstrap.InitializeApp(database.DB)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("👷 Email Worker Service Starting...")
	
	// Start Worker Pool in this process
	go app.WorkerPool.Start(ctx)

	<-ctx.Done()
	log.Println("🛑 Email Worker Service Stopped")
}
