package handler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/auth"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type AdminHandler struct {
	adminUseCase service.AdminUseCase
}

func NewAdminHandler(adminUsecase service.AdminUseCase) *AdminHandler {
	return &AdminHandler{adminUseCase: adminUsecase}
}

// // AdminSignUp godoc
// // @summary api for admin to login
// // @id AdminSignUp
// // @tags Admin Login
// // @Param input body domain.Admin{} true "inputs"
// // @Router /admin/login [post]
// // @Success 200 {object} res.Response{} "successfully logged in"
// // @Failure 400 {object} res.Response{} "invalid input"
// // @Failure 500 {object} res.Response{} "faild to generate jwt token"
func (a *AdminHandler) AdminSignUp(ctx *gin.Context) {

	var admin domain.Admin

	if err := ctx.ShouldBindJSON(&admin); err != nil {
		response := res.ErrorResponse(400, "invlaid inputs", err.Error(), admin)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := a.adminUseCase.SignUp(ctx, admin)
	if err != nil {
		response := res.ErrorResponse(400, "faild to create account fo admin", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	respone := res.SuccessResponse(200, "successfully account created for admin", nil)
	ctx.JSON(http.StatusOK, respone)
}

// AdminLogin godoc
// @summary api for admin to login
// @id AdminLogin
// @tags Admin Login
// @Param input body req.LoginStruct{} true "inputs"
// @Router /admin/login [post]
// @Success 200 {object} res.Response{} "successfully logged in"
// @Failure 400 {object} res.Response{} "invalid input"
// @Failure 500 {object} res.Response{} "faild to generate jwt token"
func (a *AdminHandler) AdminLogin(ctx *gin.Context) {

	var body req.LoginStruct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// then check all field is empty
	if body.Email == "" && body.UserName == "" {
		err := errors.New("enter email or user_name atleast")
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var admin domain.Admin
	copier.Copy(&admin, &body)
	admin, err := a.adminUseCase.Login(ctx, admin)

	if err != nil {
		response := res.ErrorResponse(400, "faild to login", err.Error(), admin)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	tokenString, err := auth.GenerateJWT(admin.ID)
	if err != nil {
		response := res.ErrorResponse(500, "faild to generate jwt token", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.SetCookie("admin-auth", tokenString["accessToken"], 60*60, "", "", false, true)

	response := res.SuccessResponse(200, "successfully logged in", nil)
	ctx.JSON(http.StatusOK, response)
}

// AdminHome godoc
// @summary api admin home
// @id AdminHome
// @tags Admin Home
// @Router /admin [get]
// @Success 200 {object} res.Response{} "successfully logged in"
func (a *AdminHandler) AdminHome(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"message":    "Welcome to Admin Home",
	})
}

// ListUsers godoc
// @summary api for admin to list users
// @id ListUsers
// @tags Admin User
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /admin/users [get]
// @Success 200 {object} res.Response{} "successfully got all users"
// @Failure 500 {object} res.Response{} "faild to get all users"
func (a *AdminHandler) ListUsers(ctx *gin.Context) {

	count, err1 := utils.StringToUint(ctx.Query("count"))
	pageNumber, err2 := utils.StringToUint(ctx.Query("page_number"))

	err1 = errors.Join(err1, err2)
	if err1 != nil {
		response := res.ErrorResponse(400, "invalid inputs", err1.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	pagination := req.ReqPagination{
		PageNumber: pageNumber,
		Count:      count,
	}

	users, err := a.adminUseCase.FindAllUser(ctx, pagination)
	if err != nil {
		respone := res.ErrorResponse(500, "faild to get all users", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, respone)
		return
	}

	// check there is no usee
	if len(users) == 0 {
		response := res.SuccessResponse(200, "there is no users to show for this page", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got all users", users)
	ctx.JSON(http.StatusOK, response)

}

// BlockUser godoc
// @summary api for admin to block or unblock user
// @id BlockUser
// @tags Admin User
// @Param input body req.BlockStruct{} true "inputs"
// @Router /admin/users/block [patch]
// @Success 200 {object} res.Response{} "Successfully changed user block_status"
// @Failure 400 {object} res.Response{} "invalid input"
func (a *AdminHandler) BlockUser(ctx *gin.Context) {

	var body req.BlockStruct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := a.adminUseCase.BlockUser(ctx, body.UserID)
	if err != nil {
		response := res.ErrorResponse(400, "faild to change user block_status", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "Successfully changed user block_status", body)
	// if successfully blocked or unblock user then response 200
	ctx.JSON(http.StatusOK, response)
}

// FullSalesReport godoc
// @summary api for admin to see full sales report and download it as csv
// @id FullSalesReport
// @tags Admin Sales
// @Param start_date query string false "Date that you wan't to start on Report"
// @Param end_date query string false "Date that you wan't to start on Report"
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /admin/sales [get]
// @Success 200 {object} res.Response{} "ecommercesalesreport.csv"
// @Failure 500 {object} res.Response{} "faild to get sales report"
func (c *AdminHandler) FullSalesReport(ctx *gin.Context) {

	// time
	startDate, err1 := utils.StringToTime(ctx.Query("start_date"))
	endDate, err2 := utils.StringToTime(ctx.Query("end_date"))

	// page
	count, err3 := utils.StringToUint(ctx.Query("count"))
	pageNumber, err4 := utils.StringToUint(ctx.Query("page_number"))

	// join all error and send it if its not nil
	err1 = errors.Join(err1, err2, err3, err4)
	if err1 != nil {
		response := res.ErrorResponse(400, "invalid inputs", err1.Error(), req.ReqSalesReport{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	reqData := req.ReqSalesReport{
		StartDate: startDate,
		EndDate:   endDate,
		Pagination: req.ReqPagination{
			Count:      count,
			PageNumber: pageNumber,
		},
	}

	salesReport, err := c.adminUseCase.GetFullSalesReport(ctx, reqData)
	if err != nil {
		respones := res.ErrorResponse(500, "faild to get sales report", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, respones)
		return
	}

	if salesReport == nil {
		response := res.SuccessResponse(200, "there is no sales report on thi period", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", "attachment;filename=ecommercesalesreport.csv")

	csvWriter := csv.NewWriter(ctx.Writer)
	headers := []string{
		"UserID", "FirstName", "Email",
		"ShopOrderID", "OrderDate", "OrderTotalPrice",
		"Discount", "OrderStatus", "PaymentType",
	}

	if err := csvWriter.Write(headers); err != nil {
		response := res.ErrorResponse(500, "faild to reponse sales report", err.Error(), nil)
		ctx.JSON(500, response)
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
			response := res.ErrorResponse(500, "faild to create error csv", err.Error(), nil)
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}
	}

	csvWriter.Flush()

}
