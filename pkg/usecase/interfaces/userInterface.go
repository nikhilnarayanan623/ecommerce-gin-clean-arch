package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

type UserUseCase interface {
	FindAllUser(ctx context.Context) ([]domain.Users, error)
	FindUserByID(ctx context.Context, id uint) (domain.Users, error)
	SaveUser(ctx context.Context, user domain.Users) (domain.Users, error)
	DeleteUser(ctx context.Context, user domain.Users) error
}
