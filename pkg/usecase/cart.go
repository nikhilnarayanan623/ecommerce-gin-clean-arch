package usecase

import (
	"context"
	"log"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type cartUseCase struct {
	cartRepo    interfaces.CartRepository
	productRepo interfaces.ProductRepository
}

func NewCartUseCase(cartRepo interfaces.CartRepository, productRepo interfaces.ProductRepository) service.CartUseCase {
	return &cartUseCase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// get user cart(it include total price and cartId)
func (c *cartUseCase) GetUserCart(ctx context.Context, userID uint) (cart domain.Cart, err error) {

	cart, err = c.cartRepo.FindCartByUserID(ctx, userID)
	if err != nil {
		return domain.Cart{}, utils.PrependMessageToError(err, "failed to get user cart")
	}

	return cart, nil
}

func (c *cartUseCase) SaveToCart(ctx context.Context, body request.Cart) error {

	productItem, err := c.productRepo.FindProductItem(ctx, body.ProductItemID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to find product items")
	}
	if productItem.ID == 0 {
		return ErrInvalidProductItemID
	}

	// check productItem is out of stock or not
	if productItem.QtyInStock == 0 {
		return ErrProductItemOutOfStock
	}

	// find the cart of user
	cart, err := c.cartRepo.FindCartByUserID(ctx, body.UserID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to find user cart")
	}
	if cart.ID == 0 { // if there is no cart is available for user then create new cart
		cart.ID, err = c.cartRepo.SaveCart(ctx, body.UserID)
		if err != nil {
			return err
		}
	}

	// check the given product item is already exit in user cart
	cartItem, err := c.cartRepo.FindCartItemByCartAndProductItemID(ctx, cart.ID, body.ProductItemID)
	if err != nil {
		return err
	}
	if cartItem.ID != 0 {
		return ErrCartItemAlreadyExist
	}

	// add productItem to cartItem
	err = c.cartRepo.SaveCartItem(ctx, cart.ID, body.ProductItemID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to save product items as cart item")
	}

	log.Printf("product_item_id %v saved on cart of  user_id %v", body.ProductItemID, body.UserID)
	return nil
}

func (c *cartUseCase) RemoveCartItem(ctx context.Context, body request.Cart) error {

	// validate the product
	productItem, err := c.productRepo.FindProductItem(ctx, body.ProductItemID)

	if err != nil {
		return err
	}
	if productItem.ID == 0 {
		return ErrInvalidProductItemID
	}

	// Find cart of user
	cart, err := c.cartRepo.FindCartByUserID(ctx, body.UserID)
	if err != nil {
		return err
	}
	if cart.ID == 0 {
		return ErrEmptyCart
	}

	// check the product_item exist on user cart
	cartItem, err := c.cartRepo.FindCartItemByCartAndProductItemID(ctx, cart.ID, body.ProductItemID)
	if err != nil {
		return err
	} else if cartItem.ID == 0 {
		return ErrCartItemNotExit
	}

	// then remove cart_item
	err = c.cartRepo.DeleteCartItem(ctx, cartItem.ID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to delete product item from cart")
	}

	log.Printf("product_item with id %v removed form user cart with user id %v", body.ProductItemID, body.UserID)
	return nil
}

func (c *cartUseCase) UpdateCartItem(ctx context.Context, body request.UpdateCartItem) error {

	const updateMaxQty = 100
	//check the given product_item_id is valid or not
	productItem, err := c.productRepo.FindProductItem(ctx, body.ProductItemID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to find product items")
	} else if productItem.ID == 0 {
		return ErrInvalidProductItemID
	}

	if body.Count < 1 {
		return ErrRequireMinimumCartItemQty
	}

	if body.Count > productItem.QtyInStock || body.Count > updateMaxQty {
		return ErrInvalidCartItemUpdateQty
	}

	// find the cart of user
	cart, err := c.cartRepo.FindCartByUserID(ctx, body.UserID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed find user cart")
	}
	if cart.ID == 0 {
		return ErrEmptyCart
	}

	// find the cart_item with given product_id and user cart_id  and check the product_item present in cart or no
	cartItem, err := c.cartRepo.FindCartItemByCartAndProductItemID(ctx, cart.ID, body.ProductItemID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to find product item from cart")
	}
	if cartItem.ID == 0 {
		return ErrCartItemNotExit
	}

	// update the cart_item qty
	if err := c.cartRepo.UpdateCartItemQty(ctx, cartItem.ID, body.Count); err != nil {
		return utils.PrependMessageToError(err, "failed to update cart item qty")
	}

	log.Printf("updated user carts product_items qty of product_item_id %v , qty %v", body.ProductItemID, body.Count)
	return nil
}

func (c *cartUseCase) GetUserCartItems(ctx context.Context, cartId uint) (cartItems []response.CartItem, err error) {
	// get the cart_items of user
	cartItems, err = c.cartRepo.FindAllCartItemsByCartID(ctx, cartId)
	if err != nil {
		return cartItems, utils.PrependMessageToError(err, "failed to find all cart items")
	}

	log.Printf("successfully got all cart_items of user with cart_id %v", cartId)
	return cartItems, nil
}
