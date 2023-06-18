package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type CartUseCase interface {
	SaveToCart(ctx context.Context, body request.Cart) error               // save product_item to cart
	RemoveCartItem(ctx context.Context, body request.Cart) error           // remove product_item from cart
	UpdateCartItem(ctx context.Context, body request.UpdateCartItem) error // edit cartItems( quantity change )
	GetUserCart(ctx context.Context, userID uint) (cart domain.Cart, err error)
	GetUserCartItems(ctx context.Context, cartId uint) (cartItems []response.CartItem, err error)
}
