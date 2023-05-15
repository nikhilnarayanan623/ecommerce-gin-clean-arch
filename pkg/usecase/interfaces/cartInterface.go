package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type CartUseCase interface {
	SaveToCart(ctx context.Context, body req.Cart) error               // save product_item to cart
	RemoveCartItem(ctx context.Context, body req.Cart) error           // remove product_item from cart
	UpdateCartItem(ctx context.Context, body req.UpdateCartItem) error // edit cartItems( quantity change )
	GetUserCart(ctx context.Context, userID uint) (cart domain.Cart, err error)
	GetUserCartItems(ctx context.Context, cartId uint) (cartItems []res.CartItem, err error)
}
