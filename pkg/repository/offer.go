package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

// Find offer by id
func (c *productDatabase) FindOfferByID(ctx context.Context, offerID uint) (offer domain.Offer, err error) {

	query := `SELECT * FROM offers WHERE id = $1`
	err = c.DB.Raw(query, offerID).Scan(&offer).Error

	return offer, err
}

// Find offer by name
func (c *productDatabase) FindOfferByName(ctx context.Context, offerName string) (offer domain.Offer, err error) {

	query := `SELECT * FROM offers WHERE name = $1`
	err = c.DB.Raw(query, offer.ID, offer.Name).Scan(&offer).Error

	return offer, err
}

// findAll offers
func (c *productDatabase) FindAllOffers(ctx context.Context, pagination request.Pagination) ([]domain.Offer, error) {
	var offers []domain.Offer
	if c.DB.Raw("SELECT * FROM offers").Scan(&offers).Error != nil {
		return offers, errors.New("faild to get all offers")
	}
	return offers, nil
}

// save a new offer
func (c *productDatabase) SaveOffer(ctx context.Context, offer domain.Offer) error {

	query := `INSERT INTO offers (offer_name,description,discount_rate,start_date,end_date) VALUES ($1,$2,$3,$4,$5)`

	if c.DB.Exec(query, offer.Name, offer.Description, offer.DiscountRate, offer.StartDate, offer.EndDate).Error != nil {
		return errors.New("faild to save offer")
	}
	return nil
}

// update an existing offer
func (c *productDatabase) UpdateOffer(ctx context.Context, offer domain.Offer) error {
	query := `UPDATE offers SET offer_name=$1,description=$2,discount_rate=$3,start_date=$4,end_date=$5 WHERE id=$6`
	if c.DB.Exec(query, offer.Name, offer.Description, offer.DiscountRate, offer.StartDate, offer.EndDate, offer.ID).Error != nil {
		return errors.New("faild to update offer")
	}
	return nil
}

// delete an offer
func (c *productDatabase) DeleteOffer(ctx context.Context, offerID uint) error {
	trx := c.DB.Begin()
	// first update all discount price to 0 for //?products
	//which are related by offer_products and offer_category
	query := `UPDATE products p SET discount_price = 0  
	FROM offer_categories oc INNER JOIN offers o ON o.id = oc.offer_id 
	WHERE oc.category_id = p.category_id AND o.id = $1`
	if trx.Exec(query, offerID).Error != nil {
		trx.Rollback()
		return errors.New("faild to remove offer price from products")
	}

	// then upate all discount price to 0 for //? product_items
	// which are related by product_offer and category_offer
	query = `UPDATE product_items AS pi SET discount_price = 0 
	FROM offer_categories oc INNER JOIN offers o ON o.id = oc.offer_id 
	INNER JOIN products p ON p.category_id = oc.category_id 
	WHERE p.id = pi.product_id AND o.id = $1`
	if trx.Exec(query, offerID).Error != nil {
		trx.Rollback()
		return errors.New("faild to remove offer price form product items")
	}

	// then remove all offer_products for the offer
	query = `DELETE FROM offer_products WHERE offer_id = $1`
	if trx.Exec(query, offerID).Error != nil {
		trx.Rollback()
		return errors.New("faild to remove offer products")
	}
	// remove all the offer_categories fot the offer
	query = `DELETE FROM offer_categories WHERE offer_id = $1`
	if trx.Exec(query, offerID).Error != nil {
		trx.Rollback()
		return errors.New("faild to remove offer cateogry")
	}

	// at last remove the offer from offer table
	query = `DELETE FROM offers WHERE id = $1`
	if trx.Exec(query, offerID).Error != nil {
		trx.Rollback()
		return errors.New("fiald to remove offer")
	}
	// commit the transaction
	if trx.Commit().Error != nil {
		trx.Rollback()
		return errors.New("faild to complete the offer removel transaction")
	}
	return nil
}

// find offer_category by category_id (for mainly checking this category have an offer existing or not)
func (c *productDatabase) FindOfferCategoryCategoryID(ctx context.Context, categoryID uint) (domain.OfferCategory, error) {
	var offerCategory domain.OfferCategory
	query := `SELECT * FROM offer_categories WHERE  category_id = ?`
	if c.DB.Raw(query, categoryID).Scan(&offerCategory).Error != nil {
		return offerCategory, errors.New("faild to find offer category")
	}
	return offerCategory, nil
}

// find offer_category by id or offer_id with category_id
func (c *productDatabase) FindOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) (domain.OfferCategory, error) {
	query := `SELECT * FROM offer_categories WHERE id = ? OR offer_id = ? AND category_id = ?`
	if c.DB.Raw(query, offerCategory.ID, offerCategory.OfferID, offerCategory.CategoryID).Scan(&offerCategory).Error != nil {
		return offerCategory, errors.New("faild to find offer category")
	}
	return offerCategory, nil
}

// find all offer_category
func (c *productDatabase) FindAllOfferCategories(ctx context.Context, pagination request.Pagination) ([]response.OfferCategory, error) {
	var offerCategories []response.OfferCategory
	query := `SELECT oc.id AS offer_category_id, oc.category_id,c.category_name,oc.offer_id,o.offer_name, o.discount_rate 
	FROM offer_categories oc INNER JOIN categories c ON c.id = oc.category_id 
	INNER JOIN offers o ON oc.offer_id = o.id`
	if c.DB.Raw(query).Scan(&offerCategories).Error != nil {
		return offerCategories, errors.New("faild to get all offer categories")
	}
	return offerCategories, nil
}

// save a new offer for category
func (c *productDatabase) SaveOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error {
	// first create the offer for category
	query := `INSERT INTO offer_categories (offer_id,category_id) VALUES ($1,$2)`
	if c.DB.Exec(query, offerCategory.OfferID, offerCategory.CategoryID).Error != nil {
		return errors.New("faild to save offer_category")
	}

	return nil
}

// remove offer_category
func (c *productDatabase) DeleteOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error {

	// first remove all discount price related offer_category
	query := `UPDATE products SET discount_price = 0 
	WHERE category_id = $1`
	if c.DB.Exec(query, offerCategory.CategoryID).Error != nil {
		return errors.New("faild to remove product discount_price for offer_category")
	}

	// then remove all discount price related to offer_category
	query = `UPDATE product_items pi SET discount_price = 0 
	FROM products p WHERE p.category_id = $1 AND pi.product_id = p.id`
	if c.DB.Exec(query, offerCategory.CategoryID).Error != nil {
		return errors.New("faild to remove product_items discount price for offer_category")
	}
	// then delete offer_category
	query = `DELETE FROM offer_categories WHERE id = $1 OR offer_id = $2 AND category_id = $3 `
	if c.DB.Exec(query, offerCategory.ID, offerCategory.OfferID, offerCategory.CategoryID).Error != nil {
		return errors.New("faild to remove offer from category on offer_cateogies")
	}

	return nil
}

// update offer_category
func (c *productDatabase) UpdateOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error {

	query := `UPDATE offer_categories SET offer_id = $1 WHERE category_id = $2`
	if c.DB.Exec(query, offerCategory.OfferID, offerCategory.CategoryID).Error != nil {
		return errors.New("faild to update offer_category")
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

// find product_offer with product_id
func (c *productDatabase) FindOfferProductByProductID(ctx context.Context, productID uint) (domain.OfferProduct, error) {
	var offerProduct domain.OfferProduct

	query := `SELECT * FROM offer_products WHERE product_id = ?`
	if c.DB.Raw(query, productID).Scan(&offerProduct).Error != nil {
		return offerProduct, errors.New("faild to find offer product")
	}
	return offerProduct, nil
}

// find all offer_products
func (c *productDatabase) FindAllOfferProducts(ctx context.Context, pagination request.Pagination) ([]response.OfferProduct, error) {
	var offerProducts []response.OfferProduct
	query := `SELECT op.id AS offer_product_id, op.product_id,p.product_name,op.offer_id,o.offer_name, o.discount_rate  
	FROM offer_products op INNER JOIN products p ON p.id = op.product_id 
	INNER JOIN offers o ON o.id = op.offer_id`
	if c.DB.Raw(query).Scan(&offerProducts).Error != nil {
		return offerProducts, errors.New("faild to find all offer products")
	}
	return offerProducts, nil
}

// save a offer for product
func (c *productDatabase) SaveOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error {

	query := `INSERT INTO offer_products (offer_id, product_id) VALUES ($1,$2)`
	err := c.DB.Exec(query, offerProduct.OfferID, offerProduct.ProductID).Error
	if err != nil {
		return errors.New("faild to save offer for offer_product")
	}

	return nil
}

// delete offer_products
func (c *productDatabase) DeleteOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error {

	// first remove discount_price of product which related to this offer_product
	query := `UPDATE products SET discount_price = 0 WHERE id = $1`
	if c.DB.Exec(query, offerProduct.ProductID).Error != nil {
		return errors.New("faild to remove discount price from product")
	}
	// then remove discount price from product_items
	query = `UPDATE product_items SET discount_price = 0 WHERE product_id = $1`
	if c.DB.Exec(query, offerProduct.ProductID).Error != nil {
		return errors.New("faild to remove discount_price from product_items")
	}

	// then delete the offer_produt
	query = `DELETE FROM offer_products WHERE id = $1 OR offer_id = $2 AND product_id = $3`
	if c.DB.Exec(query, offerProduct.ID, offerProduct.OfferID, offerProduct.ProductID).Error != nil {
		return errors.New("faild to remove offer from products on offer_products")
	}

	return nil
}

// update offer_products
func (c *productDatabase) UpdateOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error {

	query := `UPDATE offer_products SET offer_id = $1 WHERE product_id = $1`
	if c.DB.Exec(query, offerProduct.OfferID, offerProduct.ProductID).Error != nil {
		return errors.New("faild to update offer_products on offer update")
	}
	return nil
}

// update all discount price first category wise then product_wise
func (c *productDatabase) UpdateDiscountPrice(ctx context.Context) error {
	fmt.Println("\\updating product discount price//")
	trx := c.DB.Begin()

	// update the all products discount price
	query := `UPDATE products SET discount_price = (price * (100 - o.discount_rate))/100 
	FROM offer_categories oc INNER JOIN offers o ON o.id = oc.offer_id 
	WHERE oc.category_id = products.category_id`

	if trx.Exec(query).Error != nil {
		trx.Rollback()
		return errors.New("faild to update discount price on product for category offer")
	}

	// update all product_items discount price
	query = `UPDATE product_items AS pi SET discount_price = (pi.price * (100 - o.discount_rate))/100 
	FROM offer_categories oc INNER JOIN offers o ON o.id = oc.offer_id 
	INNER JOIN products p ON p.category_id = oc.category_id 
	WHERE p.id = pi.product_id`
	if trx.Exec(query).Error != nil {
		return errors.New("faild to update discount price on product_items for category offer")
	}

	query = `UPDATE products p SET discount_price = (p.price * (100 - o.discount_rate))/100 
	FROM offers o INNER JOIN offer_products op ON o.id = op.offer_id 
	WHERE p.id = op.product_id`
	if trx.Exec(query).Error != nil {
		trx.Rollback()
		return errors.New("faild to update product discount price on offer_products")
	}

	// then update product_items discont price
	query = `UPDATE product_items pi SET discount_price = (pi.price * (100 - o.discount_rate))/100 
	FROM offers o INNER JOIN offer_products op ON o.id = op.offer_id  
	WHERE pi.product_id = op.product_id`
	if trx.Exec(query).Error != nil {
		trx.Rollback()
		return errors.New("faild to update product_items discount price on offer_products")
	}

	if trx.Commit().Error != nil {
		return errors.New("faild to complete the add offer category")
	}
	return nil
}
