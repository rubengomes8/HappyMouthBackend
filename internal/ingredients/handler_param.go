package ingredients

import (
	"net/http"
	"strconv"
)

func GetBoolQueryParamWithDefault(r *http.Request, queryParam string, defaultValue bool) (bool, error) {
	strValue := r.FormValue("sort-by-name")
	if strValue == "" {
		return defaultValue, nil
	}
	boolValue, err := strconv.ParseBool(strValue)
	if err != nil {
		return false, err
	}
	return boolValue, nil
}
