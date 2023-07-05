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

func (c *productDatabase) FindCategoryByName(ctx context.Context, categoryName string) (category domain.Category, err error) {

	query := `SELECT * FROM categories WHERE name = ?`
	err = c.DB.Raw(query, categoryName).Scan(&category).Error

	return
}

// Save Category
func (c *productDatabase) SaveCategory(ctx context.Context, categoryName string) (err error) {

	query := `INSERT INTO categories (name) VALUES ($1)`
	err = c.DB.Exec(query, categoryName).Error

	return err
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

// Find variation by category id and variation name
func (c *productDatabase) FindVariationByNameAndCategoryID(ctx context.Context,
	variationName string, categoryID uint) (variation domain.Variation, err error) {

	query := `SELECT id, name FROM variations WHERE category_id = $1 AND name = $2`
	err = c.DB.Raw(query, categoryID, variationName).Scan(&variation).Error

	return
}

// Find variation option by variation id and variation value
func (c *productDatabase) FindVariationOptionByValueAndVariationID(ctx context.Context,
	variationID uint, variationValue string) (variationOption domain.VariationOption, err error) {

	query := `SELECT id, value FROM variation_options WHERE variation_id = $1 AND value = $2`
	err = c.DB.Raw(query, variationID, variationValue).Scan(&variationOption).Error

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

func (c *productDatabase) FindProductByName(ctx context.Context, productName string) (product domain.Product, err error) {

	query := `SELECT * FROM products WHERE name = $1`
	err = c.DB.Raw(query, productName).Scan(&product).Error

	return
}

func (c *productDatabase) IsProductNameAlreadyExist(ctx context.Context, productName string) (exist bool, err error) {

	query := `SELECT EXISTS(SELECT 1 FROM products WHERE name = $1) AS exist FROM products`
	err = c.DB.Raw(query, productName).Scan(&exist).Error

	return
}

// to add a new product in database
func (c *productDatabase) SaveProduct(ctx context.Context, product request.Product) error {

	query := `INSERT INTO products (name, description, category_id, price, image, created_at) 
	VALUES($1, $2, $3, $4, $5, $6)`

	createdAt := time.Now()
	err := c.DB.Exec(query, product.Name, product.Description, product.CategoryID,
		product.Price, product.Image, createdAt).Error

	return err
}

// update product
func (c *productDatabase) UpdateProduct(ctx context.Context, product domain.Product) error {

	query := `UPDATE products SET name = $1, description = $2, category_id = $3, 
	price = $4, image = $5, updated_at = $6 
	WHERE id = $7`

	updatedAt := time.Now()

	err := c.DB.Exec(query, product.Name, product.Description, product.CategoryID,
		product.Price, product.Image, updatedAt, product.ID).Error

	return err
}

// get all products from database
func (c *productDatabase) FindAllProducts(ctx context.Context, pagination request.Pagination) (products []response.Product, err error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	query := `SELECT p.id, p.name, p.description, p.price, p.discount_price, p.image, p.category_id, 
	p.image, c.name AS category_name, p.created_at, p.updated_at  
	FROM products p LEFT JOIN categories c ON p.category_id = c.id 
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
	pi.qty_in_stock, sku FROM product_items pi 
	INNER JOIN products p ON p.id = pi.product_id 
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
