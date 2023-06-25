package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

func (c *productUseCase) SaveOffer(ctx context.Context, offer request.Offer) error {

	existOffer, err := c.productRepo.FindOfferByName(ctx, offer.Name)
	if err != nil {
		return utils.PrependMessageToError(err, "failed check offer name already exist")
	}
	if existOffer.ID != 0 {
		return ErrOfferNameAlreadyExist
	}

	// check the offer end date is valid
	if time.Since(offer.EndDate) > 0 {
		return ErrInvalidOfferEndDate
	}

	err = c.productRepo.SaveOffer(ctx, offer)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to save offer")
	}

	return nil
}

func (c *productUseCase) RemoveOffer(ctx context.Context, offerID uint) error {

	err := c.productRepo.Transactions(ctx, func(repo interfaces.ProductRepository) error {
		// first delete all offer categories based on the removing offer
		err := repo.DeleteAllCategoryOffersByOfferID(ctx, offerID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to remove all category offer related to given offer")
		}
		// delete all product offer based on the removing offer
		err = repo.DeleteAllProductOffersByOfferID(ctx, offerID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to remove all product offer related to given offer")
		}

		// remove the offer
		err = repo.DeleteOffer(ctx, offerID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to remove offer")
		}

		// re-calculate products after removed offer by category offer wise
		// err = repo.ReCalculateAllProductsDiscountByCategoryOffer(ctx)
		// if err != nil {
		// 	return utils.PrependMessageToError(err, "failed to re calculate products discount by category offer")
		// }
		// // re-calculate products after removed offer by product offer wise
		// err = repo.ReCalculateAllProductsDiscountByProductOffer(ctx)
		// if err != nil {
		// 	return utils.PrependMessageToError(err, "failed to re calculate products discount by product offer")
		// }

		// // re-calculate product items after removed offer by category offer wise
		// err = repo.ReCalculateAllProductItemsDiscountByCategoryOffer(ctx)
		// if err != nil {
		// 	return utils.PrependMessageToError(err, "failed to re calculate product items discount by category offer")
		// }
		// // re-calculate product items after removed offer by product offer wise
		// err = repo.ReCalculateAllProductItemsDiscountByProductOffer(ctx)
		// if err != nil {
		// 	return utils.PrependMessageToError(err, "failed to re calculate product items discount by product offer")
		// }

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (c *productUseCase) FindAllOffers(ctx context.Context, pagination request.Pagination) ([]domain.Offer, error) {

	offers, err := c.productRepo.FindAllOffers(ctx, pagination)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed to find all offers")
	}
	return offers, nil
}

func (c *productUseCase) SaveCategoryOffer(ctx context.Context, offerCategory request.OfferCategory) error {

	offer, err := c.productRepo.FindOfferByID(ctx, offerCategory.OfferID)
	if err != nil {
		return err
	}

	//check the offer date is end or not
	if time.Since(offer.EndDate) > 0 {
		return ErrOfferAlreadyEnded
	}

	//  check the category have already offer exist or not
	category, err := c.productRepo.FindOfferCategoryCategoryID(ctx, offerCategory.CategoryID)
	if err != nil {
		return err
	}
	if category.ID != 0 {
		return ErrCategoryOfferAlreadyExist
	}

	err = c.productRepo.Transactions(ctx, func(repo interfaces.ProductRepository) error {
		// save category offer
		categoryOfferID, err := repo.SaveCategoryOffer(ctx, offerCategory)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to save category offer")
		}
		// calculate products after removed offer by category offer wise
		err = repo.UpdateProductsDiscountByCategoryOfferID(ctx, categoryOfferID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to re calculate products discount by category offer")
		}
		// calculate product items after removed offer by category offer wise
		err = repo.UpdateProductItemsDiscountByCategoryOfferID(ctx, categoryOfferID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to re calculate product items discount by category offer")
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// get all offer_category
func (c *productUseCase) FindAllCategoryOffers(ctx context.Context, pagination request.Pagination) ([]response.OfferCategory, error) {

	categoryOffers, err := c.productRepo.FindAllOfferCategories(ctx, pagination)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed to find all category offers")
	}

	return categoryOffers, nil
}

// remove offer from category
func (c *productUseCase) RemoveCategoryOffer(ctx context.Context, categoryOfferID uint) error {

	err := c.productRepo.Transactions(ctx, func(repo interfaces.ProductRepository) error {

		err := c.productRepo.DeleteCategoryOffer(ctx, categoryOfferID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to remove category offer")
		}
		// re-calculate products after removed offer by category offer wise
		// err = repo.ReCalculateAllProductsDiscountByCategoryOffer(ctx)
		// if err != nil {
		// 	return utils.PrependMessageToError(err, "failed to re calculate products discount by category offer")
		// }
		// // re-calculate product items after removed offer by category offer wise
		// err = repo.ReCalculateAllProductItemsDiscountByCategoryOffer(ctx)
		// if err != nil {
		// 	return utils.PrependMessageToError(err, "failed to re calculate product items discount by category offer")
		// }
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *productUseCase) ReplaceCategoryOffer(ctx context.Context, offerCategory domain.OfferCategory) error {

	// if offer exist then update it
	err := c.productRepo.UpdateCategoryOffer(ctx, offerCategory.CategoryID, offerCategory.OfferID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to update category offer")
	}

	// err = c.productRepo.ReCalculateAllDiscountPrice(ctx)
	// if err != nil {
	// 	utils.PrependMessageToError(err, "failed to update discount price after replace category offer")
	// }
	return nil
}

// offer on products

func (c *productUseCase) SaveProductOffer(ctx context.Context, offerProduct domain.OfferProduct) error {

	// check the any offer is already exist for the given product
	offerProduct, err := c.productRepo.FindOfferProductByProductID(ctx, offerProduct.ProductID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to check product have already offer exist")
	}
	if offerProduct.ID != 0 {
		return ErrProductOfferAlreadyExist
	}

	err = c.productRepo.Transactions(ctx, func(repo interfaces.ProductRepository) error {
		productOfferID, err := repo.SaveOfferProduct(ctx, offerProduct)
		if err != nil {
			return utils.PrependMessageToError(err, "failed save product offer")
		}

		err = repo.UpdateProductsDiscountByProductOfferID(ctx, productOfferID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to calculate product discount price for offer")
		}

		err = repo.UpdateProductItemsDiscountByProductOfferID(ctx, productOfferID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to calculate product items discount price for offer")
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// get all offers for products
func (c *productUseCase) FindAllProductOffers(ctx context.Context, pagination request.Pagination) ([]response.OfferProduct, error) {
	return c.productRepo.FindAllOfferProducts(ctx, pagination)
}

// remove offer form products
func (c *productUseCase) RemoveProductOffer(ctx context.Context, offerProducts domain.OfferProduct) error {

	offerProducts, err := c.productRepo.FindOfferProduct(ctx, offerProducts)
	if err != nil {
		return err
	} else if offerProducts.OfferID == 0 {
		return errors.New("invalid offer_product_id")
	}

	if err := c.productRepo.DeleteOfferProduct(ctx, offerProducts); err != nil {
		return err
	}

	// update the discount price
	return nil //c.productRepo.ReCalculateAllDiscountPrice(ctx)
}

// replace offer products
func (c *productUseCase) ReplaceProductOffer(ctx context.Context, offerProduct domain.OfferProduct) error {

	err := c.productRepo.UpdateOfferProduct(ctx, offerProduct)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to replace offer product")
	}

	// update the discount price
	return nil //c.productRepo.ReCalculateAllDiscountPrice(ctx)
}
