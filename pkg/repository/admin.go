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
		return admin, errors.New("admin not exist")
	}

	return admin, nil
}

func (c *adminDatabase) SaveAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, any) {

	// first check the admin already exist or not
	var dbAdmin domain.Admin
	c.DB.Raw("SELECT * FROM admins WHERE email=? OR user_name=?", admin.Email, admin.UserName).Scan(&dbAdmin)

	if dbAdmin.ID != 0 { // amdmin already exist

		//errMap := map[string]string{}
		// first check the email is already exist
		if dbAdmin.Email == admin.Email {
			// errMap["email"] = "Email Already exist"
			return admin, helper.SingleRespStruct{Error: "Email Already exist"}
		} // if email not then its user name is exist
		//errMap["user_name"] = "UserName Already exist"
		return admin, helper.SingleRespStruct{Error: "UserName Already exist"}
	}

	//if admin not exist then create it
	querry := `INSERT INTO admins (user_name,email,password) VALUES ($1,$2,$3) RETURNING user_name,email,password`
	c.DB.Raw(querry, admin.UserName, admin.Email, admin.Password).Scan(&admin)

	return admin, nil // successfully admin added
}

func (c *adminDatabase) FindAllUser(ctx context.Context) ([]domain.Users, error) {

	var users []domain.Users
	err := c.DB.Raw("SELECT * FROM users").Scan(&users).Error

	return users, err
}

func (c *adminDatabase) BlockUser(ctx context.Context, user domain.Users) (domain.Users, any) {

	var dbUser domain.Users
	// first check ther user valid or not
	c.DB.Raw("SELECT * FROM users WHERE id=?", user.ID).Scan(&dbUser)

	if dbUser.ID == 0 {
		return user, helper.SingleRespStruct{Error: "Invalid user ID user doesn't exist"}
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
	querry := `SELECT f.id, f.category_name,f.category_id, s.category_name as main_category_name FROM categories f LEFT JOIN categories s ON f.category_id=s.id`
	err := c.DB.Raw(querry).Scan(&response).Error
	return response, err
}

func (c *adminDatabase) AddCategory(ctx context.Context, category domain.Category) (helper.RespCategory, any) {

	var checkCat domain.Category
	//first check the categoryname already exisits or not
	c.DB.Raw("SELECT * FROM categories WHERE category_name=?", category.CategoryName).Scan(&checkCat)

	if checkCat.ID != 0 { // means category already exist
		return helper.RespCategory{}, map[string]string{"error": "category already exist"}
	}

	// check the given category is main or sub
	if category.CategoryID == 0 { // no catogry id means its main category
		querry := `INSERT INTO categories (category_name)VALUES($1) RETURNING category_name`
		c.DB.Raw(querry, category.CategoryName).Scan(&category)
	} else {
		// first check the category id is valid or not
		c.DB.Raw("SELECT * FROM categories WHERE id=?", category.CategoryID).Scan(&checkCat)
		if checkCat.ID == 0 { // its not a valid category
			return helper.RespCategory{}, map[string]string{"error": "category_id is not valid means provided main category is not valid"}
		}
		//otherwise add its with main category
		querry := `INSERT INTO categories (category_id,category_name)VALUES($1,$2) RETURNING category_id,category_name`
		c.DB.Raw(querry, category.CategoryID, category.CategoryName).Scan(&category)
	}

	var response helper.RespCategory

	return response, nil
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
