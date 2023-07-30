package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
)

func (c *middleware) TrimSpaces() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, int64(32<<20))

		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.Abort()
			response.ErrorResponse(ctx, http.StatusBadRequest, "failed to ready request body", err, nil)
			return
		}

		bodyBytes = trimSpacesInJson(bodyBytes)

		ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}
}

func trimSpacesInJson(data []byte) []byte {

	var mapData map[string]interface{}

	if err := json.Unmarshal(data, &mapData); err != nil {
		return data
	}

	for key, value := range mapData {
		// if the value is able to convert into string then trim its leading and tailing spaces
		if strValue, ok := value.(string); ok {
			mapData[key] = strings.TrimSpace(strValue)
		}
	}

	trimmedData, err := json.Marshal(mapData)
	if err != nil {
		return data
	}

	return trimmedData
}
