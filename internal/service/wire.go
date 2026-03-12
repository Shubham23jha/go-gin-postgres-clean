package service

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewSessionService,
	NewUserService,
	NewCampaignService,
	NewOutboxPublisher,
	// Provider for RabbitMQ connection string
	wire.Value("amqp://guest:guest@localhost:5672/"),
)
