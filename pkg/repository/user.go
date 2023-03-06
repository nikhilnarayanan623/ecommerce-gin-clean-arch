package repository

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type userDatabse struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabse{DB: DB}
}

func (c *userDatabse) FindUser(ctx context.Context, user domain.Users) (domain.Users, any) {

	// check id,email,phone any of then match i db
	c.DB.Raw("SELECT * FROM users where id=? OR email=? OR phone=?", user.ID, user.Email, user.Phone).Scan(&user)

	// if given userid then check mail is stil there otherwise phone or id
	if user.ID == 0 || user.Email == "" || user.Phone == "" {
		return user, map[string]string{"Error": "Can't find the user"}
	}
	// if found the user then retur user with nil
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

func (c *userDatabse) GetAllProducts(ctx context.Context) ([]domain.Product, any) {

	var products []domain.Product

	return products, c.DB.Raw("SELECT * FROM products").Scan(&products).Error

}

func (c *userDatabse) GetProductItems(ctx context.Context, product domain.Product) ([]domain.ProductItem, any) {

	var productItems []domain.ProductItem

	return productItems, c.DB.Raw("SELECT * FROM product_items WHERE product_id", product.ID).Scan(&productItems).Error
}
