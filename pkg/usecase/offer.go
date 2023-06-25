package usecase

import (
	"context"
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

		// re-calculate products after removed offer by category offer wise
		err := repo.RemoveProductsDiscountByCategoryOfferID(ctx, categoryOfferID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to remove products discount by category offer")
		}
		// re-calculate product items after removed offer by category offer wise
		err = repo.RemoveProductItemsDiscountByCategoryOfferID(ctx, categoryOfferID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to remove product items discount by category offer")
		}
		// last remove the offer
		err = c.productRepo.DeleteCategoryOffer(ctx, categoryOfferID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to remove category offer")
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *productUseCase) ChangeCategoryOffer(ctx context.Context, categoryOfferID, offerID uint) error {

	err := c.productRepo.Transactions(ctx, func(repo interfaces.ProductRepository) error {
		err := c.productRepo.UpdateCategoryOffer(ctx, categoryOfferID, offerID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to update category offer")
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
		// save product offer
		productOfferID, err := repo.SaveOfferProduct(ctx, offerProduct)
		if err != nil {
			return utils.PrependMessageToError(err, "failed save product offer")
		}
		// update the discount price of products
		err = repo.UpdateProductsDiscountByProductOfferID(ctx, productOfferID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to calculate product discount price for offer")
		}
		// update the discount price of products
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
	productOffers, err := c.productRepo.FindAllOfferProducts(ctx, pagination)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed to find product offers")
	}
	return productOffers, nil
}

// remove offer form products
func (c *productUseCase) RemoveProductOffer(ctx context.Context, productOfferID uint) error {

	err := c.productRepo.Transactions(ctx, func(repo interfaces.ProductRepository) error {

		err := repo.RemoveProductsDiscountByProductOfferID(ctx, productOfferID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to remove discount price of offer product")
		}
		err = repo.RemoveProductItemsDiscountByProductOfferID(ctx, productOfferID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to remove discount price of offer product items")
		}

		err = repo.DeleteOfferProduct(ctx, productOfferID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to remove product offer")
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *productUseCase) ChangeProductOffer(ctx context.Context, productOfferID, offerID uint) error {

	err := c.productRepo.Transactions(ctx, func(repo interfaces.ProductRepository) error {
		err := c.productRepo.UpdateOfferProduct(ctx, productOfferID, offerID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to update product offer")
		}
		// calculate products after removed offer by category offer wise
		err = repo.UpdateProductsDiscountByProductOfferID(ctx, productOfferID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to re calculate products discount by product offer")
		}
		// calculate product items after removed offer by category offer wise
		err = repo.UpdateProductItemsDiscountByProductOfferID(ctx, productOfferID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to re calculate product items discount by product offer")
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
