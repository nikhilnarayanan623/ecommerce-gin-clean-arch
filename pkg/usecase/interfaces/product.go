package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type ProductUseCase interface {
	FindAllCategories(ctx context.Context, pagination request.Pagination) ([]response.Category, error)
	SaveCategory(ctx context.Context, categoryName string) error
	SaveSubCategory(ctx context.Context, subCategory request.SubCategory) error

	// variations
	SaveVariation(ctx context.Context, variation request.Variation) error
	SaveVariationOption(ctx context.Context, variationOption request.VariationOption) error

	FindAllVariationsAndItsValues(ctx context.Context, categoryID uint) ([]response.Variation, error)

	// products
	FindAllProducts(ctx context.Context, pagination request.Pagination) (products []response.Product, err error)
	SaveProduct(ctx context.Context, product request.Product) error
	UpdateProduct(ctx context.Context, product domain.Product) error

	SaveProductItem(ctx context.Context, productItem request.ProductItem) error
	FindProductItems(ctx context.Context, productID uint) ([]response.ProductItems, error)

	// offer
	SaveOffer(ctx context.Context, offer request.Offer) error
	RemoveOffer(ctx context.Context, offerID uint) error
	FindAllOffers(ctx context.Context, pagination request.Pagination) ([]domain.Offer, error)

	// offer category
	SaveCategoryOffer(ctx context.Context, offerCategory request.OfferCategory) error
	FindAllCategoryOffers(ctx context.Context, pagination request.Pagination) ([]response.OfferCategory, error)
	RemoveCategoryOffer(ctx context.Context, categoryOfferID uint) error
	ChangeCategoryOffer(ctx context.Context, categoryOfferID, offerID uint) error

	// offer product
	SaveProductOffer(ctx context.Context, offerProduct domain.OfferProduct) error
	FindAllProductOffers(ctx context.Context, pagination request.Pagination) ([]response.OfferProduct, error)
	RemoveProductOffer(ctx context.Context, productOfferID uint) error
	ChangeProductOffer(ctx context.Context, productOfferID, offerID uint) error
}
