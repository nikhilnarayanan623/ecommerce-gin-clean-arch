package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

func (c *productUseCase) AddOffer(ctx context.Context, offer domain.Offer) error {

	if offer, err := c.productRepo.FindOffer(ctx, offer); err != nil {
		return err
	} else if offer.ID != 0 {
		return fmt.Errorf("offer already exist with this  %s  name", offer.Name)
	}

	// validate the offer date
	if time.Since(offer.EndDate) > 0 {
		return fmt.Errorf("given offer end_date already exceeded %v", offer.EndDate)
	}

	err := c.productRepo.SaveOffer(ctx, offer)
	if err != nil {
		return err
	}

	log.Printf("successfully offer created with offer name %v\n\n", offer.Name)
	return nil
}

func (c *productUseCase) RemoveOffer(ctx context.Context, offerID uint) error {

	offer := domain.Offer{ID: offerID}
	if offer, err := c.productRepo.FindOffer(ctx, offer); err != nil {
		return err
	} else if offer.Name == "" {
		return errors.New("invalid offer id")
	}

	if err := c.productRepo.DeleteOffer(ctx, offerID); err != nil {
		return err
	}
	// upate discount price ( after removal discount price there is chance other offer category will relate with this offers products)
	return c.productRepo.UpdateDiscountPrice(ctx)
}

func (c *productUseCase) FindAllOffers(ctx context.Context) ([]domain.Offer, error) {

	return c.productRepo.FindAllOffer(ctx)
}

func (c *productUseCase) AddOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error {

	// check the offerId is valid or not
	var offer = domain.Offer{ID: offerCategory.OfferID}

	offer, err := c.productRepo.FindOffer(ctx, offer)
	if err != nil {
		return err
	} else if offer.Name == "" {
		return errors.New("invalid offer_id")
	}

	//check the offer date is end or not
	if time.Since(offer.EndDate) > 0 {
		return fmt.Errorf("can't apply to category \noffer already ended on %v ", offer.EndDate)
	}

	category, err := c.productRepo.FindCategoryByID(ctx, offerCategory.CategoryID)
	if err != nil {
		return err
	} else if category.Name == "" {
		return errors.New("invalid category_id")
	}

	//  check the category have already offer exist or not
	checkofferCategory, err := c.productRepo.FindOfferCategoryCategoryID(ctx, offerCategory.CategoryID)
	if err != nil {
		return err
	} else if checkofferCategory.ID != 0 {
		return errors.New("an offer already exist for this category you can replace it")
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
func (c *productUseCase) FindAllOffersOfCategories(ctx context.Context) ([]response.OfferCategory, error) {

	return c.productRepo.FindAllOfferCategories(ctx)
}

// remove offer from category
func (c *productUseCase) RemoveOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error {

	offerCategory, err := c.productRepo.FindOfferCategory(ctx, offerCategory)
	if err != nil {
		return err
	} else if offerCategory.OfferID == 0 {
		return errors.New("invalid offer_category_id")
	}

	if err := c.productRepo.DeleteOfferCategory(ctx, offerCategory); err != nil {
		return err
	}
	// update the discount price
	return c.productRepo.UpdateDiscountPrice(ctx)
}

func (c *productUseCase) ReplaceOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error {
	// check the offerId is valid or not
	var offer = domain.Offer{ID: offerCategory.OfferID}
	if offer, err := c.productRepo.FindOffer(ctx, offer); err != nil {
		return err
	} else if offer.Name == "" {
		return errors.New("invalid offer_id")
	}

	category, err := c.productRepo.FindCategoryByID(ctx, offerCategory.CategoryID)
	if err != nil {
		return err
	} else if category.Name == "" {
		return errors.New("invalid category_id")
	}

	//  check the given category offer for replacing category exist or not
	if offerCategory, err := c.productRepo.FindOfferCategoryCategoryID(ctx, offerCategory.CategoryID); err != nil {
		return err
	} else if offerCategory.OfferID == 0 {
		return errors.New("there is no offer not exist this category for replacing")
	}
	// if offer exist then update it
	if err := c.productRepo.UpdateOfferCategory(ctx, offerCategory); err != nil {
		return err
	}

	// update the discount price
	return c.productRepo.UpdateDiscountPrice(ctx)
}

// offer on products

func (c *productUseCase) AddOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error {
	// check the offerId is valid or not
	var offer = domain.Offer{ID: offerProduct.OfferID}
	offer, err := c.productRepo.FindOffer(ctx, offer)
	if err != nil {
		return err
	} else if offer.Name == "" {
		return errors.New("invalid offer_id")
	}

	product, err := c.productRepo.FindProductByID(ctx, offerProduct.ProductID)
	if err != nil {
		return err
	} else if product.ID == 0 {
		return errors.New("invalid product_id")
	}

	// check the any offer is already exist for the given product
	if offerProduct, err := c.productRepo.FindOfferProductByProductID(ctx, offerProduct.ProductID); err != nil {
		return err
	} else if offerProduct.ID != 0 {
		return errors.New("this offer alredy exist for given product")
	}

	// if not exist then add it
	if err := c.productRepo.SaveOfferProduct(ctx, offerProduct); err != nil {
		return err
	}

	// update the discount price
	return c.productRepo.UpdateDiscountPrice(ctx)
}

// get all offers for products
func (c *productUseCase) FindAllOffersOfProducts(ctx context.Context) ([]response.OfferProduct, error) {
	return c.productRepo.FindAllOfferProducts(ctx)
}

// remove offer form products
func (c *productUseCase) RemoveOfferProducts(ctx context.Context, offerProdcts domain.OfferProduct) error {

	offerProducts, err := c.productRepo.FindOfferProduct(ctx, offerProdcts)
	if err != nil {
		return err
	} else if offerProducts.OfferID == 0 {
		return errors.New("invalid offer_product_id")
	}

	if err := c.productRepo.DeleteOfferProduct(ctx, offerProdcts); err != nil {
		return err
	}

	// update the discount price
	return c.productRepo.UpdateDiscountPrice(ctx)
}

// replace offer products
func (c *productUseCase) ReplaceOfferProducts(ctx context.Context, offerProduct domain.OfferProduct) error {
	// check the offerId is valid or not
	var offer = domain.Offer{ID: offerProduct.OfferID}
	offer, err := c.productRepo.FindOffer(ctx, offer)
	if err != nil {
		return err
	} else if offer.Name == "" {
		return errors.New("invalid offer_id")
	}

	// check the product id is valid or not
	product, err := c.productRepo.FindProductByID(ctx, offerProduct.ProductID)
	if err != nil {
		return err
	} else if product.ID == 0 {
		return errors.New("invalid product_id")
	}

	// check the offer is already exist for the given product
	if offerProduct, err := c.productRepo.FindOfferProductByProductID(ctx, offerProduct.ProductID); err != nil {
		return err
	} else if offerProduct.ID == 0 {
		return errors.New("invalid api call\nthere is no offer exist for this products so can't replace")
	}

	if err := c.productRepo.UpdateOfferProduct(ctx, offerProduct); err != nil {
		return err
	}

	// update the discount price
	return c.productRepo.UpdateDiscountPrice(ctx)
}
