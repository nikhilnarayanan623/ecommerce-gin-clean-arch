package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type ProductRepository interface {
	Transactions(ctx context.Context, trxFn func(repo ProductRepository) error) error
	//product
	FindProductByID(ctx context.Context, productID uint) (product domain.Product, err error)
	FindProductByName(ctx context.Context, productName string) (product domain.Product, err error)
	IsProductNameAlreadyExist(ctx context.Context, productName string) (exist bool, err error)

	FindAllProducts(ctx context.Context, pagination request.Pagination) ([]response.Product, error)
	SaveProduct(ctx context.Context, product request.Product) error
	UpdateProduct(ctx context.Context, product domain.Product) error

	// product items
	FindProductItemByID(ctx context.Context, productItemID uint) (domain.ProductItem, error)
	FindAllProductItems(ctx context.Context, productID uint) ([]response.ProductItems, error)
	IsProductItemAlreadyExist(ctx context.Context, productID, variationOptionID uint) (exist bool, err error)
	SaveProductConfiguration(ctx context.Context, productItemID, variationOptionID uint) error
	SaveProductItem(ctx context.Context, productItem domain.ProductItem) (productItemID uint, err error)

	// category
	FindCategoryByName(ctx context.Context, categoryName string) (domain.Category, error)
	FindCategoryByID(ctx context.Context, categoryID uint) (domain.Category, error)

	FindAllMainCategories(ctx context.Context, pagination request.Pagination) (categories []response.Category, err error)
	SaveCategory(ctx context.Context, categoryName string) (err error)

	FindAllSubCategories(ctx context.Context, categoryID uint) (subCategories []response.SubCategory, err error)
	SaveSubCategory(ctx context.Context, categoryID uint, categoryName string) (err error)

	// variation
	SaveVariation(ctx context.Context, categoryID uint, variationName string) error
	FindVariationByID(ctx context.Context, variationID uint) (domain.Variation, error)
	FindVariationByNameAndCategoryID(ctx context.Context,
		variationName string, categoryID uint) (variation domain.Variation, err error)
	FindAllVariationsByCategoryID(ctx context.Context, categoryID uint) ([]response.Variation, error)

	// variation values
	SaveVariationOption(ctx context.Context, variationID uint, variationValue string) error
	FindVariationOptionByValueAndVariationID(ctx context.Context, variationID uint, variationValue string) (domain.VariationOption, error)
	FindAllVariationOptionsByVariationID(ctx context.Context, variationID uint) ([]response.VariationOption, error)

	FindAllVariationValuesOfProductItem(ctx context.Context,
		productItemID uint) (productVariations []response.ProductVariationValue, err error)

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
	FindOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) (domain.OfferCategory, error)
	FindOfferCategoryCategoryID(ctx context.Context, categoryID uint) (domain.OfferCategory, error)
	FindAllOfferCategories(ctx context.Context, pagination request.Pagination) ([]response.OfferCategory, error)

	SaveCategoryOffer(ctx context.Context, categoryOffer request.OfferCategory) (categoryOfferID uint, err error)
	DeleteCategoryOffer(ctx context.Context, categoryOfferID uint) error
	UpdateCategoryOffer(ctx context.Context, categoryOfferID, offerID uint) error

	// offer products
	FindOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) (domain.OfferProduct, error)
	FindAllOfferProducts(ctx context.Context, pagination request.Pagination) ([]response.OfferProduct, error)
	FindOfferProductByProductID(ctx context.Context, productID uint) (domain.OfferProduct, error)

	SaveOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) (productOfferId uint, err error)
	DeleteOfferProduct(ctx context.Context, productOfferID uint) error
	UpdateOfferProduct(ctx context.Context, productOfferID, offerID uint) error

	//new refracted
	DeleteAllProductOffersByOfferID(ctx context.Context, offerID uint) error
	DeleteAllCategoryOffersByOfferID(ctx context.Context, offerID uint) error
}
