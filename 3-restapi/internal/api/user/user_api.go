package userhandler

import (
	"net/http"
	"restapi/internal/model/payload"
	userservice "restapi/internal/service/user"
	"restapi/internal/utils"
	"strconv"

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
	router.HandleFunc("/login", h.LoginUser).Methods("POST")
	router.HandleFunc("/user/{id}", h.GetUser).Methods("GET")
	router.HandleFunc("/user/{id}", h.UpdateUser).Methods("PUT")
	router.HandleFunc("/user/{id}", h.DeleteUser).Methods("DELETE")
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

func (h *UserApi) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	userID, err := strconv.Atoi(idString)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "Invalid User ID", err)
		return
	}

	payload := new(payload.UpdateUser)
	if err := utils.ParseJSON(r, payload); err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "BAD REQUEST", err)
		return
	}

	result, err := h.Service.UpdateUser(r.Context(), *payload, userID)
	if err != nil {
		if err == userservice.ErrNotFound {
			utils.WriteErr(w, http.StatusNotFound, "NOT FOUND", err)
			return
		}
		utils.WriteErr(w, http.StatusBadRequest, "BAD REQUEST", err)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "Success Updated", result)
}

func (h *UserApi) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	userID, err := strconv.Atoi(idString)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "Invalid User ID", err)
		return
	}

	result, err := h.Service.FindByID(r.Context(), userID)
	if err != nil {
		if err == userservice.ErrNotFound {
			utils.WriteErr(w, http.StatusNotFound, "NOT FOUND", err)
			return
		}
		utils.WriteErr(w, http.StatusBadRequest, "BAD REQUEST", err)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "Success Updated", result)
}

func (h *UserApi) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	userID, err := strconv.Atoi(idString)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "Invalid User ID", err)
		return
	}

	err = h.Service.DeleteUser(r.Context(), userID)
	if err != nil {
		if err == userservice.ErrNotFound {
			utils.WriteErr(w, http.StatusNotFound, "NOT FOUND", err)
			return
		}
		utils.WriteErr(w, http.StatusBadRequest, "BAD REQUEST", err)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "Success Delete", nil)
}

func (h *UserApi) LoginUser(w http.ResponseWriter, r *http.Request) {
	payload := new(payload.LoginUser)
	if err := utils.ParseJSON(r, payload); err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "BAD REQUEST", err)
		return
	}

	result, err := h.Service.Login(r.Context(), *payload)
	if err != nil {
		if err == userservice.ErrInternal {
			utils.WriteErr(w, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
			return
		}
		utils.WriteErr(w, http.StatusBadRequest, "BAD REQUEST", err)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "Success Login", result)
}
