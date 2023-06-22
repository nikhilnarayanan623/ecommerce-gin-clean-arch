package request

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	PageNumber uint64
	Count      uint64
}

const (
	defaultPageNumber = 1
	defaultPageCount  = 10
)

func GetPagination(ctx *gin.Context) Pagination {

	pagination := Pagination{
		PageNumber: defaultPageNumber,
		Count:      defaultPageCount,
	}

	num, err := strconv.ParseUint(ctx.Query("page_number"), 10, 64)
	if err == nil {
		pagination.PageNumber = num
	}

	num, err = strconv.ParseUint(ctx.Query("count"), 10, 64)
	if err == nil {
		pagination.Count = num
	}
	return pagination
}
