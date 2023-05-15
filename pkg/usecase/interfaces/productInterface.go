package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type ProductUseCase interface {
	GetCategory(ctx context.Context) (res.FullCategory, error)
	AddCategory(ctx context.Context, category domain.Category) error
	// variations
	AddVariation(ctx context.Context, vartaion domain.Variation) (domain.Variation, error)
	AddVariationOption(ctx context.Context, variationOption domain.VariationOption) (domain.VariationOption, error)

	// products
	GetProducts(ctx context.Context, pagination req.Pagination) (products []res.Product, err error)
	AddProduct(ctx context.Context, product domain.Product) error
	UpdateProduct(ctx context.Context, product domain.Product) error

	AddProductItem(ctx context.Context, productItem req.ProductItem) error
	GetProductItems(ctx context.Context, productID uint) ([]res.ProductItems, error)

	// offer
	AddOffer(ctx context.Context, offer domain.Offer) error
	RemoveOffer(ctx context.Context, offerID uint) error
	GetAllOffers(ctx context.Context) ([]domain.Offer, error)

	// offer category
	AddOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error
	GetAllOffersOfCategories(ctx context.Context) ([]res.OfferCategory, error)
	RemoveOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error
	ReplaceOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error

	// offer product
	AddOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error
	GetAllOffersOfProducts(ctx context.Context) ([]res.OfferProduct, error)
	RemoveOfferProducts(ctx context.Context, offerProdcts domain.OfferProduct) error
	ReplaceOfferProducts(ctx context.Context, offerProduct domain.OfferProduct) error
}
