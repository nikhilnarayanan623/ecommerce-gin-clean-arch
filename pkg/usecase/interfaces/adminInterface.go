package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

type AdminUseCase interface {
	Login(ctx context.Context, admin domain.Admin) (domain.Admin, error)
	SignUp(ctx context.Context, admin domain.Admin) error

	FindAllUser(ctx context.Context) ([]res.UserRespStrcut, error)
	BlockUser(ctx context.Context, userID uint) error
}

// GetCategory(ctx context.Context) (helper.ReqCategory, any)
// 	SetCategory(ctx context.Context, body helper.ReqCategory)
