package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
)

type UserRepository interface {
	FindUser(ctx context.Context, user domain.Users) (domain.Users, error)
	SaveUser(ctx context.Context, user domain.Users) (domain.Users, error)

	GetCartItems(ctx context.Context, userId uint) (helper.ResponseCart, error)
}
