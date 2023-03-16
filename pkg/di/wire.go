//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	http "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/db"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase"
)

func InitializeApi(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectDatbase,
		repository.NewAdminRepository, repository.NewUserRepository, repository.NewProductRepository,
		usecase.NewAdminUseCase, usecase.NewUserUseCase, usecase.NewProductUseCase,
		handler.NewAdminHandler, handler.NewUserHandler, handler.NewProductHandler,
		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
