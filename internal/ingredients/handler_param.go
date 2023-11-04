package ingredients

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func getBoolQueryParamWithDefault(ctx *gin.Context, queryParam string, defaultValue bool) (bool, error) {
	strValue := ctx.Query(queryParam)
	if strValue == "" {
		return defaultValue, nil
	}
	boolValue, err := strconv.ParseBool(strValue)
	if err != nil {
		return false, err
	}
	return boolValue, nil
}
