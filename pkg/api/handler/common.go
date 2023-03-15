package handler

import (
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

func NewHandlers(adminUseCase service.AdminUseCase, userUseCase service.UserUseCase, productUseCase service.ProductUseCase) (*AdminHandler, *UserHandler, *ProductHandler) {

	return &AdminHandler{adminUseCase: adminUseCase},
		&UserHandler{userUseCase: userUseCase},
		&ProductHandler{productUseCase: productUseCase}
}
