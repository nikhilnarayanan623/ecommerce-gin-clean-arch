package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
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
		return cart, fmt.Errorf("faild to get cart %v", err.Error())
	}

	return cart, nil
}

func (c *cartUseCase) SaveToCart(ctx context.Context, body req.Cart) (err error) {

	// get the productitem to check product is valid
	productItem, err := c.productRepo.FindProductItem(ctx, body.ProductItemID)
	if err != nil {
		return fmt.Errorf("faild to find product items error: %v", err.Error())
	} else if productItem.ID == 0 {
		return errors.New("invalid product_item id")
	}

	// check productItem is out of stock or not
	if productItem.QtyInStock == 0 {
		return errors.New("product is now out of stock")
	}

	// find the cart of user
	cart, err := c.cartRepo.FindCartByUserID(ctx, body.UserID)
	if err != nil {
		return err
	} else if cart.ID == 0 { // if there is no cart is available for user then create it
		cart.ID, err = c.cartRepo.SaveCart(ctx, body.UserID)
		if err != nil {
			return err
		}
		log.Println(cart.ID)
	}

	// check the given product item is already exit in user cart
	cartItem, err := c.cartRepo.FindCartItemByCartAndProductItemID(ctx, cart.ID, body.ProductItemID)
	if err != nil {
		return err
	} else if cartItem.ID != 0 {
		return errors.New("product_item already exist on the cart can't save product to cart")
	}

	// add productItem to cartItem
	if err := c.cartRepo.SaveCartItem(ctx, cart.ID, body.ProductItemID); err != nil {
		return err
	}

	log.Printf("product_item_id %v saved on cart of  user_id %v", body.ProductItemID, body.UserID)
	return nil
}

func (c *cartUseCase) RemoveCartItem(ctx context.Context, body req.Cart) error {

	// validate the product
	productItem, err := c.productRepo.FindProductItem(ctx, body.ProductItemID)

	if err != nil {
		return err
	} else if productItem.ID == 0 {
		return errors.New("invalid product_id")
	}

	// Find cart of user
	cart, err := c.cartRepo.FindCartByUserID(ctx, body.UserID)
	if err != nil {
		return err
	} else if cart.ID == 0 {
		return errors.New("can't remove product_item from user cart \n user cart is empty")
	}

	// check the product_item exist on user cart
	cartItem, err := c.cartRepo.FindCartItemByCartAndProductItemID(ctx, cart.ID, body.ProductItemID)
	if err != nil {
		return err
	} else if cartItem.ID == 0 {
		return fmt.Errorf("prduct_item with id %v is not exist on user cart", body.ProductItemID)
	}
	// then remvoe cart_item
	err = c.cartRepo.DeleteCartItem(ctx, cartItem.ID)
	if err != nil {
		return err
	}

	log.Printf("product_item with id %v removed form user cart with usre id %v", body.ProductItemID, body.UserID)
	return nil
}

func (c *cartUseCase) UpdateCartItem(ctx context.Context, body req.UpdateCartItem) error {

	//check the given product_item_id is valid or not
	productItem, err := c.productRepo.FindProductItem(ctx, body.ProductItemID)
	if err != nil {
		return err
	} else if productItem.ID == 0 {
		return errors.New("invalid product_item_id")
	}

	if body.Count < 1 {
		return fmt.Errorf("can't change cart_item qty to %v \n minuimun qty is 1", body.Count)
	}

	// check the given qty of product_item updation is morethan prodct_item qty
	if body.Count > productItem.QtyInStock {
		return errors.New("there is not this much quantity available in product_item")
	}

	// find the cart of user
	cart, err := c.cartRepo.FindCartByUserID(ctx, body.UserID)
	if err != nil {
		return err
	} else if cart.ID == 0 {
		return errors.New("user cart is empty")
	}

	// find the cart_item with given product_id and user cart_id  and check the product_item present in cart or no
	cartItem, err := c.cartRepo.FindCartItemByCartAndProductItemID(ctx, cart.ID, body.ProductItemID)
	if err != nil {
		return err
	} else if cartItem.ID == 0 {
		return fmt.Errorf("product_item not exist in the cart with given product_item_id %v", body.ProductItemID)
	}

	// update the cart_item qty
	if err := c.cartRepo.UpdateCartItemQty(ctx, cartItem.ID, body.Count); err != nil {
		return err
	}

	log.Printf("updated user carts product_items qty of product_item_id %v , qty %v", body.ProductItemID, body.Count)
	return nil
}

func (c *cartUseCase) GetUserCartItems(ctx context.Context, cartId uint) (cartItems []res.CartItem, err error) {
	// get the cart_items of user
	cartItems, err = c.cartRepo.FindAllCartItemsByCartID(ctx, cartId)
	if err != nil {
		return cartItems, err
	}

	log.Printf("sucessfully got all cart_items of user wtih cart_id %v", cartId)
	return cartItems, nil
}
