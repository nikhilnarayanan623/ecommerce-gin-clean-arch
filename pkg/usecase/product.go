package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
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
	var (
		resFullCategories res.RespFullCategory
		err               error
	)
	// find all categories
	if resFullCategories.Category, err = c.productRepo.FindAllCategories(ctx); err != nil {
		return resFullCategories, err
	}

	// find all variations
	if resFullCategories.VariationName, err = c.productRepo.FindAllVariations(ctx); err != nil {
		return resFullCategories, err
	}

	// find all variation values
	if resFullCategories.VariationValue, err = c.productRepo.FindAllVariationValues(ctx); err != nil {
		return resFullCategories, err
	}
	return resFullCategories, nil
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
	return c.productRepo.FindAllProducts(ctx)
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
func (c *productUseCase) GetProductItems(ctx context.Context, productID uint) ([]res.RespProductItems, error) {

	//validate the productID
	if product, err := c.productRepo.FindProduct(ctx, domain.Product{ID: productID}); err != nil {
		return []res.RespProductItems{}, err
	} else if product.ID == 0 {
		return []res.RespProductItems{}, errors.New("invalid product_id")
	}

	return c.productRepo.FindAllProductItems(ctx, productID)
}

func (c *productUseCase) UpdateProduct(ctx context.Context, product domain.Product) error {
	// validate the product_id
	if product, err := c.productRepo.FindProduct(ctx, product); err != nil {

	} else if product.ProductName == "" {
		return errors.New("invalid product_id")
	}

	// check the given product_name already exist or not
	if checkProduct, err := c.productRepo.FindProduct(ctx, domain.Product{ProductName: product.ProductName}); err != nil {
		return err
	} else if checkProduct.ID != 0 && checkProduct.ID != product.ID {
		return errors.New("can't update the product \nthere is alread a product exist with this product_name")
	}

	return c.productRepo.UpdateProduct(ctx, product)
}
