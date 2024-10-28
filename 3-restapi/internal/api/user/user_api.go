package userhandler

import (
	"net/http"
	"restapi/internal/model/payload"
	userservice "restapi/internal/service/user"
	"restapi/internal/utils"

	"github.com/gorilla/mux"
)

type UserApi struct {
	Service userservice.UserServiceInterface
}

func NewUserApi(Service userservice.UserServiceInterface) *UserApi {
	return &UserApi{
		Service: Service,
	}
}

func (h *UserApi) RegisterRouter(router *mux.Router) {
	router.HandleFunc("/register", h.Create).Methods("POST")
}

func (h *UserApi) Create(w http.ResponseWriter, r *http.Request) {
	payload := new(payload.CreateUser)
	if err := utils.ParseJSON(r, payload); err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "BAD REQUEST", err)
		return
	}

	result, err := h.Service.CreateUser(r.Context(), *payload)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "BAD REQUEST", err)
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, "Success Created", result)
}
