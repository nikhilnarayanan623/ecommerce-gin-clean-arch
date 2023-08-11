package usecase

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/service/cloud"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
)

type productUseCase struct {
	productRepo  interfaces.ProductRepository
	cloudService cloud.CloudService
}

// to get a new instance of productUseCase
func NewProductUseCase(productRepo interfaces.ProductRepository, cloudService cloud.CloudService) service.ProductUseCase {
	return &productUseCase{
		productRepo:  productRepo,
		cloudService: cloudService,
	}
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

	categoryExist, err := c.productRepo.IsCategoryNameExist(ctx, categoryName)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to check category already exist")
	}
	if categoryExist {
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

	subCatExist, err := c.productRepo.IsSubCategoryNameExist(ctx, subCategory.Name, subCategory.CategoryID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to check sub category already exist")
	}
	if subCatExist {
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

			variationExist, err := repo.IsVariationNameExistForCategory(ctx, variationName, categoryID)
			if err != nil {
				return utils.PrependMessageToError(err, "failed to check variation already exist")
			}

			if variationExist {
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

			valueExist, err := repo.IsVariationValueExistForVariation(ctx, variationValue, variationID)
			if err != nil {
				return utils.PrependMessageToError(err, "failed to check variation already exist")
			}
			if valueExist {
				return utils.PrependMessageToError(ErrVariationOptionAlreadyExist, "variation option value "+variationValue)
			}

			err = repo.SaveVariationOption(ctx, variationID, variationValue)
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
func (c *productUseCase) FindAllProducts(ctx context.Context, pagination request.Pagination) ([]response.Product, error) {
	products, err := c.productRepo.FindAllProducts(ctx, pagination)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed to get product details from database")
	}

	for i := range products {

		url, err := c.cloudService.GetFileUrl(ctx, products[i].Image)
		if err != nil {
			continue
		}
		products[i].Image = url
	}

	return products, nil
}

// to add new product
func (c *productUseCase) SaveProduct(ctx context.Context, product request.Product) error {

	productNameExist, err := c.productRepo.IsProductNameExist(ctx, product.Name)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to check product name already exist")
	}
	if productNameExist {
		return utils.PrependMessageToError(ErrProductAlreadyExist, "product name "+product.Name)
	}

	uploadID, err := c.cloudService.SaveFile(ctx, product.ImageFileHeader)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to save image on cloud storage")
	}

	err = c.productRepo.SaveProduct(ctx, domain.Product{
		Name:        product.Name,
		Description: product.Description,
		CategoryID:  product.CategoryID,
		Price:       product.Price,
		Image:       uploadID,
	})
	if err != nil {
		return utils.PrependMessageToError(err, "failed to save product")
	}
	return nil
}

// for add new productItem for a specific product
func (c *productUseCase) SaveProductItem(ctx context.Context, productID uint, productItem request.ProductItem) error {

	variationCount, err := c.productRepo.FindVariationCountForProduct(ctx, productID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to get variation count of product from database")
	}

	if len(productItem.VariationOptionIDs) != int(variationCount) {
		return ErrNotEnoughVariations
	}

	// check the given all combination already exist (Color:Red with Size:M)
	productItemExist, err := c.isProductVariationCombinationExist(productID, productItem.VariationOptionIDs)
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

		errChan := make(chan error, 2)
		newCtx, cancel := context.WithCancel(ctx) // for any of one of goroutine get error then cancel the working of other also
		defer cancel()

		go func() {
			// save all product configurations based on given variation option id
			for _, variationOptionID := range productItem.VariationOptionIDs {

				select {
				case <-newCtx.Done():
					return
				default:
					err = trxRepo.SaveProductConfiguration(ctx, productItemID, variationOptionID)
					if err != nil {
						errChan <- utils.PrependMessageToError(err, "failed to save product_item configuration")
						return
					}
				}
			}
			errChan <- nil
		}()

		go func() {
			// save all images for the given product item
			for _, imageFile := range productItem.ImageFileHeaders {

				select {
				case <-newCtx.Done():
					return
				default:
					// upload image on cloud
					uploadID, err := c.cloudService.SaveFile(ctx, imageFile)
					if err != nil {
						errChan <- utils.PrependMessageToError(err, "failed to upload image to cloud")
						return
					}
					// save upload id on database
					err = trxRepo.SaveProductItemImage(ctx, productItemID, uploadID)
					if err != nil {
						errChan <- utils.PrependMessageToError(err, "failed to save image for product item on database")
						return
					}
				}
			}
			errChan <- nil
		}()

		// wait for the both go routine to complete
		for i := 1; i <= 2; i++ {

			select {
			case <-ctx.Done():
				return nil
			case err := <-errChan:
				if err != nil { // if any of the goroutine send error then return the error
					return err
				}
				// no error then continue for the next check of select
			}
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// step 1 : get product_id and and all variation id as function parameter
// step 2 : initialize an map for storing product item id and its count(map[uint]int)
// step 3 : loop through the variation option ids
// step 4 : then find all product items ids with given product id and the loop variation option id
// step 5 : if the product item array length is zero means the configuration not exist return false
// step 6 : then loop through the product items ids array(got from database)
// step 7 : add each id on the map and increment its count
// step 8 : check if any of the product items id's count is greater than the variation options ids length then return true
// step 9 : if the loop exist means product configuration is not exist
func (c *productUseCase) isProductVariationCombinationExist(productID uint, variationOptionIDs []uint) (exist bool, err error) {

	setOfIds := map[uint]int{}

	for _, variationOptionID := range variationOptionIDs {

		productItemIds, err := c.productRepo.FindAllProductItemIDsByProductIDAndVariationOptionID(context.TODO(),
			productID, variationOptionID)
		if err != nil {
			return false, utils.PrependMessageToError(err, "failed to find product item ids from database using product id and variation option id")
		}

		if len(productItemIds) == 0 {
			return false, nil
		}

		for _, productItemID := range productItemIds {

			setOfIds[productItemID]++
			// if any of the ids count is equal to array length it means product item id of this is the existing product item of this configuration
			if setOfIds[productItemID] >= len(variationOptionIDs) {
				return true, nil
			}
		}
	}
	return false, nil
}

// for get all productItem for a specific product
func (c *productUseCase) FindAllProductItems(ctx context.Context, productID uint) ([]response.ProductItems, error) {

	productItems, err := c.productRepo.FindAllProductItems(ctx, productID)
	if err != nil {
		return productItems, err
	}

	errChan := make(chan error, 2)
	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {

		// get all variation values of each product items
		for i := range productItems {

			select { // checking each time context is cancelled or not
			case <-ctx.Done():
				return
			default:
				variationValues, err := c.productRepo.FindAllVariationValuesOfProductItem(ctx, productItems[i].ID)
				if err != nil {
					errChan <- utils.PrependMessageToError(err, "failed to find variation values product item")
					return
				}
				productItems[i].VariationValues = variationValues
			}
		}
		errChan <- nil
	}()

	go func() {
		// get all images of each product items
		for i := range productItems {

			select { // checking each time context is cancelled or not
			case <-newCtx.Done():
				return
			default:
				images, err := c.productRepo.FindAllProductItemImages(ctx, productItems[i].ID)

				imageUrls := make([]string, len(images))

				for j := range images {

					url, err := c.cloudService.GetFileUrl(ctx, images[j])
					if err != nil {
						errChan <- utils.PrependMessageToError(err, "failed to get image url from could service")
					}
					imageUrls[j] = url
				}

				if err != nil {
					errChan <- utils.PrependMessageToError(err, "failed to find images of product item")
					return
				}
				productItems[i].Images = imageUrls
			}
		}
		errChan <- nil
	}()

	// wait for the two routine to complete
	for i := 1; i <= 2; i++ {

		select {
		case <-ctx.Done():
			return nil, nil
		case err := <-errChan:
			if err != nil {
				return nil, err
			}
			// no error then continue for the next check
		}
	}

	return productItems, nil
}

func (c *productUseCase) UpdateProduct(ctx context.Context, updateDetails domain.Product) error {

	nameExistForOther, err := c.productRepo.IsProductNameExistForOtherProduct(ctx, updateDetails.Name, updateDetails.ID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to check product name already exist for other product")
	}

	if nameExistForOther {
		return utils.PrependMessageToError(ErrProductAlreadyExist, "product name "+updateDetails.Name)
	}

	// c.productRepo.FindProductByID(ctx, updateDetails.ID)

	err = c.productRepo.UpdateProduct(ctx, updateDetails)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to update product")
	}
	return nil
}
