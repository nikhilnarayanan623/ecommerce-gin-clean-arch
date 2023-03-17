package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
)

type UserUseCase interface {
	Signup(ctx context.Context, user domain.User) (domain.User, error)
	Login(ctx context.Context, user domain.User) (domain.User, error)
	LoginOtp(ctx context.Context, user domain.User) (domain.User, error)

	Home(ctx context.Context, userId uint) (domain.User, error)

	SaveToCart(ctx context.Context, body helper.ReqCart) (domain.CartItem, error)
	RemoveCartItem(ctx context.Context, body helper.ReqCart) (domain.Cart, error)
	UpdateCartItem(ctx context.Context, body helper.ReqCartCount) (domain.CartItem, error)

	GetCartItems(ctx context.Context, userId uint) (helper.ResponseCart, error)
}
