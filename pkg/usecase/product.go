package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
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
func (c *productUseCase) GetCategory(ctx context.Context) (res.RespFullCategory, error) {
	return c.productRepo.GetCategory(ctx)
}

// to add a new category or add new sub category
func (c *productUseCase) AddCategory(ctx context.Context, category domain.Category) error {

	// check the given category already exist or not
	var checkCategory = domain.Category{CategoryName: category.CategoryName}
	if checkCategory, err := c.productRepo.FindCategory(ctx, category); err != nil {
		return err
	} else if checkCategory.ID != 0 {
		return fmt.Errorf("category already exit with %s name", category.CategoryName)
	}

	// if main category is given then check it valid or not
	if category.CategoryID != 0 {
		checkCategory.CategoryName = ""
		checkCategory.ID = category.CategoryID
		if checkCategory, err := c.productRepo.FindCategory(ctx, checkCategory); err != nil {
			return err
		} else if checkCategory.CategoryName == "" {
			return errors.New("invalid categoyr id")
		}
	}

	return c.productRepo.SaveCategory(ctx, category)
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
func (c *productUseCase) GetProducts(ctx context.Context) ([]res.ResponseProduct, error) {
	return c.productRepo.GetProducts(ctx)
}

// to add new product
func (c *productUseCase) AddProduct(ctx context.Context, product domain.Product) error {
	//check product already exist or not
	if product, err := c.productRepo.FindProduct(ctx, product); err != nil {
		return err
	} else if product.ID != 0 {
		return fmt.Errorf("product already exist with %s product name", product.ProductName)
	}
	return c.productRepo.SaveProduct(ctx, product)
}

// for add new productItem for a speicific product
func (c *productUseCase) AddProductItem(ctx context.Context, productItem req.ReqProductItem) (domain.ProductItem, error) {
	return c.productRepo.AddProductItem(ctx, productItem)
}

// for get all productItem for a specific product
func (c *productUseCase) GetProductItems(ctx context.Context, product domain.Product) ([]res.RespProductItems, error) {
	return c.productRepo.GetProductItems(ctx, product)
}
