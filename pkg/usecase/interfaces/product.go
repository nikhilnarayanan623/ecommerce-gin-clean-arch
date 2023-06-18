package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type ProductUseCase interface {
	FindCategory(ctx context.Context) (response.FullCategory, error)
	SaveCategory(ctx context.Context, category domain.Category) error
	// variations
	SaveVariation(ctx context.Context, variation domain.Variation) error
	SaveVariationOption(ctx context.Context, variationOption domain.VariationOption) error

	// products
	FindAllProducts(ctx context.Context, pagination request.Pagination) (products []response.Product, err error)
	AddProduct(ctx context.Context, product domain.Product) error
	UpdateProduct(ctx context.Context, product domain.Product) error

	AddProductItem(ctx context.Context, productItem request.ProductItem) error
	FindProductItems(ctx context.Context, productID uint) ([]response.ProductItems, error)

	// offer
	AddOffer(ctx context.Context, offer domain.Offer) error
	RemoveOffer(ctx context.Context, offerID uint) error
	FindAllOffers(ctx context.Context) ([]domain.Offer, error)

	// offer category
	AddOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error
	FindAllOffersOfCategories(ctx context.Context) ([]response.OfferCategory, error)
	RemoveOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error
	ReplaceOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error

	// offer product
	AddOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error
	FindAllOffersOfProducts(ctx context.Context) ([]response.OfferProduct, error)
	RemoveOfferProducts(ctx context.Context, offerProdcts domain.OfferProduct) error
	ReplaceOfferProducts(ctx context.Context, offerProduct domain.OfferProduct) error
}
