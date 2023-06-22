package usecase

import (
	"context"

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

func (c *cartUseCase) SaveProductItemToCart(ctx context.Context, userID, productItemId uint) error {

	productItem, err := c.productRepo.FindProductItemByID(ctx, productItemId)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to find product items")
	}

	// check productItem is out of stock or not
	if productItem.QtyInStock == 0 {
		return ErrProductItemOutOfStock
	}

	// find the cart of user
	cart, err := c.cartRepo.FindCartByUserID(ctx, userID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to find user cart")
	}
	if cart.ID == 0 { // if there is no cart is available for user then create new cart
		cart.ID, err = c.cartRepo.SaveCart(ctx, userID)
		if err != nil {
			return err
		}
	}

	// check the given product item is already exit in user cart
	cartItem, err := c.cartRepo.FindCartItemByCartAndProductItemID(ctx, cart.ID, productItemId)
	if err != nil {
		return err
	}
	if cartItem.ID != 0 {
		return ErrCartItemAlreadyExist
	}

	// add productItem to cartItem
	err = c.cartRepo.SaveCartItem(ctx, cart.ID, productItemId)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to save product items as cart item")
	}

	return nil
}

func (c *cartUseCase) RemoveProductItemFromCartItem(ctx context.Context, userID, productItemId uint) error {

	// Find cart of user
	cart, err := c.cartRepo.FindCartByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if cart.ID == 0 {
		return ErrEmptyCart
	}

	// check the product_item exist on user cart
	cartItem, err := c.cartRepo.FindCartItemByCartAndProductItemID(ctx, cart.ID, productItemId)
	if err != nil {
		return err
	} else if cartItem.ID == 0 {
		return ErrCartItemNotExit
	}

	// then remove cart_item
	err = c.cartRepo.DeleteCartItem(ctx, cartItem.ID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to remove product item from cart")
	}

	return nil
}

func (c *cartUseCase) UpdateCartItem(ctx context.Context, updateDetails request.UpdateCartItem) error {

	const maxCartItemQty = 100
	//check the given product_item_id is valid or not
	productItem, err := c.productRepo.FindProductItemByID(ctx, updateDetails.ProductItemID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to find product items")
	}

	if updateDetails.Count < 1 {
		return ErrRequireMinimumCartItemQty
	}

	if updateDetails.Count > productItem.QtyInStock || updateDetails.Count > maxCartItemQty {
		return ErrInvalidCartItemUpdateQty
	}

	// find the cart of user
	cart, err := c.cartRepo.FindCartByUserID(ctx, updateDetails.UserID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed find user cart")
	}
	if cart.ID == 0 {
		return ErrEmptyCart
	}

	// find the cart_item with given product_id and user cart_id  and check the product_item present in cart or no
	cartItem, err := c.cartRepo.FindCartItemByCartAndProductItemID(ctx, cart.ID, updateDetails.ProductItemID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to find product item from cart")
	}
	if cartItem.ID == 0 {
		return ErrCartItemNotExit
	}

	// update the cart_item qty
	if err := c.cartRepo.UpdateCartItemQty(ctx, cartItem.ID, updateDetails.Count); err != nil {
		return utils.PrependMessageToError(err, "failed to update cart item qty")
	}

	return nil
}

func (c *cartUseCase) GetUserCartItems(ctx context.Context, cartId uint) (cartItems []response.CartItem, err error) {
	// get the cart_items of user
	cartItems, err = c.cartRepo.FindAllCartItemsByCartID(ctx, cartId)
	if err != nil {
		return cartItems, utils.PrependMessageToError(err, "failed to find all cart items")
	}

	return cartItems, nil
}
