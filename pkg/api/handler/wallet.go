package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

func (c *OrderHandler) GetUserWallet(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	wallet, err := c.orderUseCase.GetUserWallet(ctx, userID)

	if err != nil {
		response := res.ErrorResponse(500, "faild to get user wallet", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got user wallet", wallet)

	ctx.JSON(http.StatusOK, response)
}

func (c *OrderHandler) GetUserWalletTransactions(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	count, err1 := utils.StringToUint(ctx.Query("count"))
	pageNumber, err2 := utils.StringToUint(ctx.Query("page_number"))

	err1 = errors.Join(err1, err2)
	if err1 != nil {
		response := res.ErrorResponse(400, "invalid inputs", err1.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	pagination := req.Pagination{
		PageNumber: pageNumber,
		Count:      count,
	}

	transactions, err := c.orderUseCase.GetUserWalletTransactions(ctx, userID, pagination)

	if err != nil {
		response := res.ErrorResponse(500, "faild to get user wallet transactions", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if transactions == nil {
		resonse := res.SuccessResponse(200, "there is no wallet transaction for this page", nil)
		ctx.JSON(http.StatusOK, resonse)
		return
	}

	response := res.SuccessResponse(200, "successfully got user wallet transactions", transactions)
	ctx.JSON(http.StatusOK, response)
}
