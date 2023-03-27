package repository

import (
	"context"
	"errors"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

// find offer using id or offer_name
func (c *productDatabase) FindOffer(ctx context.Context, offer domain.Offer) (domain.Offer, error) {

	query := `SELECT * FROM offers WHERE id = ? OR offer_name = ?`
	if c.DB.Raw(query, offer.ID, offer.OfferName).Scan(&offer).Error != nil {
		return offer, errors.New("faild to find offer")
	}

	return offer, nil
}

func (c *productDatabase) FindAllOffer(ctx context.Context) ([]domain.Offer, error) {
	var offers []domain.Offer
	if c.DB.Raw("SELECT * FROM offers").Scan(&offers).Error != nil {
		return offers, errors.New("faild to get all offers")
	}
	return offers, nil
}

// save a new offer
func (c *productDatabase) SaveOffer(ctx context.Context, offer domain.Offer) error {

	query := `INSERT INTO offers (offer_name,description,discount_rate,start_date,end_date) VALUES ($1,$2,$3,$4,$5)`

	if c.DB.Exec(query, offer.OfferName, offer.Description, offer.DiscountRate, offer.StartDate, offer.EndDate).Error != nil {
		return errors.New("faild to save offer")
	}
	return nil
}

// update offer
func (c *productDatabase) UpdateOffer(ctx context.Context, offer domain.Offer) error {
	query := `UPDATE offers SET offer_name=$1,description=$2,discount_rate=$3,start_date=$4,end_date=$5 WHERE id=$6`
	if c.DB.Exec(query, offer.OfferName, offer.Description, offer.DiscountRate, offer.StartDate, offer.EndDate, offer.ID).Error != nil {
		return errors.New("faild to update offer")
	}
	return nil
}

// delet offer
func (c *productDatabase) DeleteOffer(ctx context.Context, offerID uint) error {
	trx := c.DB.Begin()
	// first update all discount price to 0 form products which are related by offer_products and offer_category
	query := `UPDATE products SET discount_price = 0  
	FROM offer_categories oc INNER JOIN offers o ON o.id = oc.offer_id 
	WHERE oc.category_id = products.category_id AND o.id = $1`
	if trx.Exec(query, offerID).Error != nil {
		trx.Rollback()
		return errors.New("faild to remove offer price from products")
	}

	// then upate all discount price on product_items discount price to 0 which are related by product_offer and category_offer
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

// find offer_category using category_id
func (c *productDatabase) FindOfferCategoryCategoryID(ctx context.Context, categoryID uint) (domain.OfferCategory, error) {
	var offerCategory domain.OfferCategory
	query := `SELECT * FROM offer_categories WHERE  category_id = ?`
	if c.DB.Raw(query, categoryID).Scan(&offerCategory).Error != nil {
		return offerCategory, errors.New("faild to find offer category")
	}
	return offerCategory, nil
}

func (c *productDatabase) FindOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) (domain.OfferCategory, error) {
	query := `SELECT * FROM offer_categories WHERE id = ? OR offer_id = ? AND category_id = ?`
	if c.DB.Raw(query, offerCategory.ID, offerCategory.OfferID, offerCategory.CategoryID).Scan(&offerCategory).Error != nil {
		return offerCategory, errors.New("faild to find offer category")
	}
	return offerCategory, nil
}

// find all offerCategory
func (c *productDatabase) FindAllOfferCategories(ctx context.Context) ([]res.ResOfferCategory, error) {
	var offerCategories []res.ResOfferCategory
	query := `SELECT oc.id AS offer_category_id, oc.category_id,c.category_name,oc.offer_id,o.offer_name FROM offer_categories oc 
	INNER JOIN categories c ON c.id = oc.category_id INNER JOIN offers o ON oc.offer_id = o.id`
	if c.DB.Raw(query).Scan(&offerCategories).Error != nil {
		return offerCategories, errors.New("faild to get all offer categories")
	}
	return offerCategories, nil
}

func (c *productDatabase) UpdateProductPriceForCategory(ctx context.Context) error {
	trx := c.DB.Begin()

	// update the all products discount price
	query := `UPDATE products SET discount_price = (price * (100 - o.discount_rate))/100 
	FROM offer_categories oc INNER JOIN offers o ON o.id = oc.offer_id 
	WHERE oc.category_id = products.category_id`

	if trx.Exec(query).Error != nil {
		trx.Rollback()
		return errors.New("faild to update discount price on product for category offer")
	}

	//remove discount price form removed offer_category
	query = `UPDATE products SET discount_price = 0 
	FROM offer_categories oc INNER JOIN offers o ON o.id = oc.offer_id 
	WHERE oc.category_id != products.category_id`

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

	// remove all discount price from product item where offer removed
	query = `UPDATE product_items AS pi SET discount_price = 0  
	FROM offer_categories oc INNER JOIN offers o ON o.id = oc.offer_id 
	INNER JOIN products p ON p.category_id = oc.category_id 
	WHERE p.id != pi.product_id`
	if trx.Exec(query).Error != nil {
		return errors.New("faild to update discount price on product_items for category offer")
	}
	// commit the transaction
	if trx.Commit().Error != nil {
		return errors.New("faild to complete the add offer category")
	}
	return nil
}

// save a new offer for category
func (c *productDatabase) SaveOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error {

	// first create the offer for category
	query := `INSERT INTO offer_categories (offer_id,category_id) VALUES ($1,$2)`
	if c.DB.Exec(query, offerCategory.OfferID, offerCategory.CategoryID).Error != nil {
		return errors.New("faild to save offer for category")
	}

	return c.UpdateProductPriceForCategory(ctx)
}

// remove offer from category

func (c *productDatabase) DeleteOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error {
	query := `DELETE FROM offer_categories WHERE id = $1 OR offer_id = $2 AND category_id = $3 `
	if c.DB.Exec(query, offerCategory.ID, offerCategory.OfferID, offerCategory.CategoryID).Error != nil {
		return errors.New("faild to remove offer from category on offer_cateogies")
	}
	return c.UpdateProductPriceForCategory(ctx)
}

// update offer category
func (c *productDatabase) UpdateOfferCategory(ctx context.Context, offerCategory domain.OfferCategory) error {

	query := `UPDATE offer_categories SET offer_id = $1 WHERE category_id = $2`
	if c.DB.Exec(query, offerCategory.OfferID, offerCategory.CategoryID).Error != nil {
		return errors.New("faild to update offer for category")
	}
	return c.UpdateProductPriceForCategory(ctx)
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
func (c *productDatabase) FindAllOfferProducts(ctx context.Context) ([]res.ResOfferProduct, error) {
	var offerProducts []res.ResOfferProduct
	query := `SELECT op.id AS offer_product_id, op.product_id,p.product_name,op.offer_id,o.offer_name FROM offer_products op INNER JOIN products p ON p.id = op.product_id 
	INNER JOIN offers o ON o.id = op.offer_id`
	if c.DB.Raw(query).Scan(&offerProducts).Error != nil {
		return offerProducts, errors.New("faild to find all offer products")
	}
	return offerProducts, nil
}

// save a offer for product
func (c *productDatabase) SaveOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error {
	query := `INSERT INTO offer_products (offer_id, product_id) VALUES ($1,$2)`
	if c.DB.Exec(query, offerProduct.OfferID, offerProduct.ProductID).Error != nil {
		return errors.New("faild to save offer for product")
	}
	return c.UpdateProductPriceForProductOffer(ctx)
}

func (c *productDatabase) DeleteOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error {
	// delete by id or offer_id and product_id
	query := `DELETE FROM offer_products WHERE id = $1 OR offer_id = $2 AND product_id = $3`
	if c.DB.Exec(query, offerProduct.ID, offerProduct.OfferID, offerProduct.ProductID).Error != nil {
		return errors.New("faild to remove offer from products on offer_products")
	}
	return c.UpdateProductPriceForProductOffer(ctx)
}

// update offer
func (c *productDatabase) UpdateOfferProduct(ctx context.Context, offerProduct domain.OfferProduct) error {
	query := `UPDATE offer_products SET offer_id = $1 WHERE product_id = $1`
	if c.DB.Exec(query, offerProduct.OfferID, offerProduct.ProductID).Error != nil {
		return errors.New("faild to update offer_products on offer update")
	}
	return c.UpdateProductPriceForProductOffer(ctx)
}
func (c *productDatabase) UpdateProductPriceForProductOffer(ctx context.Context) error {
	trx := c.DB.Begin()
	// first update all products discount price
	query := `UPDATE products p SET discount_price = (p.price * (100 - o.discount_rate))/100 
	FROM offers o INNER JOIN offer_products op ON o.id = op.offer_id 
	WHERE p.id = op.product_id`
	if trx.Exec(query).Error != nil {
		trx.Rollback()
		return errors.New("faild to update product discount price on offer_products")
	}

	//remove removed offer price removed product items
	query = `UPDATE products p SET discount_price = 0 
	FROM offers o INNER JOIN offer_products op ON o.id = op.offer_id 
	WHERE p.id != op.product_id`
	if trx.Exec(query).Error != nil {
		trx.Rollback()
		return errors.New("faild to remove product discount price on offer_products")
	}

	// then update product_items discont price
	query = `UPDATE product_items pi SET discount_price = (pi.price * (100 - o.discount_rate))/100 
	FROM offers o INNER JOIN offer_products op ON o.id = op.offer_id  
	WHERE pi.product_id = op.product_id`
	if trx.Exec(query).Error != nil {
		trx.Rollback()
		return errors.New("faild to update product_items discount price on offer_products")
	}

	//remove discount from removed offer productItems
	query = `UPDATE product_items pi SET discount_price = 0  
	FROM offers o INNER JOIN offer_products op ON o.id = op.offer_id  
	WHERE pi.product_id != op.product_id`
	if trx.Exec(query).Error != nil {
		trx.Rollback()
		return errors.New("faild to remove product_items discount price on offer_products")
	}

	//complete the updation
	if trx.Commit().Error != nil {
		trx.Rollback()
		return errors.New("faild to complete products price updation on offer_products")
	}
	return nil
}
