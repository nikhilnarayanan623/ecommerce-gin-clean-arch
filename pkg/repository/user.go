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

func (c *userDatabse) FindAllUser(ctx context.Context) ([]domain.Users, error) {
	var users []domain.Users
	err := c.DB.Raw("SELECT * FROM users").Scan(&users).Error

	return users, err
}

func (c *userDatabse) FindUserByID(ctx context.Context, id uint) (domain.Users, error) {

	var user domain.Users
	err := c.DB.Raw("SELECT * FROM user where id=?", id).Scan(user).Error

	return user, err
}

func (c *userDatabse) SaveUser(ctx context.Context, user domain.Users) (domain.Users, error) {

	err := c.DB.Save(&user).Error

	return user, err
}

func (c *userDatabse) DeleteUser(ctx context.Context, user domain.Users) error {

	return c.DB.Raw("DELETE users WHERE id=?", user.ID).Error
}
