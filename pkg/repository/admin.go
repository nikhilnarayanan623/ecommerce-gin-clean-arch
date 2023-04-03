package repository

import (
	"context"
	"errors"
	"fmt"

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

func (c *adminDatabase) SaveAdmin(ctx context.Context, admin domain.Admin) error {

	querry := `INSERT INTO admins (user_name,email,password) VALUES ($1,$2,$3)`
	if c.DB.Exec(querry, admin.UserName, admin.Email, admin.Password).Error != nil {
		return errors.New("faild to save admin")
	}

	return nil
}

func (c *adminDatabase) FindAllUser(ctx context.Context) ([]domain.User, error) {

	var users []domain.User
	err := c.DB.Raw("SELECT * FROM users").Scan(&users).Error

	return users, err
}

func (c *adminDatabase) BlockUser(ctx context.Context, userID uint) error {

	// first check ther user is valid or not
	var user domain.User
	c.DB.Raw("SELECT * FROM users WHERE id=?", userID).Scan(&user)
	if user.Email == "" { // here given id so check with email
		return errors.New("invalid user id user doesn't exist")
	}

	query := `UPDATE users SET block_status = $1 WHERE id = $2`
	if c.DB.Exec(query, !user.BlockStatus, userID).Error != nil {
		return fmt.Errorf("faild update user block_status to %v", !user.BlockStatus)
	}
	return nil
}
