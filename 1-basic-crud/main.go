package main

import (
	"basic-crud/app"
	"basic-crud/controller"
	"basic-crud/repository"
	"basic-crud/service"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func NewServer(router http.Handler) *http.Server {
	return &http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}
}

func main() {
	DB := app.NewDatabase()
	Validate := app.NewValidator()
	repository := repository.NewPersonRepository()
	service := service.NewPersonService(repository, DB, Validate)
	controller := controller.NewPersonController(service)
	router := app.NewRouter(controller)

	server := NewServer(router)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
