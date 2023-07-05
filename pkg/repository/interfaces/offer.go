package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

type OfferRepository interface {
	Transactions(ctx context.Context, trxFn func(repo OfferRepository) error) error

	// offer
	FindOfferByID(ctx context.Context, offerID uint) (domain.Offer, error)
	FindOfferByName(ctx context.Context, offerName string) (domain.Offer, error)
	FindAllOffers(ctx context.Context, pagination request.Pagination) ([]domain.Offer, error)
	SaveOffer(ctx context.Context, offer request.Offer) error
	DeleteOffer(ctx context.Context, offerID uint) error

	// to calculate the discount price and update
	UpdateProductsDiscountByCategoryOfferID(ctx context.Context, categoryOfferID uint) error
	UpdateProductItemsDiscountByCategoryOfferID(ctx context.Context, categoryOfferID uint) error
	UpdateProductsDiscountByProductOfferID(ctx context.Context, productOfferID uint) error
	UpdateProductItemsDiscountByProductOfferID(ctx context.Context, productOfferID uint) error

	// to remove the discount product price
	RemoveProductsDiscountByCategoryOfferID(ctx context.Context, categoryOfferID uint) error
	RemoveProductItemsDiscountByCategoryOfferID(ctx context.Context, categoryOfferID uint) error
	RemoveProductsDiscountByProductOfferID(ctx context.Context, productOfferID uint) error
	RemoveProductItemsDiscountByProductOfferID(ctx context.Context, productOfferID uint) error

	// offer category
	FindOfferCategoryCategoryID(ctx context.Context, categoryID uint) (domain.OfferCategory, error)
	FindAllOfferCategories(ctx context.Context, pagination request.Pagination) ([]response.OfferCategory, error)

	SaveCategoryOffer(ctx context.Context, categoryOffer request.OfferCategory) (categoryOfferID uint, err error)
	DeleteCategoryOffer(ctx context.Context, categoryOfferID uint) error
	UpdateCategoryOffer(ctx context.Context, categoryOfferID, offerID uint) error

	// offer products
	FindOfferProductByProductID(ctx context.Context, productID uint) (domain.OfferProduct, error)
	FindAllOfferProducts(ctx context.Context, pagination request.Pagination) ([]response.OfferProduct, error)

	SaveOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) (productOfferId uint, err error)
	DeleteOfferProduct(ctx context.Context, productOfferID uint) error
	UpdateOfferProduct(ctx context.Context, productOfferID, offerID uint) error

	DeleteAllProductOffersByOfferID(ctx context.Context, offerID uint) error
	DeleteAllCategoryOffersByOfferID(ctx context.Context, offerID uint) error
}
