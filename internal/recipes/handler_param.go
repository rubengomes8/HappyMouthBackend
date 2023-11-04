package recipes

import (
	"github.com/gin-gonic/gin"
)

func getStringQueryParam(ctx *gin.Context, queryParam string) (string, error) {
	strValue := ctx.Query(queryParam)
	if strValue == "" {
		return "", ErrInvalidUserID
	}
	return strValue, nil
}
