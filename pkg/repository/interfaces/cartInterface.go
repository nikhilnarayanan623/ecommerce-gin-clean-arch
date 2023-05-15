package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type CartRepository interface {
	FindCartByUserID(ctx context.Context, userID uint) (cart domain.Cart, err error)
	SaveCart(ctx context.Context, userID uint) (cartID uint, err error)
	UpdateCart(ctx context.Context, cartId, discountAmount, couponID uint) error

	FindCartItemByCartAndProductItemID(ctx context.Context, cartID, productItemID uint) (cartItem domain.CartItem, err error)
	FindAllCartItemsByCartID(ctx context.Context, cartID uint) (cartItems []res.CartItem, err error)
	SaveCartItem(ctx context.Context, cartId, productItemId uint) error
	DeleteCartItem(ctx context.Context, cartItemID uint) error
	DeleteAllCartItemsByUserID(ctx context.Context, userID uint) error // this used for order time cart remove using user_id
	DeleteAllCartItemsByCartID(ctx context.Context, cartID uint) error
	UpdateCartItemQty(ctx context.Context, cartItemId, qty uint) error

	CheckcartIsValidForOrder(ctx context.Context, userID uint) (cart domain.Cart, err error)
}
