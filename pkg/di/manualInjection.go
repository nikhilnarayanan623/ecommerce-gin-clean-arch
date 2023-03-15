package di

import (
	http "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/db"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase"
)

func InitializeApi(cfg config.Config) (*http.ServerHTTP, error) {

	db, err := db.ConnectDatbase(cfg)
	if err != nil {
		return nil, err
	}
	adminRepo := repository.NewAdminRepository(db)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)

	adminUseCase := usecase.NewAdminUseCase(adminRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)

	productUseCase := usecase.NewProductUseCase(productRepo)

	adminHandler, userHandler, productHandler := handler.NewHandlers(adminUseCase, userUseCase, productUseCase)

	server := http.NewServerHTTP(adminHandler, userHandler, productHandler)

	return server, nil
}
