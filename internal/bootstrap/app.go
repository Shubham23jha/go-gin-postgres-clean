package bootstrap

import (
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/handlers"
)

type App struct {
	AuthHandler *handlers.AuthHandler
}

func NewApp(
	authHandler *handlers.AuthHandler,
) *App {

	return &App{
		AuthHandler: authHandler,
	}
}
