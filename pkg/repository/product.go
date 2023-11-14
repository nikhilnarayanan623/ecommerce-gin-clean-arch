package repository

import (
	"context"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type productDatabase struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) interfaces.ProductRepository {
	return &productDatabase{
		DB: db,
	}
}

func (c *productDatabase) Transactions(ctx context.Context, trxFn func(repo interfaces.ProductRepository) error) error {

	trx := c.DB.Begin()

	repo := NewProductRepository(trx)

	if err := trxFn(repo); err != nil {
		trx.Rollback()
		return err
	}

	if err := trx.Commit().Error; err != nil {
		trx.Rollback()
		return err
	}
	return nil
}

// To check the category name exist
func (c *productDatabase) IsCategoryNameExist(ctx context.Context, name string) (exist bool, err error) {

	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE name = $1 AND category_id IS NULL)`
	err = c.DB.Raw(query, name).Scan(&exist).Error

	return
}

// Save Category
func (c *productDatabase) SaveCategory(ctx context.Context, categoryName string) (err error) {

	query := `INSERT INTO categories (name) VALUES ($1)`
	err = c.DB.Exec(query, categoryName).Error

	return err
}

// To check the sub category name already exist for the category
func (c *productDatabase) IsSubCategoryNameExist(ctx context.Context, name string, categoryID uint) (exist bool, err error) {

	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE name = $1 AND category_id = $2)`
	err = c.DB.Raw(query, name, categoryID).Scan(&exist).Error

	return
}

// Save Category as sub category
func (c *productDatabase) SaveSubCategory(ctx context.Context, categoryID uint, categoryName string) (err error) {

	query := `INSERT INTO categories (category_id, name) VALUES ($1, $2)`
	err = c.DB.Exec(query, categoryID, categoryName).Error

	return err
}

// Find all main category(its not have a category_id)
func (c *productDatabase) FindAllMainCategories(ctx context.Context,
	pagination request.Pagination) (categories []response.Category, err error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	query := `SELECT id, name FROM categories WHERE category_id IS NULL 
	LIMIT $1 OFFSET $2`
	err = c.DB.Raw(query, limit, offset).Scan(&categories).Error

	return
}

// Find all sub categories of a category
func (c *productDatabase) FindAllSubCategories(ctx context.Context,
	categoryID uint) (subCategories []response.SubCategory, err error) {

	query := `SELECT id, name FROM categories WHERE category_id = $1`
	err = c.DB.Raw(query, categoryID).Scan(&subCategories).Error

	return
}

// Find all variations which related to given category id
func (c *productDatabase) FindAllVariationsByCategoryID(ctx context.Context,
	categoryID uint) (variations []response.Variation, err error) {

	query := `SELECT id, name FROM variations WHERE category_id = $1`
	err = c.DB.Raw(query, categoryID).Scan(&variations).Error

	return
}

// Find all variation options which related to given variation id
func (c productDatabase) FindAllVariationOptionsByVariationID(ctx context.Context,
	variationID uint) (variationOptions []response.VariationOption, err error) {

	query := `SELECT id, value FROM variation_options WHERE variation_id = $1`
	err = c.DB.Raw(query, variationID).Scan(&variationOptions).Error

	return
}

// To check a variation exist for the given category
func (c *productDatabase) IsVariationNameExistForCategory(ctx context.Context,
	name string, categoryID uint) (exist bool, err error) {

	query := `SELECT EXISTS(SELECT 1 FROM variations WHERE name = $1 AND category_id = $2)`
	err = c.DB.Raw(query, name, categoryID).Scan(&exist).Error

	return
}

// To check a variation value exist for the given variation
func (c *productDatabase) IsVariationValueExistForVariation(ctx context.Context,
	value string, variationID uint) (exist bool, err error) {

	query := `SELECT EXISTS(SELECT 1 FROM variation_options WHERE value = $1 AND variation_id = $2)`
	err = c.DB.Raw(query, value, variationID).Scan(&exist).Error

	return
}

// Save Variation for category
func (c *productDatabase) SaveVariation(ctx context.Context, categoryID uint, variationName string) error {

	query := `INSERT INTO variations (category_id, name) VALUES($1, $2)`
	err := c.DB.Exec(query, categoryID, variationName).Error

	return err
}

// add variation option
func (c *productDatabase) SaveVariationOption(ctx context.Context, variationID uint, variationValue string) error {

	query := `INSERT INTO variation_options (variation_id, value) VALUES($1, $2)`
	err := c.DB.Exec(query, variationID, variationValue).Error

	return err
}

// find product by id
func (c *productDatabase) FindProductByID(ctx context.Context, productID uint) (product domain.Product, err error) {

	query := `SELECT * FROM products WHERE id = $1`
	err = c.DB.Raw(query, productID).Scan(&product).Error

	return
}

func (c *productDatabase) IsProductNameExistForOtherProduct(ctx context.Context,
	name string, productID uint) (exist bool, err error) {

	query := `SELECT EXISTS(SELECT id FROM products WHERE name = $1 AND id != $2)`
	err = c.DB.Raw(query, name, productID).Scan(&exist).Error

	return
}

func (c *productDatabase) IsProductNameExist(ctx context.Context, productName string) (exist bool, err error) {

	query := `SELECT EXISTS(SELECT 1 FROM products WHERE name = $1)`
	err = c.DB.Raw(query, productName).Scan(&exist).Error

	return
}

// to add a new product in database
func (c *productDatabase) SaveProduct(ctx context.Context, product domain.Product) error {

	query := `INSERT INTO products (name, description, category_id, brand_id, price, image, created_at) 
	VALUES($1, $2, $3, $4, $5, $6, $7)`

	createdAt := time.Now()
	err := c.DB.Exec(query, product.Name, product.Description, product.CategoryID, product.BrandID,
		product.Price, product.Image, createdAt).Error

	return err
}

// update product
func (c *productDatabase) UpdateProduct(ctx context.Context, product domain.Product) error {

	query := `UPDATE products SET name = $1, description = $2, category_id = $3, 
	price = $4, image = $5, brand_id = $6, updated_at = $7 
	WHERE id = $8`

	updatedAt := time.Now()

	err := c.DB.Exec(query, product.Name, product.Description, product.CategoryID,
		product.Price, product.Image, product.BrandID, updatedAt, product.ID).Error

	return err
}

// get all products from database
func (c *productDatabase) FindAllProducts(ctx context.Context, pagination request.Pagination) (products []response.Product, err error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	query := `SELECT p.id, p.name, p.description, p.price, p.discount_price, 
	p.image, p.image, p.category_id, sc.name AS category_name, 
	mc.name AS main_category_name, p.brand_id, b.name AS brand_name,
	p.created_at, p.updated_at 
	FROM products p 
	INNER JOIN categories sc ON p.category_id = sc.id 
	INNER JOIN categories mc ON sc.category_id = mc.id 
	INNER JOIN brands b ON b.id = p.brand_id 
	ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	err = c.DB.Raw(query, limit, offset).Scan(&products).Error

	return
}

// to get productItem id
func (c *productDatabase) FindProductItemByID(ctx context.Context, productItemID uint) (productItem domain.ProductItem, err error) {

	query := `SELECT * FROM product_items WHERE id = $1`
	err = c.DB.Raw(query, productItemID).Scan(&productItem).Error

	return productItem, err
}

// to get how many variations are available for a product
func (c *productDatabase) FindVariationCountForProduct(ctx context.Context, productID uint) (variationCount uint, err error) {

	query := `SELECT COUNT(v.id) FROM variations v
	INNER JOIN categories c ON c.id = v.category_id 
	INNER JOIN products p ON p.category_id = v.category_id 
	WHERE p.id = $1`

	err = c.DB.Raw(query, productID).Scan(&variationCount).Error

	return
}

// To find all product item ids which related to the given product id and variation option id
func (c *productDatabase) FindAllProductItemIDsByProductIDAndVariationOptionID(ctx context.Context, productID,
	variationOptionID uint) (productItemIDs []uint, err error) {

	query := `SELECT id FROM product_items pi 
		INNER JOIN product_configurations pc ON pi.id = pc.product_item_id 
		WHERE pi.product_id = $1 AND variation_option_id = $2`
	err = c.DB.Raw(query, productID, variationOptionID).Scan(&productItemIDs).Error

	return
}

func (c *productDatabase) SaveProductConfiguration(ctx context.Context, productItemID, variationOptionID uint) error {

	query := `INSERT INTO product_configurations (product_item_id, variation_option_id) VALUES ($1, $2)`
	err := c.DB.Exec(query, productItemID, variationOptionID).Error

	return err
}

func (c *productDatabase) SaveProductItem(ctx context.Context, productItem domain.ProductItem) (productItemID uint, err error) {

	query := `INSERT INTO product_items (product_id, qty_in_stock, price, sku,created_at) VALUES($1, $2, $3, $4, $5) 
	 RETURNING id AS product_item_id`
	createdAt := time.Now()
	err = c.DB.Raw(query, productItem.ProductID, productItem.QtyInStock, productItem.Price, productItem.SKU, createdAt).
		Scan(&productItemID).Error

	return
}

// for get all products items for a product
func (c *productDatabase) FindAllProductItems(ctx context.Context,
	productID uint) (productItems []response.ProductItems, err error) {

	// first find all product_items

	query := `SELECT p.name, pi.id,  pi.product_id, pi.price, pi.discount_price, 
	pi.qty_in_stock, pi.sku, p.category_id, sc.name AS category_name, 
	mc.name AS main_category_name, p.brand_id, b.name AS brand_name 
	FROM product_items pi 
	INNER JOIN products p ON p.id = pi.product_id 
	INNER JOIN categories sc ON p.category_id = sc.id 
	INNER JOIN categories mc ON sc.category_id = mc.id 
	INNER JOIN brands b ON b.id = p.brand_id 
	AND pi.product_id = $1`

	err = c.DB.Raw(query, productID).Scan(&productItems).Error

	return
}

// Find all variation and value of a product item
func (c *productDatabase) FindAllVariationValuesOfProductItem(ctx context.Context,
	productItemID uint) (productVariationsValues []response.ProductVariationValue, err error) {

	query := `SELECT v.id AS variation_id, v.name, vo.id AS variation_option_id, vo.value 
	FROM  product_configurations pc 
	INNER JOIN variation_options vo ON vo.id = pc.variation_option_id 
	INNER JOIN variations v ON v.id = vo.variation_id 
	WHERE pc.product_item_id = $1`
	err = c.DB.Raw(query, productItemID).Scan(&productVariationsValues).Error

	return
}

// To save image for product item
func (c *productDatabase) SaveProductItemImage(ctx context.Context, productItemID uint, image string) error {

	query := `INSERT INTO product_images (product_item_id, image) VALUES ($1, $2)`
	err := c.DB.Exec(query, productItemID, image).Error

	return err
}

// To find all images of a product item
func (c *productDatabase) FindAllProductItemImages(ctx context.Context, productItemID uint) (images []string, err error) {

	query := `SELECT image FROM product_images WHERE product_item_id = $1`

	err = c.DB.Raw(query, productItemID).Scan(&images).Error

	return
}
