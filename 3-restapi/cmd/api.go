package api

import (
	"database/sql"
	"log"
	"net/http"
	userhandler "restapi/internal/api/user"
	userrepository "restapi/internal/repository/user"
	userservice "restapi/internal/service/user"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type ApiServer struct {
	add string
	db  *sql.DB
}

func NewServerAPI(add string, db *sql.DB) *ApiServer {
	return &ApiServer{
		add: add,
		db:  db,
	}
}

func (api *ApiServer) Run() error {
	router := mux.NewRouter()

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	validate := *validator.New()
	userRepo := userrepository.NewUserRepository()
	userService := userservice.NewUserService(userRepo, &validate, api.db)
	userHandler := userhandler.NewUserApi(userService)
	userHandler.RegisterRouter(subrouter)

	log.Println("Server is running on localhost", api.add)
	return http.ListenAndServe(api.add, router)
}
