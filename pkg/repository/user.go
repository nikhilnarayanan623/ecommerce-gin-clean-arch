package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type userDatabse struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabse{DB: DB}
}

func (c *userDatabse) FindUser(ctx context.Context, user domain.Users) (domain.Users, error) {
	fmt.Println("user", user)

	// check id,email,phone any of then match i db
	c.DB.Raw("SELECT * FROM users where id=? OR email=? OR phone=?", user.ID, user.Email, user.Phone).Scan(&user)

	// if given userid then check mail is stil there otherwise phone or id
	if user.ID == 0 || user.Email == "" || user.Phone == "" {
		return user, errors.New("Can't find the user")
	}
	// if found the user then return user with nil
	return user, nil
}

func (c *userDatabse) SaveUser(ctx context.Context, user domain.Users) (domain.Users, any) {

	// check whether user is already exisist
	c.DB.Raw("SELECT * FROM users WHERE email=? OR phone=?", user.Email, user.Phone).Scan(&user)
	//if exist then return message as user exist
	if user.ID != 0 {
		return user, map[string]string{"msg": "User Already Exist"}
	}

	// return user with save status
	return user, c.DB.Save(&user).Error
}

func (c *userDatabse) GetAllProducts(ctx context.Context) ([]helper.ResponseProduct, any) {

	var products []helper.ResponseProduct

	querry := `SELECT f.product_name,f.description,f.price,s.category_name,f.image FROM products f LEFT JOIN categories s ON f.category_id=s.id`

	err := c.DB.Raw(querry).Scan(&products).Error

	return products, err

}

func (c *userDatabse) GetProductItems(ctx context.Context, product domain.Product) ([]domain.ProductItem, any) {

	var productItems []domain.ProductItem

	return productItems, c.DB.Raw("SELECT * FROM product_items WHERE product_id", product.ID).Scan(&productItems).Error
}

func (c *userDatabse) GetCartItems(ctx context.Context, userId uint) (helper.ResCart, any) {

	var (
		user = domain.Users{ID: userId}
		// resCart = helper.ResCart{CartItems: make([]helper.ResCartItem, 0)}
		resCart helper.ResCart
		cart    domain.Cart
		//cartItems []domain.CartItem
	)

	//first find the user
	user, err := c.FindUser(ctx, user)

	if err != nil {
		return resCart, err
	}

	// then get cart id of user
	if c.DB.Raw("SELECT * FROM carts WHERE users_id=?", userId).Scan(&cart); cart.ID == 0 {
		return resCart, helper.SingleRespStruct{Error: "User Have no cart"} // I think I want to delete it later
	}

	// add total price to response
	resCart.TotalPrice = cart.TotalPrice

	//then get all cart items of user
	var resCartItem []helper.ResCartItem

	c.DB.Raw("SELECT * FROM cart_items WHERE cart_id=?", cart.ID).Scan(&resCartItem)

	// assign it to resCart
	resCart.CartItems = resCartItem

	return resCart, nil
}
