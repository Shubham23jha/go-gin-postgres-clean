package bootstrap

import (
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/handlers"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/repository"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/service"
	"gorm.io/gorm"
)

type App struct {
	AuthHandler *handlers.AuthHandler
}

func NewApp(db *gorm.DB) *App {

	// ===== Repositories =====
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// ===== Services =====
	sessionService := service.NewSessionService(sessionRepo)

	userService := service.NewUserService(userRepo, sessionRepo)

	// ===== Handlers =====
	authHandler := handlers.NewAuthHandler(userService, sessionService)

	return &App{
		AuthHandler: authHandler,
	}
}
