package usecase

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type productUseCase struct {
	productRepo interfaces.ProductRepository
}

// to get a new instance of productUseCase
func NewProductUseCase(productRepo interfaces.ProductRepository) service.ProductUseCase {
	return &productUseCase{productRepo: productRepo}
}

func (c *productUseCase) FindAllCategories(ctx context.Context, pagination request.Pagination) ([]response.Category, error) {

	categories, err := c.productRepo.FindAllMainCategories(ctx, pagination)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed find all main categories")
	}

	for i, category := range categories {

		subCategory, err := c.productRepo.FindAllSubCategories(ctx, category.ID)
		if err != nil {
			return nil, utils.PrependMessageToError(err, "failed to find sub categories")
		}
		categories[i].SubCategory = subCategory
	}

	return categories, nil
}

// Save category
func (c *productUseCase) SaveCategory(ctx context.Context, categoryName string) error {

	existCategory, err := c.productRepo.FindCategoryByName(ctx, categoryName)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to check category already exist")
	}
	if existCategory.ID != 0 {
		return ErrCategoryAlreadyExist
	}

	err = c.productRepo.SaveCategory(ctx, categoryName)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to save category")
	}

	return nil
}

// Save Sub category
func (c *productUseCase) SaveSubCategory(ctx context.Context, subCategory request.SubCategory) error {

	existSubCat, err := c.productRepo.FindCategoryByName(ctx, subCategory.Name)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to check sub category already exist")
	}
	if existSubCat.ID != 0 {
		return ErrCategoryAlreadyExist
	}

	err = c.productRepo.SaveSubCategory(ctx, subCategory.CategoryID, subCategory.Name)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to save sub category")
	}

	return nil
}

// to add new variation for a category
func (c *productUseCase) SaveVariation(ctx context.Context, categoryID uint, variationNames []string) error {

	err := c.productRepo.Transactions(ctx, func(repo interfaces.ProductRepository) error {

		for _, variationName := range variationNames {

			existVariation, err := c.productRepo.FindVariationByNameAndCategoryID(ctx, variationName, categoryID)
			if err != nil {
				return utils.PrependMessageToError(err, "failed to check variation already exist")
			}

			if existVariation.ID != 0 {
				return utils.PrependMessageToError(ErrVariationAlreadyExist, "variation name "+variationName)
			}

			err = c.productRepo.SaveVariation(ctx, categoryID, variationName)
			if err != nil {
				return utils.PrependMessageToError(err, "failed to save variation")
			}
		}
		return nil
	})

	return err
}

// to add new variation value for variation
func (c *productUseCase) SaveVariationOption(ctx context.Context, variationID uint, variationOptionValues []string) error {

	err := c.productRepo.Transactions(ctx, func(repo interfaces.ProductRepository) error {
		for _, variationValue := range variationOptionValues {

			existVarOpt, err := c.productRepo.FindVariationOptionByValueAndVariationID(ctx, variationID, variationValue)
			if err != nil {
				return utils.PrependMessageToError(err, "failed to check variation already exist")
			}
			if existVarOpt.ID != 0 {
				return utils.PrependMessageToError(ErrVariationOptionAlreadyExist, "variation option value "+variationValue)
			}

			err = c.productRepo.SaveVariationOption(ctx, variationID, variationValue)
			if err != nil {
				return utils.PrependMessageToError(err, "failed to save variation option")
			}
		}
		return nil
	})

	return err
}

func (c *productUseCase) FindAllVariationsAndItsValues(ctx context.Context, categoryID uint) ([]response.Variation, error) {

	variations, err := c.productRepo.FindAllVariationsByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed to find all variations of category")
	}

	// get all variation values of each variations
	for i, variation := range variations {

		variationOption, err := c.productRepo.FindAllVariationOptionsByVariationID(ctx, variation.ID)
		if err != nil {
			return nil, utils.PrependMessageToError(err, "failed to get variation option")
		}
		variations[i].VariationOptions = variationOption
	}
	return variations, nil
}

// to get all product
func (c *productUseCase) FindAllProducts(ctx context.Context, pagination request.Pagination) (products []response.Product, err error) {
	return c.productRepo.FindAllProducts(ctx, pagination)
}

// to add new product
func (c *productUseCase) SaveProduct(ctx context.Context, product request.Product) error {

	exist, err := c.productRepo.IsProductNameAlreadyExist(ctx, product.Name)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to check product name already exist")
	}
	if exist {
		return utils.PrependMessageToError(ErrProductAlreadyExist, "product name "+product.Name)
	}

	err = c.productRepo.SaveProduct(ctx, product)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to save product")
	}
	return nil
}

// for add new productItem for a specific product
func (c *productUseCase) SaveProductItem(ctx context.Context, productID uint, productItem request.ProductItem) error {

	// check the given all combination already exist (Color:Red with Size:M)
	productItemExist, err := c.isAllVariationCombinationExist(productID, productItem.VariationOptionIDs)
	if err != nil {
		return err
	}
	if productItemExist {
		return ErrProductItemAlreadyExist
	}

	err = c.productRepo.Transactions(ctx, func(trxRepo interfaces.ProductRepository) error {

		sku := utils.GenerateSKU()
		newProductItem := domain.ProductItem{
			ProductID:  productID,
			QtyInStock: productItem.QtyInStock,
			Price:      productItem.Price,
			SKU:        sku,
		}

		productItemID, err := trxRepo.SaveProductItem(ctx, newProductItem)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to save product item")
		}

		// save all product configurations based on given variation option id
		for _, variationOptionID := range productItem.VariationOptionIDs {

			err = trxRepo.SaveProductConfiguration(ctx, productItemID, variationOptionID)
			if err != nil {
				return utils.PrependMessageToError(err, "failed to save product_item configuration")
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// To check all variation option is exist for the product
func (c *productUseCase) isAllVariationCombinationExist(productID uint, variationOptionIDs []uint) (exist bool, err error) {

	for _, variationOptionID := range variationOptionIDs {
		exist, err := c.productRepo.IsProductItemAlreadyExist(context.Background(), productID, variationOptionID)
		if err != nil {
			return false, utils.PrependMessageToError(err, "failed to check product item already exist with given configuration")
		}
		// one of the combination not exist then return
		if !exist {
			return false, nil
		}
	}

	return true, nil
}

// for get all productItem for a specific product
func (c *productUseCase) FindProductItems(ctx context.Context, productID uint) (productItems []response.ProductItems, err error) {

	productItems, err = c.productRepo.FindAllProductItems(ctx, productID)
	if err != nil {
		return productItems, err
	}

	// get all variation values of the product items
	for i, productItem := range productItems {
		variationValues, err := c.productRepo.FindAllVariationValuesOfProductItem(ctx, productItem.ID)
		if err != nil {
			return nil, utils.PrependMessageToError(err, "failed to find variation values product item")
		}
		productItems[i].VariationValues = variationValues
	}

	// for _, productItem := range productItems {
	// 	images, err := c.productRepo.FindAllProductItemImages(ctx, productItem.ID)

	// 	if err != nil {
	// 		return productItems, err
	// 	}
	// 	fmt.Println(images, "images")
	// }

	return productItems, nil
}

func (c *productUseCase) UpdateProduct(ctx context.Context, product domain.Product) error {

	// check the given product_name already exist or not
	existProduct, err := c.productRepo.FindProductByName(ctx, product.Name)
	if err != nil {
		return err
	}
	// if we found a product with this name but not for this id means another product exist with this name
	if existProduct.ID != 0 && existProduct.ID != product.ID {
		return utils.PrependMessageToError(ErrProductAlreadyExist, "product name "+product.Name)
	}

	err = c.productRepo.UpdateProduct(ctx, product)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to update product")
	}
	return nil
}
