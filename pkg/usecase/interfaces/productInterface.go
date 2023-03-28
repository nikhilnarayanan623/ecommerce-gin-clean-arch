package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

type ProductUseCase interface {
	GetCategory(ctx context.Context) (res.RespFullCategory, error)
	AddCategory(ctx context.Context, category domain.Category) error

	AddVariation(ctx context.Context, vartaion domain.Variation) (domain.Variation, error)
	AddVariationOption(ctx context.Context, variationOption domain.VariationOption) (domain.VariationOption, error)

	GetProducts(ctx context.Context) ([]res.ResponseProduct, error)
	AddProduct(ctx context.Context, product domain.Product) error

	AddProductItem(ctx context.Context, productItem req.ReqProductItem) (domain.ProductItem, error)
	GetProductItems(ctx context.Context, product domain.Product) ([]res.RespProductItems, error)

	// offer
	AddOffer(ctx context.Context, offer domain.Offer) error
	RemoveOffer(ctx context.Context, offerID uint) error
	GetAllOffers(ctx context.Context) ([]domain.Offer, error)

	// offer category
	AddOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error
	GetAllOffersOfCategories(ctx context.Context) ([]res.ResOfferCategory, error)
	RemoveOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error
	ReplaceOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error

	// offer product
	AddOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error
	GetAllOffersOfProducts(ctx context.Context) ([]res.ResOfferProduct, error)
	RemoveOfferProducts(ctx context.Context, offerProdcts domain.OfferProduct) error
	ReplaceOfferProducts(ctx context.Context, offerProduct domain.OfferProduct) error
}
