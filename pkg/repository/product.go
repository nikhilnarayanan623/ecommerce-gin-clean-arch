package repository

import (
	"context"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
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

	return category, err
}

func (c *productDatabase) FindCategoryByID(ctx context.Context, categoryID uint) (category domain.Category, err error) {

	query := `SELECT * FROM categories WHERE id = ?`
	err = c.DB.Raw(query, categoryID).Scan(&category).Error

	return category, err
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

// Find all categories
func (c *productDatabase) FindAllCategories(ctx context.Context,
	pagination request.Pagination) (categories []response.Category, err error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	query := `SELECT id, name FROM categories LIMIT $1 OFFSET $2`
	err = c.DB.Raw(query, limit, offset).Scan(&categories).Error

	return categories, err
}

// Find variation by id
func (c *productDatabase) FindVariationByID(ctx context.Context, variationID uint) (variation domain.Variation, err error) {

	query := `SELECT id, name FROM variations WHERE id = $1`
	err = c.DB.Raw(query, variationID).Scan(&variation).Error

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
	variationOptionValue string, categoryID uint) (variationOption domain.VariationOption, err error) {

	query := `SELECT id, value FROM variation_options WHERE variation_id = $1 AND value = $2`
	err = c.DB.Raw(query, categoryID, variationOptionValue).Scan(&variationOption).Error

	return
}

// Save Variation for category
func (c *productDatabase) SaveVariation(ctx context.Context, variation request.Variation) error {

	query := `INSERT INTO variations (category_id, name) VALUES($1, $2)`
	err := c.DB.Exec(query, variation.CategoryID, variation.Name).Error

	return err
}


// add variation option
func (c *productDatabase) SaveVariationOption(ctx context.Context, variationOption request.VariationOption) error {

	query := `INSERT INTO variation_options (variation_id, value) VALUES($1, $2)`
	err := c.DB.Raw(query, variationOption.VariationID, variationOption.Value).Scan(&variationOption).Error

	return err
}

// find product by id
func (c *productDatabase) FindProductByID(ctx context.Context, productID uint) (product domain.Product, err error) {

	query := `SELECT * FROM products WHERE id = $1`
	err = c.DB.Raw(query, productID).Scan(&product).Error

	return product, err
}

func (c *productDatabase) FindProductByName(ctx context.Context, productName string) (product domain.Product, err error) {

	query := `SELECT * FROM products WHERE name = $1`
	err = c.DB.Raw(query, productName).Scan(&product).Error

	return product, err
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

	return products, err
}

// to get productItem id
func (c *productDatabase) FindProductItemByID(ctx context.Context, productItemID uint) (productItem domain.ProductItem, err error) {

	query := `SELECT * FROM product_items WHERE id = $1`
	err = c.DB.Raw(query, productItemID).Scan(&productItem).Error

	return productItem, err
}

func (c *productDatabase) IsProductItemAlreadyExist(ctx context.Context, productID, variationOptionID uint) (exist bool, err error) {

	query := `SELECT CASE WHEN id != 0 THEN 'T' ELSE 'F' END AS exist FROM product_items pi 
	INNER JOIN product_configurations pc  ON pi.id = pc.product_item_id AND pc.variation_option_id = $1 
	AND pi.product_id = $2`
	err = c.DB.Raw(query, variationOptionID, productID).Scan(&exist).Error

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
func (c *productDatabase) FindAllProductItems(ctx context.Context, productID uint) (productItems []response.ProductItems, err error) {

	// first find all product_items

	query := `SELECT p.name, pi.id,  pi.product_id, pi.price, pi.discount_price, 
	pi.qty_in_stock, sku FROM product_items pi 
	INNER JOIN products p ON p.id = pi.product_id 
	AND pi.product_id = $1`

	err = c.DB.Raw(query, productID).Scan(&productItems).Error

	return productItems, err
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

// func (c *productDatabase) FindAllProductItemImages(ctx context.Context, productItemID uint) (images []string, err error) {

// 	query := `SELECT image FROM product_images WHERE product_item_id = $1`

// 	var imagess domain.ProductImage

// 	err = c.DB.Raw(query, productItemID).Scan(&imagess).Error

// 	if err != nil {
// 		return images, fmt.Errorf("faild to find image of product_item with product_item_id %v", productItemID)
// 	}

// 	fmt.Println(imagess.Image)
// 	return images, nil
// }
