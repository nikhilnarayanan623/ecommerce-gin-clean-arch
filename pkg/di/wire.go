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

	}

	adminRepo := repository.NewAdminRepository(db)
	userRepo := repository.NewUserRepository(db)

	adminUseCase := usecase.NewAdminUseCase(adminRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)

	adminHandler, userHandler := handler.NewHandlers(adminUseCase, userUseCase)

	return http.NewServerHTTP(adminHandler, userHandler), nil

	// wire.Build(
	// 	db.ConnectDatbase,
	// 	repository.NewAdminRepository, repository.NewUserRepository,
	// 	usecase.NewAdminUseCase, usecase.NewUserUseCase,
	// 	handler.NewHandlers, http.NewServerHTTP,
	// )

	//return &http.ServerHTTP{}, nil
}
