package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

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
func (c *productUseCase) GetCategory(ctx context.Context) (res.FullCategory, error) {
	var (
		resFullCategories res.FullCategory
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
	err := c.productRepo.SaveCategory(ctx, category)
	if err != nil {
		return err
	}

	return nil
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
func (c *productUseCase) GetProducts(ctx context.Context, pagination req.Pagination) (products []res.Product, err error) {
	return c.productRepo.FindAllProducts(ctx, pagination)
}

// to add new product
func (c *productUseCase) AddProduct(ctx context.Context, product domain.Product) error {
	//check product already exist or not
	if product, err := c.productRepo.FindProduct(ctx, product); err != nil {
		return err
	} else if product.ID != 0 {
		return fmt.Errorf("product already exist with %s product name", product.ProductName)
	}
	log.Printf("successfully product saved\n\n")
	return c.productRepo.SaveProduct(ctx, product)
}

// for add new productItem for a speicific product
func (c *productUseCase) AddProductItem(ctx context.Context, productItem req.ProductItem) error {

	// validate the product_id
	product, err := c.productRepo.FindProductByID(ctx, productItem.ProductID)
	if err != nil {
		return err
	} else if product.ID == 0 {
		return fmt.Errorf("invalid product_id %v", productItem.ProductID)
	}

	// save the product item
	err = c.productRepo.SaveProductItem(ctx, productItem)
	if err != nil {
		return err
	}

	log.Printf("successfully product_item saved for product_id %v\n\n", productItem.ProductID)
	return nil
}

// for get all productItem for a specific product
func (c *productUseCase) GetProductItems(ctx context.Context, productID uint) (productItems []res.ProductItems, err error) {

	//validate the productID
	if product, err := c.productRepo.FindProduct(ctx, domain.Product{ID: productID}); err != nil {
		return productItems, err
	} else if product.ProductName == "" {
		return productItems, errors.New("invalid product_id")
	}

	productItems, err = c.productRepo.FindAllProductItems(ctx, productID)
	if err != nil {
		return productItems, err
	}

	// for _, productItem := range productItems {
	// 	images, err := c.productRepo.FindAllProductItemImages(ctx, productItem.ID)

	// 	if err != nil {
	// 		return productItems, err
	// 	}
	// 	fmt.Println(images, "images")
	// }

	log.Printf("successfully got all prouctItems for product_id %v", productID)
	return productItems, nil
}

func (c *productUseCase) UpdateProduct(ctx context.Context, product domain.Product) error {
	// validate the product_id
	checkProduct, err := c.productRepo.FindProductByID(ctx, product.ID)
	if err != nil {
		return err
	} else if checkProduct.ProductName == "" {
		return errors.New("invalid product_id")
	}

	// check the given product_name already exist or not
	checkProduct, err = c.productRepo.FindProduct(ctx, domain.Product{ProductName: product.ProductName})
	if err != nil {
		return err
	} else if checkProduct.ID != 0 && checkProduct.ID != product.ID {
		return errors.New("can't update the product \nthere is alread a product exist with this product_name")
	}

	return c.productRepo.UpdateProduct(ctx, product)
}
