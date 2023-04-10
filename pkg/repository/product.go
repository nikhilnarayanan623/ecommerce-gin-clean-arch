package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
	"gorm.io/gorm"
)

type productDatabase struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) interfaces.ProductRepository {
	return &productDatabase{DB: db}
}

func (c *productDatabase) FindCategory(ctx context.Context, category domain.Category) (domain.Category, error) {

	if c.DB.Raw("SELECT * FROM categories WHERE id = ? OR category_name=?", category.ID, category.CategoryName).Scan(&category).Error != nil {
		return category, errors.New("faild to get category")
	}
	return category, nil
}

// add category
func (c *productDatabase) SaveCategory(ctx context.Context, category domain.Category) (err error) {

	// check the given category is main or sub
	if category.CategoryID == 0 { // no catogry id means its main category
		querry := `INSERT INTO categories (category_name)VALUES($1)`
		err = c.DB.Exec(querry, category.CategoryName).Error
	} else {
		//otherwise add its with main category
		querry := `INSERT INTO categories (category_id,category_name)VALUES($1,$2)`
		err = c.DB.Exec(querry, category.CategoryID, category.CategoryName).Error
	}

	if err != nil {
		return errors.New("faild to save category")
	}

	return nil
}

func (c *productDatabase) FindAllCategories(ctx context.Context) ([]res.Category, error) {
	var categories []res.Category
	querry := `SELECT s.id,s.category_name,s.category_id,m.category_name AS main_category_name 
	FROM categories s LEFT JOIN categories m ON s.category_id = m.id`
	if c.DB.Raw(querry).Scan(&categories).Error != nil {
		return categories, errors.New("faild to get categories")
	}
	return categories, nil
}

func (c *productDatabase) FindAllVariations(ctx context.Context) ([]res.VariationName, error) {
	var variationName []res.VariationName
	querry := `SELECT v.id,v.variation_name,v.category_id,c.category_name FROM variations v 
	LEFT JOIN categories c ON v.category_id=c.id`
	if c.DB.Raw(querry).Scan(&variationName).Error != nil {
		return variationName, errors.New("faild to get variations")
	}
	return variationName, nil
}

func (c productDatabase) FindAllVariationValues(ctx context.Context) ([]res.VariationValue, error) {
	var variaionValues []res.VariationValue
	querry := `SELECT vo.id,vo.variation_value,vo.variation_id,v.variation_name FROM variation_options vo 
	LEFT JOIN variations v ON vo.variation_id=v.id`
	if c.DB.Raw(querry).Scan(&variaionValues).Error != nil {
		return variaionValues, errors.New("faild to get variation options")
	}
	return variaionValues, nil
}

// add variation
func (c *productDatabase) AddVariation(ctx context.Context, variation domain.Variation) (domain.Variation, error) {

	//firs variation already exist or not
	c.DB.Raw("SELECT * FROM variations WHERE variation_name=?", variation.VariationName).Scan(&variation)
	if variation.ID != 0 {
		return variation, errors.New("variation already exist")
	}

	// then check the category provided for variaion is valid or not
	var cat domain.Category
	c.DB.Raw("SELECT * FROM categories WHERE id=?", variation.CategoryID).Scan(&cat)
	if cat.ID == 0 {
		return variation, errors.New("invalid category_id")
	}

	// if everything ok then add variation
	querry := `INSERT INTO variations (category_id,variation_name) VALUES($1,$2) RETURNING id, category_id,variation_name`
	if c.DB.Raw(querry, variation.CategoryID, variation.VariationName).Scan(&variation).Error != nil {
		return variation, errors.New("faild to add variation")
	}
	return variation, nil
}

// add variation option
func (c *productDatabase) AddVariationOption(ctx context.Context, variationOption domain.VariationOption) (domain.VariationOption, error) {

	// first check the variationOption already exist or not
	c.DB.Raw("SELECT * FROM variation_options WHERE variation_value=?", variationOption.VariationValue).Scan(&variationOption)
	if variationOption.ID != 0 {
		return variationOption, errors.New("given variation value already exist")
	}

	// then check the given variation is exist or not
	var variation domain.Variation
	c.DB.Raw("SELECT * FROM variations WHERE id=?", variationOption.VariationID).Scan(&variation)
	if variation.ID == 0 {
		return variationOption, errors.New("given variation dosen't exist")
	}

	//if everything ok then add the variation value
	querry := `INSERT INTO variation_options (variation_id,variation_value) VALUES($1,$2) RETURNING id, variation_id,variation_value`
	if c.DB.Raw(querry, variationOption.VariationID, variationOption.VariationValue).Scan(&variationOption).Error != nil {
		return variationOption, errors.New("faild to add variation value")
	}

	return variationOption, nil
}

// find product by id
func (c *productDatabase) FindProductByID(ctx context.Context, productID uint) (product domain.Product, err error) {
	query := `SELECT * FROM products WHERE id = $1`
	err = c.DB.Raw(query, productID).Scan(&product).Error
	if err != nil {
		return product, fmt.Errorf("faild find product with prduct_id %v", productID)
	}
	return product, nil
}

// find produc by name
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
func (c *productDatabase) FindAllProducts(ctx context.Context, pagination req.ReqPagination) (products []res.ResponseProduct, err error) {

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

// add a new product Items on database
func (c *productDatabase) SaveProductItem(ctx context.Context, reqProductItem req.ReqProductItem) error {

	trx := c.DB.Begin()

	var productItemItemID uint
	// first check the product item already exist

	querry := ` SELECT DISTINCT pi.id AS product_item_id FROM product_items pi INNER JOIN product_configurations pc ON pi.id = pc.product_item_id 
	WHERE pi.product_id= $1 AND pc.variation_option_id= $2`
	if trx.Raw(querry, reqProductItem.ProductID, reqProductItem.VariationOptionID).Scan(&productItemItemID).Error != nil {
		trx.Rollback()
		return errors.New("faild to check product_item already exist with this configuration")
	}
	fmt.Println(productItemItemID)
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
func (c *productDatabase) FindAllProductItems(ctx context.Context, productID uint) (productItems []res.RespProductItems, err error) {

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
