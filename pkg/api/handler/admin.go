package handler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	usecaseInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
)

type adminHandler struct {
	adminUseCase usecaseInterface.AdminUseCase
}

func NewAdminHandler(adminUsecase usecaseInterface.AdminUseCase) interfaces.AdminHandler {
	return &adminHandler{
		adminUseCase: adminUsecase,
	}
}

// // AdminSignUp godoc
// // @summary api for admin to login
// // @id AdminSignUp
// // @tags Admin Login
// // @Param input body domain.Admin{} true "inputs"
// // @Router /admin/login [post]
// // @Success 200 {object} response.Response{} "successfully logged in"
// // @Failure 400 {object} response.Response{} "invalid input"
// // @Failure 500 {object} response.Response{} "failed to generate jwt token"
func (a *adminHandler) AdminSignUp(ctx *gin.Context) {

	var body domain.Admin

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	err := a.adminUseCase.SignUp(ctx, body)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create account for admin", err, nil)
		return
	}

	response.SuccessResponse(ctx, 200, "Successfully account created for admin", nil)
}

// GetAllUsers godoc
//
//	@Summary		Get all users
//	@Security		BearerAuth
//	@Description	API for admin to get all user details
//	@Id				GetAllUsers
//	@Tags			Admin User
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count"
//	@Router			/admin/users [get]
//	@Success		200	{object}	response.Response{}	"Successfully got all users"
//	@Success		204	{object}	response.Response{}	"No users found"
//	@Failure		500	{object}	response.Response{}	"Failed to find all users"
func (a *adminHandler) GetAllUsers(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	users, err := a.adminUseCase.FindAllUser(ctx, pagination)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all users", err, nil)
		return
	}

	if len(users) == 0 {
		response.SuccessResponse(ctx, http.StatusNoContent, "No users found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found all users", users)
}

// BlockUser godoc
//
//	@summary 	api for admin to block or unblock user
//	@Security	BearerAuth
//	@id			BlockUser
//	@tags		Admin User
//	@Param		input	body	request.BlockUser{}	true	"inputs"
//	@Router		/admin/users/block [patch]
//	@Success	200	{object}	response.Response{}	"Successfully changed block status of user"
//	@Failure	400	{object}	response.Response{}	"invalid input"
func (a *adminHandler) BlockUser(ctx *gin.Context) {

	var body request.BlockUser

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	err := a.adminUseCase.BlockOrUnBlockUser(ctx, body)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to change block status of user", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully changed block status of user")
}

// GetFullSalesReport godoc
//
//	@Summary		Get full sales report (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all sales report for a specific period in csv form
//	@id				GetFullSalesReport
//	@tags			Admin Sales
//	@Param			start_date	query	string	false	"Sales report starting date"
//	@Param			end_date	query	string	false	"Sales report ending date"
//	@Param			page_number	query	int		false	"Page Number"
//	@Param			count		query	int		false	"Count"
//	@Router			/admin/sales [get]
//	@Success		200	{object}	response.Response{}	"ecommerce_sales_report.csv"
//	@Success		204	{object}	response.Response{}	"No sales report found"
//	@Failure		500	{object}	response.Response{}	"failed to get sales report"
func (c *adminHandler) GetFullSalesReport(ctx *gin.Context) {

	// time
	startDate, err1 := utils.StringToTime(ctx.Query("start_date"))
	endDate, err2 := utils.StringToTime(ctx.Query("end_date"))

	// join all error and send it if its not nil
	err := errors.Join(err1, err2)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindQueryFailMessage, err1, nil)
		return
	}

	pagination := request.GetPagination(ctx)

	reqData := request.SalesReport{
		StartDate:  startDate,
		EndDate:    endDate,
		Pagination: pagination,
	}

	salesReport, err := c.adminUseCase.GetFullSalesReport(ctx, reqData)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get full sales report", err, nil)
		return
	}

	if len(salesReport) == 0 {
		response.SuccessResponse(ctx, http.StatusNoContent, "No sales report found", nil)
		return
	}

	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", "attachment;filename=ecommerce_sales_report.csv")

	csvWriter := csv.NewWriter(ctx.Writer)
	headers := []string{
		"UserID", "FirstName", "Email",
		"ShopOrderID", "OrderDate", "OrderTotalPrice",
		"Discount", "OrderStatus", "PaymentType",
	}

	if err := csvWriter.Write(headers); err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to write sales report on csv", err, nil)
		return
	}

	for _, sales := range salesReport {
		row := []string{
			fmt.Sprintf("%v", sales.UserID),
			sales.FirstName,
			sales.Email,
			fmt.Sprintf("%v", sales.ShopOrderID),
			sales.OrderDate.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("%v", sales.OrderTotalPrice),
			fmt.Sprintf("%v", sales.Discount),
			sales.OrderStatus,
			sales.PaymentType,
		}

		if err := csvWriter.Write(row); err != nil {
			response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create write sales report to csv", err, nil)
			return
		}
	}

	csvWriter.Flush()

}
