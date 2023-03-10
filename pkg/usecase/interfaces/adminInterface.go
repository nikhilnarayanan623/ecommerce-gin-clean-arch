package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
)

type AdminUseCase interface {
	Login(ctx context.Context, admin domain.Admin) (domain.Admin, any)
	SignUp(ctx context.Context, admin domain.Admin) (domain.Admin, any)
	FindAllUser(ctx context.Context) ([]helper.UserRespStrcut, error)
	BlockUser(ctx context.Context, user helper.BlockStruct) (domain.Users, any)
	AddCategory(ctx context.Context, productCategory domain.Category) (domain.Category, any)
}
