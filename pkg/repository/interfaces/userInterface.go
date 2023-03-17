package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
)

type UserRepository interface {
	FindUser(ctx context.Context, user domain.User) (domain.User, error)
	SaveUser(ctx context.Context, user domain.User) (domain.User, error)
	//
	FindProductItem(ctx context.Context, productItemID uint) (domain.ProductItem, error)
	FindCart(ctx context.Context, userId uint) (domain.Cart, error)
	UpdateCartPrice(ctx context.Context, cart domain.Cart) (domain.Cart, error)
	//cart
	FindCartItem(ctx context.Context, cartID, productItemID uint) (domain.CartItem, error)
	SaveCartItem(ctx context.Context, cartID, productItemID uint) (domain.CartItem, error)
	RemoveCartItem(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error)
	UpdateCartItem(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error)
	GetCartItems(ctx context.Context, userId uint) (helper.ResponseCart, error)
	//address
	FindAddress(ctx context.Context, address domain.Address) (domain.Address, error)
}
