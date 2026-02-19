//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/Shubham23jha/go-gin-postgres-clean/internal/handlers"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/repository"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/service"
)

func InitializeApp(db *gorm.DB) (*App, error) {

	wire.Build(
		repository.ProviderSet,
		service.ProviderSet,
		handlers.ProviderSet,
		NewApp,
	)

	return &App{}, nil
}
