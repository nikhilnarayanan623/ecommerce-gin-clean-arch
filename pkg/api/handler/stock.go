package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	usecaseInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type stockHandler struct {
	stockUseCase usecaseInterface.StockUseCase
}

func NewStockHandler(stockUseCase usecaseInterface.StockUseCase) interfaces.StockHandler {
	return &stockHandler{
		stockUseCase: stockUseCase,
	}
}

// GetAllStocks godoc
//	@Summary		Get all stocks (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all stocks
//	@Id				GetAllStocks
//	@Tags			Admin Stock
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count"
//	@Router			/admin/stocks [get]
//	@Success		200	{object}	response.Response{}	"Successfully found all stocks"
//	@Success		204	{object}	response.Response{}	"No stocks found"
//	@Failure		500	{object}	response.Response{}	"Failed to Get all stocks"
func (c *stockHandler) GetAllStocks(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	stocks, err := c.stockUseCase.GetAllStockDetails(ctx, pagination)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to Get all stocks", err, nil)
		return
	}

	if len(stocks) == 0 {
		response.SuccessResponse(ctx, http.StatusNoContent, "No stocks found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found all stocks", stocks)
}

// UpdateStock godoc
//	@Summary		Update stocks (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to update stock details
//	@Id				UpdateStock
//	@Tags			Admin Stock
//	@Param			input	body	request.UpdateStock{}	true	"Update stock details"
//	@Router			/admin/stocks [patch]
//	@Success		200	{object}	response.Response{}	"Successfully updated sock"
//	@Failure		400	{object}	response.Response{}	"Failed to bind input"
//	@Failure		500	{object}	response.Response{}	"Failed to update stock"
func (c *stockHandler) UpdateStock(ctx *gin.Context) {

	var body request.UpdateStock

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	err = c.stockUseCase.UpdateStockBySKU(ctx, body)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update stock", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully updated sock", nil)
}
