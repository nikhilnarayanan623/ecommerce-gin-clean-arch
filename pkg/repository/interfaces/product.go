package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

type ProductRepository interface {
	Transactions(ctx context.Context, trxFn func(repo ProductRepository) error) error

	// category
	FindCategoryByName(ctx context.Context, categoryName string) (domain.Category, error)
	FindAllMainCategories(ctx context.Context, pagination request.Pagination) ([]response.Category, error)
	SaveCategory(ctx context.Context, categoryName string) error

	// sub category
	FindAllSubCategories(ctx context.Context, categoryID uint) ([]response.SubCategory, error)
	SaveSubCategory(ctx context.Context, categoryID uint, categoryName string) error

	// variation
	SaveVariation(ctx context.Context, categoryID uint, variationName string) error
	FindVariationByNameAndCategoryID(ctx context.Context, variationName string, categoryID uint) (domain.Variation, error)
	FindAllVariationsByCategoryID(ctx context.Context, categoryID uint) ([]response.Variation, error)

	// variation values
	SaveVariationOption(ctx context.Context, variationID uint, variationValue string) error
	FindVariationOptionByValueAndVariationID(ctx context.Context, variationID uint, variationValue string) (domain.VariationOption, error)
	FindAllVariationOptionsByVariationID(ctx context.Context, variationID uint) ([]response.VariationOption, error)

	FindAllVariationValuesOfProductItem(ctx context.Context, productItemID uint) ([]response.ProductVariationValue, error)
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
	FindVariationCountForProduct(ctx context.Context, productID uint) (variationCount uint, err error) // to check the product config already exist
	FindAllProductItemIDsByProductIDAndVariationOptionID(ctx context.Context, productID, variationOptionID uint) ([]uint, error)
	SaveProductConfiguration(ctx context.Context, productItemID, variationOptionID uint) error
	SaveProductItem(ctx context.Context, productItem domain.ProductItem) (productItemID uint, err error)
	// product item image
	FindAllProductItemImages(ctx context.Context, productItemID uint) (images []string, err error)
	SaveProductItemImage(ctx context.Context, productItemID uint, image string) error
}
