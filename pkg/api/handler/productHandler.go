package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type ProductHandler struct {
	productUseCase service.ProductUseCase
}

func NewProductHandler(productUsecase interfaces.ProductUseCase) *ProductHandler {
	return &ProductHandler{productUseCase: productUsecase}
}

func (p *ProductHandler) AllCategories(ctx *gin.Context) {

	categories, err := p.productUseCase.GetCategory(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "Faild to get categories",
			"error":      err.Error(),
		})
		return
	}

	// check category is empty or not
	if categories.Category == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"StatusCode": 200,
			"msg":        "there is no category to show",
		})
		return

	}
	ctx.JSON(http.StatusOK, gin.H{
		"StatsuCode":         200,
		"msg":                "Category Page",
		"categories":         categories.Category,
		"variations":         categories.VariationName,
		"variations options": categories.VariationValue,
	})
}

// add a category
func (p *ProductHandler) AddCategory(ctx *gin.Context) {

	var productCategory domain.Category

	if err := ctx.ShouldBindJSON(&productCategory); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatsuCode": 400,
			"msg":        "Error to bind the input",
			"error":      err.Error(),
		})
		return
	}

	err := p.productUseCase.AddCategory(ctx, productCategory)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatsuCode": 400,
			"msg":        "category can't be added",
			"err":        err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatsuCode": 200,
		"msg":        "category added",
	})
}

// for get all variations with its related category

// for add a variation like size / color/ ram/ memory
func (p *ProductHandler) VariationPost(ctx *gin.Context) {

	var variation domain.Variation
	if err := ctx.ShouldBindJSON(&variation); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't bind json",
			"error":      err.Error(),
		})
		return
	}

	variation, err := p.productUseCase.AddVariation(ctx, variation)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't add variation",
			"error":      err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "variation added",
		"variation":  variation,
	})
}

// for add a value for varaion color:blue,red; size:M,S,L; RAM:2gb,4gb;
func (p *ProductHandler) VariationOptionPost(ctx *gin.Context) {

	var variationOption domain.VariationOption
	if err := ctx.ShouldBindJSON(&variationOption); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Faild to bind json input",
			"error":      err.Error(),
		})
		return
	}

	variationOption, err := p.productUseCase.AddVariationOption(ctx, variationOption)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't add variation option",
			"error":      err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode":       200,
		"msg":              "Successfully variationOpion added",
		"variation_option": variationOption,
	})
}

// to show get all products
func (p *ProductHandler) ListProducts(ctx *gin.Context) {

	products, err := p.productUseCase.GetProducts(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "can't show products",
			"error":      err.Error(),
		})
		return
	}

	if products == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"StatusCode": 200,
			"msg":        "there is no products to show",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "List Products",
		"produtcts":  products,
	})

}

// AddProducts godoc
// @summary api for admin to update a product
// @id AddProducts
// @tags Products
// @Param input body req.ReqProduct{} true "inputs"
// @Router /admin/products [put]
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
// @tags Products
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
}

// AddProductItem godoc
// @summary api for admin to add product-items for a specific product
// @id AddProductItem
// @tags Products
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

	productItem, err := p.productUseCase.AddProductItem(ctx, body)

	if err != nil {
		response := res.ErrorResponse(400, "faild to add product_item", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully product item added", productItem)
	ctx.JSON(http.StatusOK, response)
}

// @summary api for get all product_items for a prooduct
// @id GetProductItems
// @tags Products
// @param product_id path int true "product_id"
// @Router /admin/products/product-items [get]
// @Success 200 {object} res.Response{} "successfully got all product_items for given product_id"
// @Failure 400 {object} res.Response{} "invalid input on params"
func (p *ProductHandler) GetProductItems(ctx *gin.Context) {

	productID, err := helper.StringToUint(ctx.Param("product_id"))

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
