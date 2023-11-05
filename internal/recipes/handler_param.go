package recipes

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func getStringQueryParam(ctx *gin.Context, queryParam string) (string, error) {
	strValue := ctx.Query(queryParam)
	if strValue == "" {
		return "", ErrRequiredParam
	}
	return strValue, nil
}

func getIntQueryParam(ctx *gin.Context, queryParam string) (int, error) {
	strValue := ctx.Query(queryParam)
	if strValue == "" {
		return 0, ErrRequiredParam
	}
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, ErrInvalidInt
	}
	return intValue, nil
}
