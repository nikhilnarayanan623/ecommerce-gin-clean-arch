package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
	"gorm.io/gorm"
)

type productDatabase struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) interfaces.ProductRepository {
	return &productDatabase{DB: db}
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

// find product by name
func (c *productDatabase) FindProduct(ctx context.Context, product domain.Product) (domain.Product, error) {

	if c.DB.Raw("SELECT * FROM products WHERE id = ? OR product_name=?", product.ID, product.ProductName).Scan(&product).Error != nil {
		return product, errors.New("faild to get product")
	}
	return product, nil
}

// to add a new product in database
func (c *productDatabase) SaveProduct(ctx context.Context, product domain.Product) error {

	querry := `INSERT INTO products (product_name, description, category_id, price, image, created_at) 
	VALUES($1, $2, $3, $4, $5, $6)`

	createdAt := time.Now()
	if c.DB.Exec(querry, product.ProductName, product.Description, product.CategoryID,
		product.Price, product.Image, createdAt).Error != nil {
		return errors.New("faild to insert product on database")
	}

	return nil
}

// update product
func (c *productDatabase) UpdateProduct(ctx context.Context, product domain.Product) error {
	query := `UPDATE products SET product_name = $1, description = $2, category_id = $3, 
	price = $4, image = $5, updated_at = $6 WHERE id = $7`

	updatedAt := time.Now()

	if c.DB.Exec(query, product.ProductName, product.Description, product.CategoryID,
		product.Price, product.Image, updatedAt, product.ID).Error != nil {
		return errors.New("faild to update product")
	}

	return nil
}

// get all products from database
func (c *productDatabase) FindAllProducts(ctx context.Context, pagination request.Pagination) (products []response.Product, err error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	// aliase :: p := product; c := category
	querry := `SELECT p.id, p.product_name, p.description, p.price, p.discount_price, p.image, p.category_id, 
	p.image, c.category_name, p.created_at, p.updated_at  
	FROM products p LEFT JOIN categories c ON p.category_id=c.id 
	ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	if c.DB.Raw(querry, limit, offset).Scan(&products).Error != nil {
		return products, errors.New("faild to get products from database")
	}

	return products, nil
}

// to get productItem id
func (c *productDatabase) FindProductItem(ctx context.Context, productItemID uint) (productItem domain.ProductItem, err error) {

	query := `SELECT * FROM product_items WHERE id = $1`
	err = c.DB.Raw(query, productItemID).Scan(&productItem).Error

	return productItem, err
}

// add a new product Items on database
func (c *productDatabase) SaveProductItem(ctx context.Context, reqProductItem request.ProductItem) error {

	trx := c.DB.Begin()

	var productItemItemID uint

	querry := ` SELECT DISTINCT pi.id AS product_item_id FROM product_items pi INNER JOIN product_configurations pc ON pi.id = pc.product_item_id 
	WHERE pi.product_id= $1 AND pc.variation_option_id= $2`
	if trx.Raw(querry, reqProductItem.ProductID, reqProductItem.VariationOptionID).Scan(&productItemItemID).Error != nil {
		trx.Rollback()
		return errors.New("faild to check product_item already exist with this configuration")
	}

	// if product item already exist with this productId
	if productItemItemID != 0 {
		trx.Rollback()
		return fmt.Errorf("a product_item already for this product \nwith given configuration as product_item_id %v", productItemItemID)
	}

	// insert the product_item
	createdAt := time.Now()
	sku := utils.GenerateSKU()

	querry = `INSERT INTO product_items (product_id, qty_in_stock, price, sku, created_at) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id AS product_item_id`
	err := c.DB.Raw(querry, reqProductItem.ProductID, reqProductItem.QtyInStock, reqProductItem.Price, sku, createdAt).Scan(&productItemItemID).Error
	if err != nil {
		trx.Rollback()
		return fmt.Errorf("faild to save product_item for product with product_id %v", reqProductItem.ProductID)
	}

	querry = `INSERT INTO product_images (product_item_id,image) VALUES ($1,$2)`
	// loop to insert all images from the array
	for _, img := range reqProductItem.Images {

		err = c.DB.Exec(querry, productItemItemID, img).Error
		if err != nil {
			trx.Rollback()
			return fmt.Errorf("faid to add image for product_item of product with product_id %v", reqProductItem.ProductID)
		}
	}

	querry = `INSERT INTO product_configurations (product_item_id, variation_option_id) VALUES ($1, $2)`

	err = c.DB.Exec(querry, productItemItemID, reqProductItem.VariationOptionID).Error
	if err != nil {
		trx.Rollback()
		return fmt.Errorf("faild to add product configuration of product_item for product with prodcut_id %v", reqProductItem.ProductID)
	}

	err = trx.Commit().Error
	if err != nil {
		trx.Rollback()
		return fmt.Errorf("faild to complete the product_item save for product with product_id %v", reqProductItem.ProductID)
	}

	return nil
}

// for get all products items for a product
func (c *productDatabase) FindAllProductItems(ctx context.Context, productID uint) (productItems []response.ProductItems, err error) {

	// first find all product_items

	query := `SELECT p.product_name, pi.id,  pi.product_id, pi.price, pi.discount_price, 
	pi.qty_in_stock, sku, vo.id AS variation_option_id, vo.variation_value 
	FROM product_items pi INNER JOIN products p ON p.id = pi.product_id 
	INNER JOIN product_configurations pc ON pc.product_item_id = pi.id 
	INNER JOIN variation_options vo ON vo.id = pc.variation_option_id 
	AND pi.product_id = $1`

	err = c.DB.Raw(query, productID).Scan(&productItems).Error
	if err != nil {
		return productItems, fmt.Errorf("faild to find all product_items for product_id %v", productID)
	}

	return productItems, nil
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
