package errorhandling

import (
	"basic-crud/model/web"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func ErrorHandling(w http.ResponseWriter, r *http.Request, err interface{}) {
	if ValidationError(w, r, err) {
		return
	}
}

func ValidationError(w http.ResponseWriter, r *http.Request, err interface{}) bool {
	exception, OK := err.(validator.ValidationErrors)

	if OK {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   exception.Error(),
		}

		encoder := json.NewEncoder(w)
		errEncode := encoder.Encode(webResponse)
		if errEncode != nil {
			http.Error(w, errEncode.Error(), http.StatusInternalServerError)
		}
		return true
	}
	return false
}

func ErrorResponse(w http.ResponseWriter, code int, status string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	webResponse := web.WebResponse{
		Code:   code,
		Status: status,
		Data:   err.Error(),
	}
	json.NewEncoder(w).Encode(webResponse)
}
