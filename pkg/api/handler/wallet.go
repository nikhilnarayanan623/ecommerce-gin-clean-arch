package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

func (c *OrderHandler) FindUserWallet(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	wallet, err := c.orderUseCase.FindUserWallet(ctx, userID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find user wallet", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found user wallet", wallet)
}

func (c *OrderHandler) FindUserWalletTransactions(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)
	pagination := request.GetPagination(ctx)

	transactions, err := c.orderUseCase.FindUserWalletTransactions(ctx, userID, pagination)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find user wallet transactions", err, nil)
		return
	}

	if len(transactions) == 0 {
		response.SuccessResponse(ctx, http.StatusOK, "No user wallet transaction found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found user wallet transactions", transactions)
}
