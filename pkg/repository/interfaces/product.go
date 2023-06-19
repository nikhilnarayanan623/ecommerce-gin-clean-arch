package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type ProductRepository interface {
	//product
	FindProductByID(ctx context.Context, productID uint) (product domain.Product, err error)
	FindProduct(ctx context.Context, product domain.Product) (domain.Product, error)

	FindAllProducts(ctx context.Context, pagination request.Pagination) ([]response.Product, error)
	SaveProduct(ctx context.Context, product domain.Product) error
	UpdateProduct(ctx context.Context, product domain.Product) error

	// product items
	FindProductItem(ctx context.Context, productItemID uint) (domain.ProductItem, error)
	FindAllProductItems(ctx context.Context, productID uint) ([]response.ProductItems, error)
	//FindAllProductItemImages(ctx context.Context, productItemID uint) (images []string, err error)
	SaveProductItem(ctx context.Context, productItem request.ProductItem) error

	// category
	FindCategoryByName(ctx context.Context, categoryName string) (domain.Category, error)
	FindCategoryByID(ctx context.Context, categoryID uint) (domain.Category, error)

	FindAllCategories(ctx context.Context, pagination request.Pagination) ([]response.Category, error)

	SaveCategory(ctx context.Context, categoryName string) (err error)
	SaveSubCategory(ctx context.Context, categoryID uint, categoryName string) (err error)

	// variation
	SaveVariation(ctx context.Context, variation request.Variation) error
	FindVariationByID(ctx context.Context, variationID uint) (domain.Variation, error)
	FindVariationByNameAndCategoryID(ctx context.Context,
		variationName string, categoryID uint) (variation domain.Variation, err error)
	FindAllVariationsByCategoryID(ctx context.Context, categoryID uint) ([]response.Variation, error)

	// variation values
	SaveVariationOption(ctx context.Context, variationOption request.VariationOption) error
	FindVariationOptionByValueAndVariationID(ctx context.Context,
		variationOptionValue string, categoryID uint) (variationOption domain.VariationOption, err error)
	FindAllVariationOptionsByVariationID(ctx context.Context, variationID uint) ([]response.VariationOption, error)

	// offer
	FindOffer(ctx context.Context, offer domain.Offer) (domain.Offer, error)
	FindAllOffer(ctx context.Context) ([]domain.Offer, error)
	SaveOffer(ctx context.Context, offer domain.Offer) error
	DeleteOffer(ctx context.Context, offerID uint) error

	// offer discount price update
	UpdateDiscountPrice(ctx context.Context) error

	// offer category
	FindOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) (domain.OfferCategory, error)
	FindOfferCategoryCategoryID(ctx context.Context, categoryID uint) (domain.OfferCategory, error)
	FindAllOfferCategories(ctx context.Context) ([]response.OfferCategory, error)

	SaveOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error
	DeleteOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error
	UpdateOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error

	// offer products
	FindOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) (domain.OfferProduct, error)
	FindAllOfferProducts(ctx context.Context) ([]response.OfferProduct, error)
	FindOfferProductByProductID(ctx context.Context, productID uint) (domain.OfferProduct, error)

	SaveOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error
	DeleteOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error
	UpdateOfferProduct(ctx context.Context, productOffer domain.OfferProduct) error
}
