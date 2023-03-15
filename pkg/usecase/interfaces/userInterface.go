package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
)

type UserUseCase interface {
	Signup(ctx context.Context, user domain.Users) (domain.Users, error)
	Login(ctx context.Context, user domain.Users) (domain.Users, error)
	LoginOtp(ctx context.Context, user domain.Users) (domain.Users, error)

	Home(ctx context.Context, userId uint) (domain.Users, error)
	GetCartItems(ctx context.Context, userId uint) (helper.ResponseCart, any)
}
