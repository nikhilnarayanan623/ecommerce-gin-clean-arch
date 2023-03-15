//go:build wireinject
// +build wireinject

package di

// func InitializeApi(cfg config.Config) (*http.ServerHTTP, error) {
// 	wire.Build(
// 		db.ConnectDatbase,
// 		repository.NewAdminRepository, repository.NewUserRepository,
// 		usecase.NewAdminUseCase, usecase.NewUserUseCase,
// 		handler.NewHandlers, http.NewServerHTTP,
// 	)

// 	return &http.ServerHTTP{}, nil
// }
