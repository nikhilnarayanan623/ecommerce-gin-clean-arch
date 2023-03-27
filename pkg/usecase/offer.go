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

	return c.productRepo.DeleteOffer(ctx, offerID)
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

func (c *productUseCase) OfferCategoryPage(ctx context.Context) (res.ResOfferCategoryPage, error) {

	var (
		resOfferCategoryData res.ResOfferCategoryPage
		err                  error
	)
	// get all offers
	if resOfferCategoryData.Offers, err = c.productRepo.FindAllOffer(ctx); err != nil {
		return resOfferCategoryData, err
	}

	// get all categories
	if resOfferCategoryData.Categories, err = c.productRepo.FindAllCategories(ctx); err != nil {
		return resOfferCategoryData, err
	}
	return resOfferCategoryData, nil
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

// remove offer from category
func (c *productUseCase) RemoveOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error {

	offerCategory, err := c.productRepo.FindOfferCategory(ctx, offerCategory)
	if err != nil {
		return err
	} else if offerCategory.OfferID == 0 {
		return errors.New("invalid offer_category_id")
	}

	return c.productRepo.DeleteOfferCategory(ctx, offerCategory)
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
	return c.productRepo.UpdateOfferCategory(ctx, offerCategory)
}

// offer on products

func (c *productUseCase) OfferProductsPage(ctx context.Context) (res.ResOfferProductsPage, error) {
	var (
		resOfferProductsData res.ResOfferProductsPage
		err                  error
	)
	// find all offers
	if resOfferProductsData.Offers, err = c.productRepo.FindAllOffer(ctx); err != nil {
		return resOfferProductsData, err
	}
	// find all products
	if resOfferProductsData.Products, err = c.productRepo.GetProducts(ctx); err != nil {
		return resOfferProductsData, err
	}

	return resOfferProductsData, nil
}

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
	return c.productRepo.SaveOfferProduct(ctx, offerProduct)
}

// remove offer form products
func (c *productUseCase) RemoveOfferProducts(ctx context.Context, offerProdcts domain.OfferProduct) error {

	offerProdcts, err := c.productRepo.FindOfferProduct(ctx, offerProdcts)
	if err != nil {
		return err
	} else if offerProdcts.OfferID == 0 {
		return errors.New("invalid offer_product_id")
	}

	return c.productRepo.DeleteOfferProduct(ctx, offerProdcts)
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

	return c.productRepo.UpdateOfferProduct(ctx, offerProduct)
}
