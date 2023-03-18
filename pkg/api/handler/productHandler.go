package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
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

	category, err := p.productUseCase.AddCategory(ctx, productCategory)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatsuCode": 400,
			"msg":        "category can't be added",
			"err":        err.Error(),
		})
		return
	}

	var response res.RespCategory
	copier.Copy(&response, &category)

	ctx.JSON(http.StatusOK, gin.H{
		"StatsuCode": 200,
		"msg":        "category added",
		"categoty":   response,
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

// to add a new product
func (p *ProductHandler) AddProducts(ctx *gin.Context) {

	var body req.ReqProduct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatsuCode": 400,
			"msg":        "can't bind the input",
			"error":      err.Error(),
		})
		return
	}

	var product domain.Product
	copier.Copy(&product, &body)

	product, err := p.productUseCase.AddProduct(ctx, product)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatsuCode": 400,
			"msg":        "product can't be add",
			"err":        err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatsuCode": 200,
		"msg":        "successfully product added",
		"product":    product,
	})
}

// for add a specific product item
func (p *ProductHandler) AddProductItem(ctx *gin.Context) {
	// signle variation_value multiple images
	var reqProductItem req.ReqProductItem

	if err := ctx.ShouldBindJSON(&reqProductItem); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 200,
			"msg":        "can't bind the json",
			"error":      err.Error(),
		})
		return
	}

	productItem, err := p.productUseCase.AddProductItem(ctx, reqProductItem)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't add product item",
			"error":      err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode":   200,
		"msg":          "Successfully product item added",
		"product_item": productItem,
	})
}

func (p *ProductHandler) GetProductItems(ctx *gin.Context) {

	// product_id
	var body struct {
		ID uint `json:"id"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't bind the json",
			"errors":     err.Error(),
		})
		return
	}

	var product domain.Product

	copier.Copy(&product, &body)

	response, err := p.productUseCase.GetProductItems(ctx, product)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"error":      err.Error(),
		})
		return
	}

	// check the product have productItem exist or not
	if response == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"StatusCode": 200,
			"msg":        "there is no product items avialable for the product",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode":    200,
		"msg":           "Successfully product items got",
		"productItems ": response,
	})
}
