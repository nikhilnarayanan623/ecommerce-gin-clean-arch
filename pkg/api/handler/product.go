package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase"
	usecaseInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type ProductHandler struct {
	productUseCase usecaseInterface.ProductUseCase
}

func NewProductHandler(productUsecase usecaseInterface.ProductUseCase) interfaces.ProductHandler {
	return &ProductHandler{
		productUseCase: productUsecase,
	}
}

// FindAllCategories godoc
// @summary api for admin find all categories
// @security ApiKeyAuth
// @tags Admin Category
// @id FindAllCategories
// @Param page_number query int false "Page Number"
// @Param count query int false "Count"
// @Router /admin/category [get]
// @Success 200 {object} response.Response{} "Successfully found all categories"
// @Failure 500 {object} response.Response{} "Failed to find all categories"
func (p *ProductHandler) FindAllCategories(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	categories, err := p.productUseCase.FindAllCategories(ctx, pagination)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all categories", err, nil)
		return
	}

	if len(categories) == 0 {
		response.SuccessResponse(ctx, http.StatusOK, "No categories found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found all categories", categories)
}

// SaveCategory godoc
// @summary api for admin to add a new category
// @security ApiKeyAuth
// @id SaveCategory
// @tags Admin Category
// @Param input body domain.Category{} true "Input field"
// @Router /admin/category [post]
// @Success 200 {object} response.Response{} "Successfully category added"
// @Failure 400 {object} response.Response{} "invalid input"
func (p *ProductHandler) SaveCategory(ctx *gin.Context) {

	var body request.Category

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := p.productUseCase.SaveCategory(ctx, body.Name)

	if err != nil {

		statusCode := http.StatusInternalServerError

		if errors.Is(err, usecase.ErrCategoryAlreadyExist) {
			statusCode = http.StatusConflict
		}

		response.ErrorResponse(ctx, statusCode, "Failed to add category", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusCreated, "Successfully category added")
}

// SaveSubCategory godoc
// @summary api for admin add a new sub category
// @security ApiKeyAuth
// @id SaveSubCategory
// @tags Admin Category
// @Param input body domain.Category{} true "Input field"
// @Router /admin/category/sub-category [post]
// @Success 200 {object} response.Response{} "Successfully sub category added"
// @Failure 400 {object} response.Response{} "invalid input"
// @Failure 500 {object} response.Response{} "Failed to add sub category"
func (p *ProductHandler) SaveSubCategory(ctx *gin.Context) {

	var body request.SubCategory
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := p.productUseCase.SaveSubCategory(ctx, body)

	if err != nil {

		statusCode := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrCategoryAlreadyExist) {
			statusCode = http.StatusConflict
		}

		response.ErrorResponse(ctx, statusCode, "Failed to add sub category", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusCreated, "Successfully sub category added")
}

// SaveVariation godoc
// @summary api for admin add a new variation
// @security ApiKeyAuth
// @tags Admin Category
// @id SaveVariation
// @Param input body request.Variation{} true "Input field"
// @Router /admin/category/variation [post]
// @Success 200 {object} response.Response{} "successfully variation added"
// @Failure 400 {object} response.Response{} "invalid input"
func (p *ProductHandler) SaveVariation(ctx *gin.Context) {

	var body request.Variation

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := p.productUseCase.SaveVariation(ctx, body)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to add variation", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusCreated, "successfully variation added")
}

// SaveVariationOption godoc
// @summary api for admin add a new variation options
// @security ApiKeyAuth
// @tags Admin Category
// @id SaveVariationOption
// @Param input body request.VariationOption{} true "Input field"
// @Router /admin/category/variation-option [post]
// @Success 200 {object} response.Response{} "successfully added variation option"
// @Failure 400 {object} response.Response{} "invalid input"
func (p *ProductHandler) SaveVariationOption(ctx *gin.Context) {

	var body request.VariationOption

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := p.productUseCase.SaveVariationOption(ctx, body)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to add variation option", err, nil)
		return
	}
	response.SuccessResponse(ctx, http.StatusCreated, "Successfully variation option added")
}

// FindAllVariations godoc
// @summary api for admin to find all variations and its values for a specific category
// @security ApiKeyAuth
// @tags Admin Category
// @id FindAllVariations
// @Param category_id path  int true "Category ID"
// @Router /admin/category/variation/{category_id} [get]
// @Success 200 {object} response.Response{} "Successfully found all variations and its values"
// @Failure 400 {object} response.Response{} "invalid input"
func (c *ProductHandler) FindAllVariations(ctx *gin.Context) {

	categoryID, err := request.GetParamAsUint(ctx, "category_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	variations, err := c.productUseCase.FindAllVariationsAndItsValues(ctx, categoryID)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find variations and its values", err, nil)
		return
	}

	if len(variations) == 0 {
		response.SuccessResponse(ctx, http.StatusOK, "No variations found")
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found all variations and its values", variations)
}

// FindAllProductsAdmin godoc
// @summary api for admin to find all products
// @security ApiKeyAuth
// @tags Admin Products
// @id FindAllProductsAdmin
// @Param page_number query int false "Page Number"
// @Param count query int false "Count"
// @Router /admin/products [get]
// @Success 200 {object} response.Response{} "Successfully found all products"
// @Failure 500 {object} response.Response{}  "Failed to find all products"
func (p *ProductHandler) FindAllProductsAdmin() func(ctx *gin.Context) {
	return p.findAllProducts()
}

// FindAllProductsUser godoc
// @summary api for user to find all products
// @security ApiKeyAuth
// @tags User Products
// @id FindAllProductsUser
// @Param page_number query int false "Page Number"
// @Param count query int false "Count"
// @Router /products [get]
// @Success 200 {object} response.Response{} "Successfully found all products"
// @Failure 500 {object} response.Response{}  "Failed to find all products"
func (p *ProductHandler) FindAllProductsUser() func(ctx *gin.Context) {
	return p.findAllProducts()
}

// this is the common functionality of find product for admin and user
func (p *ProductHandler) findAllProducts() func(ctx *gin.Context) {

	return func(ctx *gin.Context) {
		pagination := request.GetPagination(ctx)

		products, err := p.productUseCase.FindAllProducts(ctx, pagination)

		if err != nil {
			response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all products", err, nil)
			return
		}

		if len(products) == 0 {
			response.SuccessResponse(ctx, http.StatusOK, "No products found", nil)
			return
		}

		response.SuccessResponse(ctx, http.StatusOK, "Successfully found all products", products)
	}

}

// SaveProduct godoc
// @summary api for admin to add a new product
// @id SaveProduct
// @tags Admin Products
// @Param input body request.Product{} true "inputs"
// @Router /admin/products [post]
// @Success 200 {object} response.Response{} "successfully product added"
// @Failure 400 {object} response.Response{} "invalid input"
func (p *ProductHandler) SaveProduct(ctx *gin.Context) {

	var body request.Product

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := p.productUseCase.SaveProduct(ctx, body)

	if err != nil {

		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to add product", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusCreated, "Successfully product added")
}

// UpdateProduct godoc
// @summary api for admin to update a product
// @id UpdateProduct
// @tags Admin Products
// @Param input body request.UpdateProduct{} true "inputs"
// @Router /admin/products [put]
// @Success 200 {object} response.Response{} "successfully product updated"
// @Failure 400 {object} response.Response{} "invalid input"
func (c *ProductHandler) UpdateProduct(ctx *gin.Context) {

	var body request.UpdateProduct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var product domain.Product
	copier.Copy(&product, &body)

	err := c.productUseCase.UpdateProduct(ctx, product)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update product", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully product updated", nil)
}

// SaveProductItem godoc
// @summary api for admin to add product item for a specific product
// @id SaveProductItem
// @tags Admin Products
// @Param input body request.ProductItem{} true "inputs"
// @Router /admin/products/product-item [post]
// @Success 200 {object} response.Response{} "Successfully product item added"
// @Failure 400 {object} response.Response{} "invalid input"
func (p *ProductHandler) SaveProductItem(ctx *gin.Context) {

	var body request.ProductItem

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := p.productUseCase.SaveProductItem(ctx, body)

	if err != nil {

		statusCode := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrProductItemAlreadyExist) {
			statusCode = http.StatusConflict
		}

		response.ErrorResponse(ctx, statusCode, "Failed to add product item", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusCreated, "Successfully product item added", nil)
}

// @summary api for admin to find all product items for a specific product
// @id FindAllProductItemsAdmin
// @tags User Products
// @param product_id path int true "product_id"
// @Router /admin/products/product-items/{product_id} [get]
// @Success 200 {object} response.Response{} "successfully got all product_items for given product_id"
// @Failure 400 {object} response.Response{} "invalid input on params"
func (p *ProductHandler) FindAllProductItemsAdmin() func(ctx *gin.Context) {
	return p.findAllProductItems()
}

// @summary api for user to find all product items for a specific produc
// @id FindAllProductItemsUser
// @tags User Products
// @param product_id path int true "product_id"
// @Router /products/product-items/{product_id} [get]
// @Success 200 {object} response.Response{} "successfully got all product_items for given product_id"
// @Failure 400 {object} response.Response{} "invalid input on params"
func (p *ProductHandler) FindAllProductItemsUser() func(ctx *gin.Context) {
	return p.findAllProductItems()
}

// same functionality of finding all product items for admin and user
func (p *ProductHandler) findAllProductItems() func(ctx *gin.Context) {

	return func(ctx *gin.Context) {

		productID, err := request.GetParamAsUint(ctx, "product_id")
		if err != nil {
			response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		}

		productItems, err := p.productUseCase.FindProductItems(ctx, productID)

		if err != nil {

			response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find product_items", err, nil)
			return
		}

		// check the product have productItem exist or not
		if len(productItems) == 0 {
			response.SuccessResponse(ctx, http.StatusOK, "No product items found")
			return
		}

		response.SuccessResponse(ctx, http.StatusOK, "Successfully found all product items ", productItems)
	}
}
