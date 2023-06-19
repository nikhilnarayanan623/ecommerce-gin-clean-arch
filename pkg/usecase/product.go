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

	categories, err := c.productRepo.FindAllCategories(ctx, pagination)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed find all categories")
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

	mainCat, err := c.productRepo.FindCategoryByID(ctx, subCategory.CategoryID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to verify category id")
	}

	if mainCat.ID == 0 {
		return ErrInvalidCategoryID
	}

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
func (c *productUseCase) SaveVariation(ctx context.Context, variation request.Variation) error {

	category, err := c.productRepo.FindCategoryByID(ctx, variation.CategoryID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to verify category id")
	}
	if category.ID == 0 {
		return ErrInvalidCategoryID
	}

	existVariation, err := c.productRepo.FindVariationByNameAndCategoryID(ctx, variation.Name, variation.CategoryID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to check variation already exist")
	}

	if existVariation.ID != 0 {
		return ErrVariationAlreadyExist
	}

	err = c.productRepo.SaveVariation(ctx, variation)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to save variation")
	}
	return nil
}

// to add new variation value for variation
func (c *productUseCase) SaveVariationOption(ctx context.Context, variationOption request.VariationOption) error {

	variation, err := c.productRepo.FindVariationByID(ctx, variationOption.VariationID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to verify variation id")
	}
	if variation.ID == 0 {
		return ErrInvalidVariationID
	}

	existVarOpt, err := c.productRepo.FindVariationOptionByValueAndVariationID(ctx, variationOption.Value, variationOption.VariationID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to check variation already exist")
	}
	if existVarOpt.ID != 0 {
		return ErrVariationOptionAlreadyExist
	}

	err = c.productRepo.SaveVariationOption(ctx, variationOption)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to save variation option")
	}
	return nil
}

func (c *productUseCase) FindAllVariationsAndItsValues(ctx context.Context, categoryID uint) ([]response.Variation, error) {

	category, err := c.productRepo.FindCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed to verify category id")
	}
	if category.ID == 0 {
		return nil, ErrInvalidCategoryID
	}

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

	category, err := c.productRepo.FindCategoryByID(ctx, product.CategoryID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to verify category id")
	}
	if category.ID == 0 {
		return ErrInvalidCategoryID
	}

	exist, err := c.productRepo.IsProductNameAlreadyExist(ctx, product.Name)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to check product name already exist")
	}
	if exist {
		return ErrProductAlreadyExist
	}

	err = c.productRepo.SaveProduct(ctx, product)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to save product")
	}
	return nil
}

// for add new productItem for a specific product
func (c *productUseCase) SaveProductItem(ctx context.Context, productItem request.ProductItem) error {

	product, err := c.productRepo.FindProductByID(ctx, productItem.ProductID)
	if err != nil {
		return err
	}
	if product.ID == 0 {
		return ErrInvalidProductID
	}

	// check the given all combination already exist (Color:Red with Size:M)
	productItemExist, err := c.isAllVariationCombinationExist(productItem.ProductID, productItem.VariationOptionID)
	if err != nil {
		return err
	}
	if productItemExist {
		return ErrProductItemAlreadyExist
	}

	err = c.productRepo.Transactions(ctx, func(trxRepo interfaces.ProductRepository) error {

		sku := utils.GenerateSKU()
		newProductItem := domain.ProductItem{
			ProductID:  productItem.ProductID,
			QtyInStock: productItem.QtyInStock,
			Price:      productItem.Price,
			SKU:        sku,
		}

		productItemID, err := trxRepo.SaveProductItem(ctx, newProductItem)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to save product item")
		}

		// save all product configurations based on given variation option id
		for _, variationOptionID := range productItem.VariationOptionID {

			valid, err := c.productRepo.IsValidVariationOptionID(ctx, variationOptionID)
			if err != nil {
				return ErrInvalidVariationID
			}
			if !valid {
				return ErrInvalidVariationOptionID
			}

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

	//validate the productID
	product, err := c.productRepo.FindProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if product.Name == "" {
		return nil, ErrInvalidProductID
	}

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
	// validate the product_id
	checkProduct, err := c.productRepo.FindProductByID(ctx, product.ID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to verify product id")
	}
	if checkProduct.ID == 0 {
		return ErrInvalidProductItemID
	}

	// check the given product_name already exist or not
	checkProduct, err = c.productRepo.FindProductByName(ctx, product.Name)
	if err != nil {
		return err
	}
	if checkProduct.ID != 0 && checkProduct.ID != product.ID {
		return ErrProductAlreadyExist
	}

	err = c.productRepo.UpdateProduct(ctx, product)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to update product")
	}
	return nil
}
