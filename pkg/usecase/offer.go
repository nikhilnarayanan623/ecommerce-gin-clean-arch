package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

func (c *productUseCase) AddOffer(ctx context.Context, offer domain.Offer) error {

	if offer, err := c.productRepo.FindOffer(ctx, offer); err != nil {
		return err
	} else if offer.ID != 0 {
		return fmt.Errorf("offer already exist with this  %s  name", offer.OfferName)
	}

	return c.productRepo.SaveOffer(ctx, offer)
}

func (c *productUseCase) RemoveOffer(ctx context.Context, offerID uint) error {

	offer := domain.Offer{ID: offerID}
	if offer, err := c.productRepo.FindOffer(ctx, offer); err != nil {
		return err
	} else if offer.OfferName == "" {
		return errors.New("invalid offer id")
	}

	if err := c.productRepo.DeleteOffer(ctx, offerID); err != nil {
		return err
	}
	// upate discount price ( after removal discount price there is chance other offer category will relate with this offers products)
	return c.productRepo.UpdateDiscountPrice(ctx)
}

func (c *productUseCase) GetAllOffers(ctx context.Context) ([]domain.Offer, error) {

	return c.productRepo.FindAllOffer(ctx)
}

func (c *productUseCase) AddOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error {

	// check the offerId is valid or not
	var offer = domain.Offer{ID: offerCategory.OfferID}
	if offer, err := c.productRepo.FindOffer(ctx, offer); err != nil {
		return err
	} else if offer.OfferName == "" {
		return errors.New("invalid offer_id")
	}

	// check the categoy id is valid or not
	var category = domain.Category{ID: offerCategory.CategoryID}
	if category, err := c.productRepo.FindCategory(ctx, category); err != nil {
		return err
	} else if category.CategoryName == "" {
		return errors.New("invalid category_id")
	}

	//  check the category have already offer exist or not
	if offerCategory, err := c.productRepo.FindOfferCategoryCategoryID(ctx, offerCategory.CategoryID); err != nil {
		return err
	} else if offerCategory.ID != 0 {
		return errors.New("an offer already exist for this category you can replace it")
	}

	// if it not exist then add it
	if err := c.productRepo.SaveOfferCategory(ctx, offerCategory); err != nil {
		return err
	}

	// last update discount price
	return c.productRepo.UpdateDiscountPrice(ctx)
}

// get all offer_category
func (c *productUseCase) GetAllOffersOfCategories(ctx context.Context) ([]res.ResOfferCategory, error) {

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
	} else if offer.OfferName == "" {
		return errors.New("invalid offer_id")
	}

	// check the categoy id is valid or not
	var category = domain.Category{ID: offerCategory.CategoryID}
	if category, err := c.productRepo.FindCategory(ctx, category); err != nil {
		return err
	} else if category.CategoryName == "" {
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
	} else if offer.OfferName == "" {
		return errors.New("invalid offer_id")
	}

	// check the product id is valid or not
	var product = domain.Product{ID: offerProduct.ProductID}
	product, err = c.productRepo.FindProduct(ctx, product)
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
func (c *productUseCase) GetAllOffersOfProducts(ctx context.Context) ([]res.ResOfferProduct, error) {
	return c.productRepo.FindAllOfferProducts(ctx)
}

// remove offer form products
func (c *productUseCase) RemoveOfferProducts(ctx context.Context, offerProdcts domain.OfferProduct) error {

	offerProdcts, err := c.productRepo.FindOfferProduct(ctx, offerProdcts)
	if err != nil {
		return err
	} else if offerProdcts.OfferID == 0 {
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
	} else if offer.OfferName == "" {
		return errors.New("invalid offer_id")
	}

	// check the product id is valid or not
	var product = domain.Product{ID: offerProduct.ProductID}
	product, err = c.productRepo.FindProduct(ctx, product)
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
