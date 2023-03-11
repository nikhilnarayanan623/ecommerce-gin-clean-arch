package repository

import (
	"context"
	"errors"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type adminDatabase struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminDatabase{DB: DB}
}

func (c *adminDatabase) FindAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, error) {

	c.DB.Raw("SELECT * FROM admins WHERE email=? OR user_name=?", admin.Email, admin.UserName).Scan(&admin)

	//check the admin got or not
	if admin.ID == 0 {
		return admin, errors.New("admin not exist with this details")
	}

	return admin, nil
}

func (c *adminDatabase) SaveAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, error) {

	// first check the admin already exist or not
	var dbAdmin domain.Admin
	c.DB.Raw("SELECT * FROM admins WHERE email=? OR user_name=?", admin.Email, admin.UserName).Scan(&dbAdmin)

	if dbAdmin.ID != 0 { // amdmin already exist

		//errMap := map[string]string{}
		// first check the email is already exist
		if dbAdmin.Email == admin.Email {
			// errMap["email"] = "Email Already exist"
			return admin, errors.New("admin already exist with this email")
		} // if email not then its user name is exist
		//errMap["user_name"] = "UserName Already exist"
		return admin, errors.New("admin already exist with this user_name")
	}

	//if admin not exist then create it
	querry := `INSERT INTO admins (user_name,email,password) VALUES ($1,$2,$3) RETURNING user_name,email,password`
	if c.DB.Raw(querry, admin.UserName, admin.Email, admin.Password).Scan(&admin).Error != nil {
		return admin, errors.New("faild to create account for admin")
	}

	c.DB.Raw("SELECT * FROM admins WHERE email=?", admin.Email).Scan(&admin)
	return admin, nil // successfully admin added
}

func (c *adminDatabase) FindAllUser(ctx context.Context) ([]domain.Users, error) {

	var users []domain.Users
	err := c.DB.Raw("SELECT * FROM users").Scan(&users).Error

	return users, err
}

func (c *adminDatabase) BlockUser(ctx context.Context, user domain.Users) (domain.Users, error) {

	var dbUser domain.Users
	// first check ther user valid or not
	c.DB.Raw("SELECT * FROM users WHERE id=?", user.ID).Scan(&dbUser)

	if dbUser.ID == 0 {
		return user, errors.New("invalid user id user doesn't exist")
	}

	// if user is blocked then unblock
	if dbUser.BlockStatus {
		c.DB.Raw("UPDATE users SET block_status='F' WHERE id=?", user.ID).Scan(&dbUser)
		dbUser.BlockStatus = false
		return dbUser, nil
	}
	c.DB.Raw("UPDATE users SET block_status='T' WHERE id=?", user.ID).Scan(&dbUser)
	dbUser.BlockStatus = true
	return dbUser, nil
}

func (c *adminDatabase) GetCategory(ctx context.Context) ([]helper.RespCategory, any) {

	var response []helper.RespCategory
	// left join to get all category and main category
	querry := `SELECT s.id, s.category_name,s.category_id, m.category_name as main_category_name FROM categories s LEFT JOIN categories m ON s.category_id=s.id`

	err := c.DB.Raw(querry).Scan(&response).Error

	return response, err
}

// add category
func (c *adminDatabase) AddCategory(ctx context.Context, category domain.Category) (helper.RespCategory, error) {

	var checkCat domain.Category
	//first check the categoryname already exisits or not
	c.DB.Raw("SELECT * FROM categories WHERE category_name=?", category.CategoryName).Scan(&checkCat)

	if checkCat.ID != 0 { // means category already exist
		return helper.RespCategory{}, errors.New("category already exist")
	}

	// check the given category is main or sub
	if category.CategoryID == 0 { // no catogry id means its main category
		querry := `INSERT INTO categories (category_name)VALUES($1) RETURNING category_name`
		c.DB.Raw(querry, category.CategoryName).Scan(&category)
	} else {
		// first check the category id is valid or not
		c.DB.Raw("SELECT * FROM categories WHERE id=?", category.CategoryID).Scan(&checkCat)
		if checkCat.ID == 0 { // its not a valid category
			return helper.RespCategory{}, errors.New("category_id is invalid main category not exist")
		}
		//otherwise add its with main category
		querry := `INSERT INTO categories (category_id,category_name)VALUES($1,$2) RETURNING category_id,category_name`
		c.DB.Raw(querry, category.CategoryID, category.CategoryName).Scan(&category)
	}

	// response as categoryName main category
	var response helper.RespCategory

	c.DB.Raw(`SELECT * FROM categories WHERE category_name=?`, category.CategoryName).Scan(&response)
	c.DB.Raw(`SELECT category_name AS main_category_name FROM categories WHERE id=?`, category.CategoryID).Scan(&response)

	return response, nil
}

// add variation
func (c *adminDatabase) AddVariation(ctx context.Context, variation domain.Variation) (domain.Variation, error) {

	//firs variation already exist or not
	c.DB.Raw("SELECT * FROM variations WHERE variation_name=?", variation.VariationName).Scan(&variation)
	if variation.ID != 0 {
		return variation, errors.New("variation already exist")
	}

	// then check the category provided for variaion is valid or not
	var cat domain.Category
	c.DB.Raw("SELECT * FROM categories WHERE id=?", variation.CategoryID).Scan(&cat)
	if cat.ID == 0 {
		return variation, errors.New("invalid category_id")
	}

	// if everything ok then add variation
	querry := `INSERT INTO variations (category_id,variation_name) VALUES($1,$2) RETURNING category_id,variation_name`
	if c.DB.Raw(querry, variation.CategoryID, variation.VariationName).Scan(&variation).Error != nil {
		return variation, errors.New("faild to add variation")
	}
	return variation, nil
}
func (c *adminDatabase) AddProducts(ctx context.Context, product domain.Product) (domain.Product, any) {

	// first check the product already exist
	var checkProduct domain.Product
	c.DB.Raw("SELECT * FROM products WHERE product_name=?", product.ProductName).Scan(&checkProduct)

	if checkProduct.ID != 0 {
		return domain.Product{}, "Product already exist"
	}
	querry := `INSERT INTO products (product_name,description,category_id,price,image)VALUES($1,$2,$3,$4,$5) RETURNING product_name,description,category_id,price,image`
	err := c.DB.Raw(querry, product.ProductName, product.Description, product.CategoryID, product.Price, product.Image).Scan(&product).Error

	return product, err
}
