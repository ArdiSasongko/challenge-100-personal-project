package userhandler

import "net/http"

type UserApiInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
}
