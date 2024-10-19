package user

import (
	"basic-rest-api/model/domain"
	"basic-rest-api/model/web"
	"basic-rest-api/repository"
	"basic-rest-api/service/auth"
	"basic-rest-api/utils"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	repository repository.UserRepoInterface
	validate   *validator.Validate
}

func NewHandler(repository repository.UserRepoInterface, validate *validator.Validate) *Handler {
	return &Handler{
		repository: repository,
		validate:   validate,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Implement login logic here
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// Parse JSON payload
	payload := new(web.RegisterUser)
	if err := utils.ParseJSON(r, payload); err != nil {
		utils.WriteErr(w, http.StatusBadRequest, err)
		return
	}

	// Validate payload
	err := h.validate.Struct(payload)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErr(w, http.StatusBadRequest, fmt.Errorf("invalid payload %s", errors))
		return
	}

	// Check if user already exists
	_, err = h.repository.GetUserByEmail(payload.Email)

	// If no error, it means the user exists
	if err == nil {
		utils.WriteErr(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	// If the error is not sql.ErrNoRows, it means there's a database issue
	if err != sql.ErrNoRows {
		utils.WriteErr(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	// Hash password
	hash, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteErr(w, http.StatusInternalServerError, fmt.Errorf("failed to generate hash password"))
		return
	}

	// Save user to database
	err = h.repository.CreateUser(domain.User{
		Username: payload.Username,
		Password: hash,
		Email:    payload.Email,
		Name:     payload.Name,
	})

	if err != nil {
		utils.WriteErr(w, http.StatusInternalServerError, err)
		return
	}

	// Send success response with user ID
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "user created successfully",
	})
}
