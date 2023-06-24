package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

func (c *productUseCase) SaveOffer(ctx context.Context, offer domain.Offer) error {

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

	err := c.productRepo.DeleteOffer(ctx, offerID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to remove offer")
	}

	err = c.productRepo.UpdateDiscountPrice(ctx)
	if err != nil {
		utils.PrependMessageToError(err, "failed to update discount price after removing offer")
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

func (c *productUseCase) SaveCategoryOffer(ctx context.Context, offerCategory domain.OfferCategory) error {

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
		return ErrOfferAlreadyExistForCategory
	}

	// if it not exist then add it
	err = c.productRepo.SaveOfferCategory(ctx, offerCategory)
	if err != nil {
		return err
	}

	// last update discount price
	err = c.productRepo.UpdateDiscountPrice(ctx)

	if err != nil {
		return err
	}

	log.Printf("successfully offer applied to category of category_id %v with offer_if %v", offerCategory.CategoryID, offerCategory.OfferID)
	return nil
}

// get all offer_category
func (c *productUseCase) FindAllCategoryOffers(ctx context.Context, pagination request.Pagination) ([]response.OfferCategory, error) {

	return c.productRepo.FindAllOfferCategories(ctx, pagination)
}

// remove offer from category
func (c *productUseCase) RemoveCategoryOffer(ctx context.Context, offerCategory domain.OfferCategory) error {

	offerCategory, err := c.productRepo.FindOfferCategory(ctx, offerCategory)
	if err != nil {
		return err
	} else if offerCategory.OfferID == 0 {
		return errors.New("invalid offer_category_id")
	}

	if err := c.productRepo.DeleteOfferCategory(ctx, offerCategory); err != nil {
		return err
	}

	err = c.productRepo.UpdateDiscountPrice(ctx)
	if err != nil {
		utils.PrependMessageToError(err, "failed to update discount price after removing category offer")
	}
	return nil
}

func (c *productUseCase) ReplaceCategoryOffer(ctx context.Context, offerCategory domain.OfferCategory) error {

	// if offer exist then update it
	err := c.productRepo.UpdateOfferCategory(ctx, offerCategory)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to update category offer")
	}

	err = c.productRepo.UpdateDiscountPrice(ctx)
	if err != nil {
		utils.PrependMessageToError(err, "failed to update discount price after replace category offer")
	}
	return nil
}

// offer on products

func (c *productUseCase) SaveProductOffer(ctx context.Context, offerProduct domain.OfferProduct) error {

	// check the any offer is already exist for the given product
	if offerProduct, err := c.productRepo.FindOfferProductByProductID(ctx, offerProduct.ProductID); err != nil {
		return err
	} else if offerProduct.ID != 0 {
		return errors.New("this offer already exist for given product")
	}

	// if not exist then add it
	if err := c.productRepo.SaveOfferProduct(ctx, offerProduct); err != nil {
		return err
	}

	// update the discount price
	return c.productRepo.UpdateDiscountPrice(ctx)
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
	return c.productRepo.UpdateDiscountPrice(ctx)
}

// replace offer products
func (c *productUseCase) ReplaceProductOffer(ctx context.Context, offerProduct domain.OfferProduct) error {

	err := c.productRepo.UpdateOfferProduct(ctx, offerProduct)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to replace offer product")
	}

	// update the discount price
	return c.productRepo.UpdateDiscountPrice(ctx)
}
