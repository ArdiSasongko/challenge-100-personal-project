package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"restapi/internal/model/web"
)

func ParseJSON(r *http.Request, payload interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteErr(w http.ResponseWriter, status int, message string, err error) {
	WriteJSON(w, status, web.FailedResponse{
		StatusCode: status,
		Message:    message,
		Error:      err,
	})
}

func WriteSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	WriteJSON(w, status, web.SuccessResponse{
		StatusCode: status,
		Message:    message,
		Data:       data,
	})
}
