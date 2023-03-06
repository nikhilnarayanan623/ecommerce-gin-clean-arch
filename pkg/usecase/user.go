package usecase

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type userUserCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(repo interfaces.UserRepository) service.UserUseCase {
	return &userUserCase{userRepo: repo}
}

func (c *userUserCase) FindAllUser(ctx context.Context) ([]domain.Users, error) {

	users, err := c.userRepo.FindAllUser(ctx)

	if err != nil {
		// do something
	}
	//logic

	return users, err
}

func (c *userUserCase) FindUserByID(ctx context.Context, id uint) (domain.Users, error) {

	user, err := c.userRepo.FindUserByID(ctx, id)

	if err != nil {

	}

	return user, err

}

func (c *userUserCase) SaveUser(ctx context.Context, user domain.Users) (domain.Users, error) {

	user, err := c.userRepo.SaveUser(ctx, user)

	if err != nil {

	}

	return user, err
}

func (c *userUserCase) DeleteUser(ctx context.Context, user domain.Users) error {

	err := c.userRepo.DeleteUser(ctx, user)

	if err != nil {

	}

	return err
}
