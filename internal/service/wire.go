package service

import (
	"fmt"
	"os"

	"github.com/google/wire"
)

func ProvideRabbitMQURL() string {
	host := os.Getenv("RABBITMQ_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("RABBITMQ_USER")
	if user == "" {
		user = "guest"
	}
	pass := os.Getenv("RABBITMQ_PASS")
	if pass == "" {
		pass = "guest"
	}
	port := os.Getenv("RABBITMQ_PORT")
	if port == "" {
		port = "5672"
	}
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
	fmt.Printf("🔌 RabbitMQ URL: %s (Host from env: %s)\n", url, os.Getenv("RABBITMQ_HOST"))
	return url
}

var ProviderSet = wire.NewSet(
	NewSessionService,
	NewUserService,
	NewCampaignService,
	NewOutboxPublisher,
	NewWorkerPool,
	ProvideRabbitMQURL,
	// Add static worker count for now
	wire.Value(5),
)
