package db

import (
	"fmt"
	"log"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func to connect data base using config(database config) and return address of a new instnce of gorm DB
func ConnectDatbase(cfg config.Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatal("Faild to connect with database")
		return nil, err
	}

	// migrate the database tables
	db.AutoMigrate(
		//user
		domain.User{},
		domain.Country{},
		domain.Address{},
		domain.UserAddress{},

		//admin
		domain.Admin{},

		//product
		domain.Category{},
		domain.Product{},
		domain.Variation{},
		domain.VariationOption{},
		domain.ProductItem{},
		domain.ProductConfiguration{},
		domain.ProductImage{},

		//wish list
		domain.WishList{},

		//cart
		domain.Cart{},
		domain.CartItem{},
	)

	return db, err
}
