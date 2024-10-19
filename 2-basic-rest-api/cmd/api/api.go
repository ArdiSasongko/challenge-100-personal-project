package api

import (
	"basic-rest-api/repository"
	"basic-rest-api/service/user"
	"basic-rest-api/utils"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	validate := utils.NewValidate()
	userRepo := repository.NewUserRepository(s.db)
	userHandler := user.NewHandler(userRepo, validate)
	userHandler.RegisterRoutes(subrouter)

	log.Println("Server is running on localhost", s.addr)
	return http.ListenAndServe(s.addr, router)
}
