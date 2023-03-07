package repository

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type adminDatabase struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminDatabase{DB: DB}
}

func (c adminDatabase) FindAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, any) {

	c.DB.Raw("SELECT * FROM admins WHERE email=? OR user_name=?", admin.Email, admin.UserName).Scan(&admin)

	//check the admin got or not
	if admin.ID == 0 {
		return admin, map[string]string{"msg": "Can't find the admin"}
	}

	return admin, nil
}

func (c *adminDatabase) FindAllUser(ctx context.Context) ([]domain.Users, error) {

	var users []domain.Users
	err := c.DB.Raw("SELECT * FROM users").Scan(&users).Error

	return users, err
}

func (c adminDatabase) AddCategory(ctx context.Context, category domain.Category) (domain.Category, any) {

	var checkCat domain.Category
	//first check the categoryname already exisits or not
	c.DB.Raw("SELECT * FROM categories WHERE category_name=?", category.CategoryName).Scan(&checkCat)

	if checkCat.ID != 0 { // means category already exist
		return checkCat, map[string]string{"error": "category already exist"}
	}

	// check the given category is main or sub
	if category.CategoryID == 0 { // no catogry id means its main category
		querry := `INSERT INTO categories (category_name)VALUES($1) RETURNING category_name`
		c.DB.Raw(querry, category.CategoryName).Scan(&category)
	} else {
		// first check the category id is valid or not
		c.DB.Raw("SELECT * FROM categories WHERE id=?", category.CategoryID).Scan(&checkCat)
		if checkCat.ID == 0 { // its not a valid category
			return checkCat, map[string]string{"error": "category_id is not valid means provided main category is not valid"}
		}
		//otherwise add its with main category
		querry := `INSERT INTO categories (category_id,category_name)VALUES($1,$2) RETURNING category_id,category_name`
		c.DB.Raw(querry, category.CategoryID, category.CategoryName).Scan(&category)
	}

	return category, nil
}
