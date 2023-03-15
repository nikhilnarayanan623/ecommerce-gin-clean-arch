package usecase

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type productUseCase struct {
	productRepo interfaces.ProductRepository
}

// to get a new instance of productUseCase
func NewProductUseCase(productRepo interfaces.ProductRepository) service.ProductUseCase {
	return &productUseCase{productRepo: productRepo}
}

// to get all Category , all variation , all variation value
func (c *productUseCase) GetCategory(ctx context.Context) (helper.RespFullCategory, error) {
	return c.productRepo.GetCategory(ctx)
}

// to add a new category or add new sub category
func (c *productUseCase) AddCategory(ctx context.Context, category domain.Category) (domain.Category, error) {

	return c.productRepo.AddCategory(ctx, category)
}

// to add new variation for a category
func (c *productUseCase) AddVariation(ctx context.Context, vartaion domain.Variation) (domain.Variation, error) {

	return c.productRepo.AddVariation(ctx, vartaion)
}

// to add new variation value for varation
func (c *productUseCase) AddVariationOption(ctx context.Context, variationOption domain.VariationOption) (domain.VariationOption, error) {
	return c.productRepo.AddVariationOption(ctx, variationOption)
}

// to get all product
func (c *productUseCase) GetProducts(ctx context.Context) ([]helper.ResponseProduct, error) {
	return c.productRepo.GetProducts(ctx)
}

// to add new product
func (c *productUseCase) AddProducts(ctx context.Context, product domain.Product) (domain.Product, error) {
	return c.productRepo.AddProducts(ctx, product)
}

// for add new productItem for a speicific product
func (c *productUseCase) AddProductItem(ctx context.Context, productItem helper.ReqProductItem) (domain.ProductItem, error) {
	return c.productRepo.AddProductItem(ctx, productItem)
}

// for get all productItem for a specific product
func (c *productUseCase) GetProductItems(ctx context.Context, product domain.Product) ([]helper.RespProductItems, error) {
	return c.productRepo.GetProductItems(ctx, product)
}
