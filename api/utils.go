package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	mapset "github.com/deckarep/golang-set"
	"github.com/research-pal/backend/db/notes"
)

// isAllowedQueryParam validates if the given parameter list has any params not supported.
// if found any, returns false along with the list of invalid params
// else returns true
func isAllowedQueryParam(params url.Values) (bool, []string) {
	validParams := mapset.NewSetFromSlice([]interface{}{"url", "assignee", "status", "group", "priority_order"})
	//TODO: need to make use of the common function isValid()
	incorrectQueryParam := []string{}

	for k := range params {
		if !validParams.Contains(k) {
			incorrectQueryParam = append(incorrectQueryParam, k)
		}
	}
	if len(incorrectQueryParam) > 0 {
		return false, incorrectQueryParam
	}
	return true, nil
}

func isValidPatchData(data map[string]interface{}) (bool, []string) {
	return isValid(data, mapset.NewSetFromSlice([]interface{}{"assignee", "status", "group", "priority_order", "notes"}))
}

func isValid(data map[string]interface{}, validFields mapset.Set) (bool, []string) {
	incorrectFields := []string{}
	for field := range data {
		if !validFields.Contains(field) {
			fmt.Printf("field:%v", field)
			incorrectFields = append(incorrectFields, field)
		}
	}

	if len(incorrectFields) > 0 {
		return false, incorrectFields
	}
	return true, nil
}

// convertToHTTPStatus converts a non nil error into the corresponding http status code
func convertToHTTPStatus(err error) int {
	switch {
	case errors.Is(err, notes.ErrorInvalidData) || errors.Is(err, notes.ErrorAlreadyExist):
		return http.StatusBadRequest
	case errors.Is(err, notes.ErrorNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
