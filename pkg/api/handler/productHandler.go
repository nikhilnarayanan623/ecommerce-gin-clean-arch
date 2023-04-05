package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type ProductHandler struct {
	productUseCase service.ProductUseCase
}

func NewProductHandler(productUsecase interfaces.ProductUseCase) *ProductHandler {
	return &ProductHandler{productUseCase: productUsecase}
}

// GetAlllCategories godoc
// @summary api for adminn get all categories
// @security ApiKeyAuth
// @tags Admin Category
// @id GetAlllCategories
// @Router /admin/category [get]
// @Success 200 {object} res.Response{} "successfully update the coupon"
// @Failure 500 {object} res.Response{} "faild to get all cateogires"
func (p *ProductHandler) GetAlllCategories(ctx *gin.Context) {

	categories, err := p.productUseCase.GetCategory(ctx)

	if err != nil {
		response := res.ErrorResponse(500, "faild to get all cateogires", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// check category is empty or not
	if categories.Category == nil {
		response := res.SuccessResponse(200, "there is no category to show", nil)
		ctx.JSON(http.StatusOK, response)
		return

	}

	response := res.SuccessResponse(200, "successfully got all cateogries", categories)
	ctx.JSON(http.StatusOK, response)
}

// AddCategory godoc
// @summary api for adminn add a new category
// @security ApiKeyAuth
// @tags Admin Category
// @id AddCategory
// @Param input body domain.Category{} true "Input field"
// @Router /admin/category [post]
// @Success 200 {object} res.Response{} "successfully added a new category"
// @Failure 400 {object} res.Response{} "invalid input"
func (p *ProductHandler) AddCategory(ctx *gin.Context) {

	var productCategory domain.Category

	if err := ctx.ShouldBindJSON(&productCategory); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), productCategory)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := p.productUseCase.AddCategory(ctx, productCategory)
	if err != nil {
		response := res.ErrorResponse(400, "faild to add cateogy", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully added a new category")
	ctx.JSON(http.StatusOK, response)
}

// for get all variations with its related category
// for add a variation like size / color/ ram/ memory

// AddVariation godoc
// @summary api for adminn add a new variation
// @security ApiKeyAuth
// @tags Admin Category
// @id AddVariation
// @Param input body req.ReqVariation{} true "Input field"
// @Router /admin/category/variation [post]
// @Success 200 {object} res.Response{} "successfully variation added"
// @Failure 400 {object} res.Response{} "invalid input"
func (p *ProductHandler) AddVariation(ctx *gin.Context) {

	var body req.ReqVariation
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid inputs", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var variation domain.Variation
	copier.Copy(&variation, &body)

	variation, err := p.productUseCase.AddVariation(ctx, variation)

	if err != nil {
		respones := res.ErrorResponse(400, "faild to add variation", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, respones)
		return
	}

	respones := res.SuccessResponse(200, "successfully variation added", variation)
	ctx.JSON(http.StatusOK, respones)
}

// for add a value for varaion color:blue,red; size:M,S,L; RAM:2gb,4gb;

// AddVariationOption godoc
// @summary api for adminn add a new variation options
// @security ApiKeyAuth
// @tags Admin Category
// @id AddVariationOption
// @Param input body req.ReqVariationOption{} true "Input field"
// @Router /admin/category/variation-option [post]
// @Success 200 {object} res.Response{} "successfully added variation option"
// @Failure 400 {object} res.Response{} "invalid input"
func (p *ProductHandler) AddVariationOption(ctx *gin.Context) {

	var body req.ReqVariationOption
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var variationOption domain.VariationOption
	copier.Copy(&variationOption, body)
	variationOption, err := p.productUseCase.AddVariationOption(ctx, variationOption)

	if err != nil {
		response := res.ErrorResponse(400, "faild to add variation option", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully added variation option")
	ctx.JSON(http.StatusOK, response)
}

// ListProducts godoc
// @summary api for admin and user to show products
// @security ApiKeyAuth
// @tags User Products
// @id ListProducts
// @Router /products [get]
// @Success 200 {object} res.Response{} "successfully got all products"
// @Failure 500 {object} res.Response{}  "faild to get all products"
func (p *ProductHandler) ListProducts(ctx *gin.Context) {

	products, err := p.productUseCase.GetProducts(ctx)

	if err != nil {
		response := res.ErrorResponse(500, "faild to get all products", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if products == nil {
		response := res.SuccessResponse(200, "there is no products to show", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	respones := res.SuccessResponse(200, "successfully got all products", products)
	ctx.JSON(http.StatusOK, respones)

}

// AddProducts godoc
// @summary api for admin to update a product
// @id AddProducts
// @tags Admin Products
// @Param input body req.ReqProduct{} true "inputs"
// @Router /admin/products [post]
// @Success 200 {object} res.Response{} "successfully product added"
// @Failure 400 {object} res.Response{} "invalid input"
func (p *ProductHandler) AddProducts(ctx *gin.Context) {

	var body req.ReqProduct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		respones := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, respones)
		return
	}

	var product domain.Product
	copier.Copy(&product, &body)

	err := p.productUseCase.AddProduct(ctx, product)

	if err != nil {
		response := res.ErrorResponse(400, "faild to add product", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully product added", product)
	ctx.JSON(http.StatusOK, response)
}

// UpdateProduct godoc
// @summary api for admin to update a product
// @id UpdateProduct
// @tags Admin Products
// @Param input body req.ReqProduct{} true "inputs"
// @Router /admin/products [put]
// @Success 200 {object} res.Response{} "successfully product updated"
// @Failure 400 {object} res.Response{} "invalid input"
func (c *ProductHandler) UpdateProduct(ctx *gin.Context) {

	var body req.ReqProduct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var product domain.Product
	copier.Copy(&product, &body)

	err := c.productUseCase.UpdateProduct(ctx, product)
	if err != nil {
		response := res.ErrorResponse(400, "faild to update product", err.Error(), product)
		ctx.JSON(400, response)
		return
	}

	response := res.SuccessResponse(200, "successfully product updated", product)
	ctx.JSON(200, response)

	ctx.Abort()
}

// AddProductItem godoc
// @summary api for admin to add product-items for a specific product
// @id AddProductItem
// @tags Admin Products
// @Param input body req.ReqProductItem{} true "inputs"
// @Router /admin/products/product-items [post]
// @Success 200 {object} res.Response{} "Successfully product item added"
// @Failure 400 {object} res.Response{} "invalid input"
func (p *ProductHandler) AddProductItem(ctx *gin.Context) {
	// signle variation_value multiple images
	var body req.ReqProductItem

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	_, err := p.productUseCase.AddProductItem(ctx, body)

	if err != nil {
		response := res.ErrorResponse(400, "faild to add product_item", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully product item added", nil)
	ctx.JSON(http.StatusOK, response)
}

// @summary api for get all product_items for a prooduct
// @id GetProductItems
// @tags User Products
// @param product_id path int true "product_id"
// @Router /products/product-items [get]
// @Success 200 {object} res.Response{} "successfully got all product_items for given product_id"
// @Failure 400 {object} res.Response{} "invalid input on params"
func (p *ProductHandler) GetProductItems(ctx *gin.Context) {

	productID, err := utils.StringToUint(ctx.Param("product_id"))

	if err != nil {
		response := res.ErrorResponse(400, "invalid input on params", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	productItems, err := p.productUseCase.GetProductItems(ctx, productID)

	if err != nil {
		response := res.ErrorResponse(400, "faild to get product_items", err.Error(), nil)
		ctx.JSON(400, response)
		return
	}

	// check the product have productItem exist or not
	if productItems == nil {
		response := res.SuccessResponse(200, "there is no product items available for given product_id")
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got all product_items for given product_id", productItems)
	ctx.JSON(http.StatusOK, response)
}
