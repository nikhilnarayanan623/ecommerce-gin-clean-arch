package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type authDatabase struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) interfaces.AuthRepository {
	return &authDatabase{
		DB: db,
	}
}

func (c *authDatabase) SaveRefreshSession(ctx context.Context, refreshSession domain.RefreshSession) error {
	query := `INSERT INTO refresh_sessions (token_id, user_id, refresh_token, expire_at) 
VALUES ($1, $2, $3, $4)`
	err := c.DB.Exec(query, refreshSession.TokenID, refreshSession.UserID, refreshSession.RefreshToken, refreshSession.ExpireAt).Error

	return err
}
func (c *authDatabase) FindRefreshSessionByTokenID(ctx context.Context, tokenID uuid.UUID) (refreshSession domain.RefreshSession, err error) {
	query := `SELECT * FROM refresh_sessions WHERE token_id = $1`

	err = c.DB.Raw(query, tokenID).Scan(&refreshSession).Error

	return
}
