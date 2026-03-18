//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/Shubham23jha/digital-post-office/internal/handlers"
	"github.com/Shubham23jha/digital-post-office/internal/repository"
	"github.com/Shubham23jha/digital-post-office/internal/service"
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
