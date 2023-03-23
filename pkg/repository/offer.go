package repository

import (
	"context"
	"errors"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

// find offer using id or offer_name
func (c *productDatabase) FindOffer(ctx context.Context, offer domain.Offer) (domain.Offer, error) {

	query := `SELECT * FROM offers WHERE id = ? OR offer_name = ?`
	if c.DB.Raw(query, offer.ID, offer.OfferName).Scan(&offer).Error != nil {
		return offer, errors.New("faild to find offer")
	}

	return offer, nil
}

// save a new offer
func (c *productDatabase) SaveOffer(ctx context.Context, offer domain.Offer) error {

	query := `INSERT INTO offers (offer_name,description,discount_rate,start_date,end_date) VALUES ($1,$2,$3,$4,$5)`

	if c.DB.Exec(query, offer.OfferName, offer.Description, offer.DiscountRate, offer.StartDate, offer.EndDate).Error != nil {
		return errors.New("faild to save offer")
	}
	return nil
}

// find offer_category using is or offer_id and category_id
func (c *productDatabase) FindOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) (domain.OfferCategory, error) {

	query := `SELECT * FROM offer_categories WHERE id = ? OR offer_id = ? AND category_id = ?`
	if c.DB.Raw(query, offerCategory.ID, offerCategory.OfferID, offerCategory.CategoryID).Scan(&offerCategory).Error != nil {
		return offerCategory, errors.New("faild to find offer category")
	}
	return offerCategory, nil
}

// save a new offer for category
func (c *productDatabase) SaveOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error {

	query := `INSERT INTO offer_categories (offer_id,category_id) VALUES ($1,$2)`
	if c.DB.Exec(query, offerCategory.OfferID, offerCategory.CategoryID).Error != nil {
		return errors.New("faild to save offer for category")
	}
	return nil
}

// find product_offer with id or offer_id and product_id
func (c *productDatabase) FindOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) (domain.OfferProduct, error) {

	query := `SELECT * FROM offer_products WHERE id = ? OR offer_id = ? AND product_id = ?`
	if c.DB.Raw(query, offerProduct.ID, offerProduct.OfferID, offerProduct.ProductID).Scan(&offerProduct).Error != nil {
		return offerProduct, errors.New("faild to find offer product")
	}
	return offerProduct, nil
}

// save a offer for product
func (c *productDatabase) SaveOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error {
	query := `INSERT INTO offer_products (offer_id, product_id) VALUES ($1,$2)`
	if c.DB.Exec(query, offerProduct.OfferID, offerProduct.ProductID).Error != nil {
		return errors.New("faild to save offer for product")
	}
	return nil
}
