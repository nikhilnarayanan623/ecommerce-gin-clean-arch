package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
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
// @summary api for admin get all categories
// @security ApiKeyAuth
// @tags Admin Category
// @id FindAllCategories
// @Router /admin/category [get]
// @Success 200 {object} response.Response{} "Successfully found all categories"
// @Failure 500 {object} response.Response{} "Failed to find all categories"
func (p *ProductHandler) FindAllCategories(ctx *gin.Context) {

	categories, err := p.productUseCase.FindCategory(ctx)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all categories", err, nil)
		return
	}

	if categories.Category == nil || len(categories.Category) == 0 {
		response.SuccessResponse(ctx, http.StatusOK, "No categories found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found all categories", categories)
}

// SaveCategory godoc
// @summary api for admin add a new category
// @security ApiKeyAuth
// @id SaveCategory
// @Param input body domain.Category{} true "Input field"
// @Router /admin/category [post]
// @Success 200 {object} res.Response{} "Successfully category added"
// @Failure 400 {object} res.Response{} "invalid input"
func (p *ProductHandler) SaveCategory(ctx *gin.Context) {

	var body domain.Category

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := p.productUseCase.SaveCategory(ctx, body)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to add category", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusCreated, "Successfully category added")
}

// for get all variations with its related category
// for add a variation like size / color/ ram/ memory

// SaveVariation godoc
// @summary api for admin add a new variation
// @security ApiKeyAuth
// @tags Admin Category
// @id SaveVariation
// @Param input body req.Variation{} true "Input field"
// @Router /admin/category/variation [post]
// @Success 200 {object} response.Response{} "successfully variation added"
// @Failure 400 {object} response.Response{} "invalid input"
func (p *ProductHandler) SaveVariation(ctx *gin.Context) {

	var body request.Variation

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var variation domain.Variation
	copier.Copy(&variation, &body)

	err := p.productUseCase.SaveVariation(ctx, variation)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to add variation", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusCreated, "successfully variation added")
}

// for add a value for variation color:blue,red; size:M,S,L; RAM:2gb,4gb;

// SaveVariationOption godoc
// @summary api for admin add a new variation options
// @security ApiKeyAuth
// @tags Admin Category
// @id SaveVariationOption
// @Param input body req.VariationOption{} true "Input field"
// @Router /admin/category/variation-option [post]
// @Success 200 {object} response.Response{} "successfully added variation option"
// @Failure 400 {object} response.Response{} "invalid input"
func (p *ProductHandler) SaveVariationOption(ctx *gin.Context) {

	var body request.VariationOption

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var variationOption domain.VariationOption
	copier.Copy(&variationOption, &body)

	err := p.productUseCase.SaveVariationOption(ctx, variationOption)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to add variation option", err, nil)
	}
	response.SuccessResponse(ctx, http.StatusCreated, "Successfully variation option added")
}

// FindAllProducts-Admin godoc
// @summary api for admin to show products
// @security ApiKeyAuth
// @tags Admin Products
// @id FindAllProducts-Admin
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /admin/products [get]
// @Success 200 {object} response.Response{} "successfully got all products"
// @Failure 500 {object} response.Response{}  "faild to get all products"

// FindAllProducts-User godoc
// @summary api for user to show products
// @security ApiKeyAuth
// @tags User Products
// @id FindAllProducts-User
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /products [get]
// @Success 200 {object} response.Response{} "successfully got all products"
// @Failure 500 {object} response.Response{}  "faild to get all products"
func (p *ProductHandler) FindAllProducts(ctx *gin.Context) {

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

// SaveProducts godoc
// @summary api for admin to update a product
// @id SaveProducts
// @tags Admin Products
// @Param input body req.Product{} true "inputs"
// @Router /admin/products [post]
// @Success 200 {object} response.Response{} "successfully product added"
// @Failure 400 {object} response.Response{} "invalid input"
func (p *ProductHandler) SaveProducts(ctx *gin.Context) {

	var body request.Product

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var product domain.Product
	copier.Copy(&product, &body)

	err := p.productUseCase.AddProduct(ctx, product)

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
// @Param input body req.UpdateProduct{} true "inputs"
// @Router /admin/products [put]
// @Success 200 {object} res.Response{} "successfully product updated"
// @Failure 400 {object} res.Response{} "invalid input"
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

	response.SuccessResponse(ctx, http.StatusOK, "Successfully product updated", body)
}

// SaveProductItem godoc
// @summary api for admin to add product-items for a specific product
// @id SaveProductItem
// @tags Admin Products
// @Param input body req.ProductItem{} true "inputs"
// @Router /admin/products/product-items [post]
// @Success 200 {object} res.Response{} "Successfully product item added"
// @Failure 400 {object} res.Response{} "invalid input"
func (p *ProductHandler) SaveProductItem(ctx *gin.Context) {

	var body request.ProductItem

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := p.productUseCase.AddProductItem(ctx, body)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to add product item", err, body)
		return
	}

	response.SuccessResponse(ctx, http.StatusCreated, "Successfully product item added", nil)
}

// @summary api for get all product_items for a product
// @id FindProductItems
// @tags User Products
// @param product_id path int true "product_id"
// @Router /products/product-items [get]
// @Success 200 {object} res.Response{} "successfully got all product_items for given product_id"
// @Failure 400 {object} res.Response{} "invalid input on params"
func (p *ProductHandler) FindProductItems(ctx *gin.Context) {

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
