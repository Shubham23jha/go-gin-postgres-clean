// Package bootstrap handles the application initialization and dependency injection.
package bootstrap

import (
	"github.com/Shubham23jha/digital-post-office/internal/handlers"
	"github.com/Shubham23jha/digital-post-office/internal/service"
)

type App struct {
	AuthHandler     *handlers.AuthHandler
	CampaignHandler *handlers.CampaignHandler
	Publisher       service.OutboxPublisher
	WorkerPool       service.WorkerPool
}

func NewApp(
	authHandler *handlers.AuthHandler,
	campaignHandler *handlers.CampaignHandler,
	publisher service.OutboxPublisher,
	workerPool service.WorkerPool,
) *App {

	return &App{
		AuthHandler:     authHandler,
		CampaignHandler: campaignHandler,
		Publisher:       publisher,
		WorkerPool:      workerPool,
	}
}
