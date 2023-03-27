package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

type ProductRepository interface {
	FindProduct(ctx context.Context, product domain.Product) (domain.Product, error)
	GetProducts(ctx context.Context) ([]res.ResponseProduct, error)
	SaveProduct(ctx context.Context, product domain.Product) error

	// product items
	GetProductItems(ctx context.Context, product domain.Product) ([]res.RespProductItems, error)
	AddProductItem(ctx context.Context, productItem req.ReqProductItem) (domain.ProductItem, error)

	// category
	FindCategory(ctx context.Context, category domain.Category) (domain.Category, error)
	FindAllCategories(ctx context.Context) ([]res.Category, error)
	SaveCategory(ctx context.Context, productCategory domain.Category) error

	// variation
	FindAllVariations(ctx context.Context) ([]res.VariationName, error)
	AddVariation(ctx context.Context, variation domain.Variation) (domain.Variation, error)

	// variation values
	AddVariationOption(ctx context.Context, variationOption domain.VariationOption) (domain.VariationOption, error)
	FindAllVariationValues(ctx context.Context) ([]res.VariationValue, error)

	// offer
	FindOffer(ctx context.Context, offer domain.Offer) (domain.Offer, error)
	FindAllOffer(ctx context.Context) ([]domain.Offer, error)
	SaveOffer(ctx context.Context, offer domain.Offer) error
	DeleteOffer(ctx context.Context, offerID uint) error

	// offer category
	FindOfferCategory(ctx context.Context, offerCateogy domain.OfferCategory) (domain.OfferCategory, error)
	FindOfferCategoryCategoryID(ctx context.Context, categoryID uint) (domain.OfferCategory, error)
	FindAllOfferCategories(ctx context.Context) ([]res.ResOfferCategory, error)
	SaveOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error
	DeleteOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error
	UpdateOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error
	UpdateProductPriceForCategory(ctx context.Context) error

	// offer productsss
	FindOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) (domain.OfferProduct, error)
	FindAllOfferProducts(ctx context.Context) ([]res.ResOfferProduct, error)
	FindOfferProductByProductID(ctx context.Context, productID uint) (domain.OfferProduct, error)

	SaveOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error
	DeleteOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error
	UpdateOfferProduct(ctx context.Context, productOffer domain.OfferProduct) error
	UpdateProductPriceForProductOffer(ctx context.Context) error
}
