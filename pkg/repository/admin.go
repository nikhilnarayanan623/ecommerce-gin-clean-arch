package repository

import (
	"context"
	"errors"

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

func (c *adminDatabase) FindAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, error) {

	if c.DB.Raw("SELECT * FROM admins WHERE email=? OR user_name=?", admin.Email, admin.UserName).Scan(&admin).Error != nil {
		return admin, errors.New("faild to find admin")
	} 

	return admin, nil
}

func (c *adminDatabase) SaveAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, error) {

	// first check the admin already exist or not
	var dbAdmin domain.Admin
	c.DB.Raw("SELECT * FROM admins WHERE email=? OR user_name=?", admin.Email, admin.UserName).Scan(&dbAdmin)

	if dbAdmin.ID != 0 { // amdmin already exist

		// first check the email is already exist
		if dbAdmin.Email == admin.Email {
			// errMap["email"] = "Email Already exist"
			return admin, errors.New("admin already exist with this email")
		} // if email not then its user name is exist
		return admin, errors.New("admin already exist with this user_name")
	}

	//if admin not exist then create it
	querry := `INSERT INTO admins (user_name,email,password) VALUES ($1,$2,$3) RETURNING id,user_name,email,password`
	if c.DB.Raw(querry, admin.UserName, admin.Email, admin.Password).Scan(&admin).Error != nil {
		return admin, errors.New("faild to create account for admin")
	}

	//c.DB.Raw("SELECT * FROM admins WHERE email=?", admin.Email).Scan(&admin)
	return admin, nil // successfully admin added
}

func (c *adminDatabase) FindAllUser(ctx context.Context) ([]domain.User, error) {

	var users []domain.User
	err := c.DB.Raw("SELECT * FROM users").Scan(&users).Error

	return users, err
}

func (c *adminDatabase) BlockUser(ctx context.Context, user domain.User) (domain.User, error) {

	// first check ther user is valid or not
	c.DB.Raw("SELECT * FROM users WHERE id=?", user.ID).Scan(&user)

	if user.Email == "" { // here given id so check with email
		return user, errors.New("invalid user id user doesn't exist")
	}

	// if user is blocked then unblock
	if user.BlockStatus {
		c.DB.Raw("UPDATE users SET block_status='F' WHERE id=? RETURNING block_status", user.ID).Scan(&user)
		//user.BlockStatus = false
		return user, nil
	}
	c.DB.Raw("UPDATE users SET block_status='T' WHERE id=? RETURNING block_status", user.ID).Scan(&user)
	//user.BlockStatus = true
	return user, nil
}
