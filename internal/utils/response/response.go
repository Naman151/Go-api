package response
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string
	Error string
}

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
   return Response{
	  Status: "Bad Request",
	  Error: err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Response{
		Status: "Bad Rwquest",
		Error: strings.Join(errMsg, ","),
	}
}
