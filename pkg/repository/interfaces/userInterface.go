package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
)

type UserRepository interface {
	FindUser(ctx context.Context, user domain.User) (domain.User, error)
	SaveUser(ctx context.Context, user domain.User) (domain.User, error)

	AddToCart(ctx context.Context, body helper.ReqCart) (domain.Cart, error)
	GetCartItems(ctx context.Context, userId uint) (helper.ResponseCart, error)
}
