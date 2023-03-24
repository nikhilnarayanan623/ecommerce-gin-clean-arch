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

	GetProductItems(ctx context.Context, product domain.Product) ([]res.RespProductItems, error)
	AddProductItem(ctx context.Context, productItem req.ReqProductItem) (domain.ProductItem, error)

	FindCategory(ctx context.Context, category domain.Category) (domain.Category, error)
	GetCategory(ctx context.Context) (res.RespFullCategory, error)
	SaveCategory(ctx context.Context, productCategory domain.Category) error

	AddVariation(ctx context.Context, variation domain.Variation) (domain.Variation, error)
	AddVariationOption(ctx context.Context, variationOption domain.VariationOption) (domain.VariationOption, error)
	// offer
	FindOffer(ctx context.Context, offer domain.Offer) (domain.Offer, error)
	FindAllOffer(ctx context.Context) ([]domain.Offer, error)
	SaveOffer(ctx context.Context, offer domain.Offer) error

	FindOfferCategoryCategoryID(ctx context.Context, categoryID uint) (domain.OfferCategory, error)
	FindAllOfferCategories(ctx context.Context) ([]res.ResOfferCategory, error)
	SaveOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error
	UpdateOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error

	FindOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) (domain.OfferProduct, error)
	FindAllOfferProducts(ctx context.Context) ([]res.ResOfferProduct, error)
	SaveOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error
	//DeleteOfferProducts(ctx context.Context, offerID uint) error
}
