package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type productDatabase struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) interfaces.ProductRepository {

	return &productDatabase{DB: db}
}

func (c *productDatabase) GetProducts(ctx context.Context) ([]helper.ResponseProduct, error) {

	var products []helper.ResponseProduct
	// aliase :: p := product; c := category
	querry := `SELECT p.id,p.product_name,p.description,p.price,p.image,p.category_id,p.image,c.category_name FROM products p LEFT JOIN categories c ON p.category_id=c.id`

	err := c.DB.Raw(querry).Scan(&products).Error

	return products, err

}
func (c *productDatabase) AddProducts(ctx context.Context, product domain.Product) (domain.Product, error) {

	// first check the product already exist
	var checkProduct domain.Product
	c.DB.Raw("SELECT * FROM products WHERE product_name=?", product.ProductName).Scan(&checkProduct)

	if checkProduct.ID != 0 {
		return product, errors.New("product already exist")
	}

	querry := `INSERT INTO products (product_name,description,category_id,price,image)VALUES($1,$2,$3,$4,$5) RETURNING id,product_name,description,category_id,price,image`
	err := c.DB.Raw(querry, product.ProductName, product.Description, product.CategoryID, product.Price, product.Image).Scan(&product).Error

	return product, err
}

// for get all products items for a product
func (c *productDatabase) GetProductItems(ctx context.Context, product domain.Product) ([]helper.RespProductItems, error) {

	var RespProductItems []helper.RespProductItems

	querry := `SELECT pi.id,pi.product_id,pi.price,pi.qty_in_stock,p.product_name FROM product_items pi INNER JOIN products p ON pi.product_id=p.id AND p.id=?`

	if c.DB.Raw(querry, product.ID).Scan(&RespProductItems).Error != nil {
		return RespProductItems, errors.New("faild to get product_items for product given product id")
	}

	// then get each productItems variationId from productCofiguration using productItemId
	querry = `SELECT vo.id AS variation_option_id,vo.variation_value FROM product_configurations pc JOIN variation_options vo ON pc.variation_option_id=vo.id AND pc.product_item_id=?`
	fmt.Println("herer")
	for i, productItem := range RespProductItems {

		c.DB.Raw(querry, productItem.ID).Scan(&productItem)

		RespProductItems[i].VariationOptionID = productItem.VariationOptionID
		RespProductItems[i].VariationValue = productItem.VariationValue
	}

	return RespProductItems, nil
}

func (c *productDatabase) AddProductItem(ctx context.Context, reqProductItem helper.ReqProductItem) (domain.ProductItem, error) {
	var productItem domain.ProductItem
	// first check the product item already exist

	querry := `SELECT * FROM product_items p JOIN product_configurations pc ON p.id=pc.product_item_id AND pc.variation_option_id=?`
	if c.DB.Raw(querry, reqProductItem.VariationOptionID).Scan(&productItem).Error != nil {
		return productItem, errors.New("faild to check product item")
	}

	if productItem.ID != 0 {
		return productItem, errors.New("this product configuration already exist")
	}

	//then insert product_id ,quantity and price
	// var productItem domain.ProductItem
	querry = `INSERT INTO product_items (product_id,qty_in_stock,price) VALUES ($1,$2,$3) RETURNING id, product_id, qty_in_stock, price`
	if c.DB.Raw(querry, reqProductItem.ProductID, reqProductItem.QtyInStock, reqProductItem.Price).Scan(&productItem).Error != nil {
		return productItem, errors.New("faild to add product item")
	}

	// add all images in db with this productItemID
	var producImage []domain.ProductImage
	querry = `INSERT INTO product_images (product_item_id,image) VALUES ($1,$2) RETURNING id`

	// loop to insert all images from the array
	for _, img := range reqProductItem.Images {
		if c.DB.Raw(querry, productItem.ID, img).Scan(&producImage).Error != nil {
			return productItem, errors.New("faild to add image on db")
		}
	}

	// atlast cofigure productItem in productConfiguration
	var pCofig domain.ProductConfiguration
	querry = `INSERT INTO product_configurations (product_item_id,variation_option_id) VALUES ($1,$2) RETURNING product_item_id,variation_option_id`
	if c.DB.Raw(querry, productItem.ID, reqProductItem.VariationOptionID).Scan(&pCofig).Error != nil {
		return productItem, errors.New("faild to add product configuration on db")
	}

	return productItem, nil
}

// to get all categories, all variations and all variation value
func (c *productDatabase) GetCategory(ctx context.Context) (helper.RespFullCategory, error) {

	var response helper.RespFullCategory

	// first find all categories (aliase :: s:= category(means sub category); m:= category (main category) )
	querry := `SELECT s.id,s.category_name,s.category_id,m.category_name AS main_category_name FROM categories s LEFT JOIN categories m ON s.category_id = m.id`
	if c.DB.Raw(querry).Scan(&response.Category).Error != nil {
		return response, errors.New("faild to get categories")
	}

	// find all variations (aliase ::  v := variations; c:= category)
	querry = `SELECT v.id,v.variation_name,v.category_id,c.category_name FROM variations v LEFT JOIN categories c ON v.category_id=c.id`
	if c.DB.Raw(querry).Scan(&response.VariationName).Error != nil {
		return response, errors.New("faild to get variations")
	}

	// find all variations value (aliase :: vo:= variation_options; v:= variations)
	querry = `SELECT vo.id,vo.variation_value,vo.variation_id,v.variation_name FROM variation_options vo LEFT JOIN variations v ON vo.variation_id=v.id`
	if c.DB.Raw(querry).Scan(&response.VariationValue).Error != nil {
		return response, errors.New("faild to get variation options")
	}

	return response, nil
}

// add category
func (c *productDatabase) AddCategory(ctx context.Context, category domain.Category) (domain.Category, error) {

	var checkCat domain.Category
	//first check the categoryname already exisits or not
	c.DB.Raw("SELECT * FROM categories WHERE category_name=?", category.CategoryName).Scan(&checkCat)

	if checkCat.ID != 0 { // means category already exist
		return category, errors.New("category already exist")
	}

	// check the given category is main or sub
	if category.CategoryID == 0 { // no catogry id means its main category
		querry := `INSERT INTO categories (category_name)VALUES($1) RETURNING id, category_id, category_name`
		c.DB.Raw(querry, category.CategoryName).Scan(&category)
	} else {
		// first check the category id is valid or not
		c.DB.Raw("SELECT * FROM categories WHERE id=?", category.CategoryID).Scan(&checkCat)
		if checkCat.ID == 0 { // its not a valid category
			return category, errors.New("category_id is invalid main category not exist")
		}
		//otherwise add its with main category
		querry := `INSERT INTO categories (category_id,category_name)VALUES($1,$2) RETURNING id, category_id,category_name`
		c.DB.Raw(querry, category.CategoryID, category.CategoryName).Scan(&category)
	}

	return category, nil
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
