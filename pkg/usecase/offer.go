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
	fmt.Println(offer)
	return c.productRepo.SaveOffer(ctx, offer)
}

func (c *productUseCase) GetAllOffers(ctx context.Context) (res.ResOffer, error) {

	var (
		resOffer res.ResOffer
		err      error
	)
	// find all offers
	if resOffer.Offers, err = c.productRepo.FindAllOffer(ctx); err != nil {
		return resOffer, err
	}
	// find all offer categories
	if resOffer.OfferCategories, err = c.productRepo.FindAllOfferCategories(ctx); err != nil {
		return resOffer, err
	}
	// find all offer products
	if resOffer.OfferProducts, err = c.productRepo.FindAllOfferProducts(ctx); err != nil {
		return resOffer, err
	}

	return resOffer, nil
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
	return c.productRepo.SaveOfferCategory(ctx, offerCategory)
}

func (c *productUseCase) ReplaceOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error {
	// check the offerId is valid or not
	var offer = domain.Offer{ID: offerCategory.OfferID}
	if offer, err := c.productRepo.FindOffer(ctx, offer); err != nil {
		return err
	} else if offer.OfferName != "" {
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
	} else if offerCategory.ID == 0 {
		return errors.New("there is no offer not exist this category for replacing")
	}
	// if offer exist then update it
	return c.productRepo.UpdateOfferCategory(ctx, offerCategory)
}

func (c *productUseCase) AddOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error {
	// check the offerId is valid or not
	var offer = domain.Offer{ID: offerProduct.OfferID}
	if offer, err := c.productRepo.FindOffer(ctx, offer); err != nil {
		return err
	} else if offer.OfferName == "" {
		return errors.New("invalid offer_id")
	}
	// check the product id is valid or not
	var product = domain.Product{ID: offerProduct.ProductID}
	if product, err := c.productRepo.FindProduct(ctx, product); err != nil {
		return err
	} else if product.ID == 0 {
		return errors.New("invalid product_id")
	}
	// check the offer is already exist for the given product
	if offerProduct, err := c.productRepo.FindOfferProduct(ctx, offerProduct); err != nil {
		return err
	} else if offerProduct.ID != 0 {
		return errors.New("this offer alredy exist for given product")
	}
	// if not exist then add it
	return c.productRepo.SaveOfferProduct(ctx, offerProduct)
}
