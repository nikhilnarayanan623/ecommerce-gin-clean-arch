package repository

import (
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type adminDatabase struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminDatabase{DB: DB}
}
