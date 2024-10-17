package app

import (
	"basic-crud/controller"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(p controller.PersonControllerInterface) *httprouter.Router {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data := map[string]string{"message": "Welcome to the API"}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err := encoder.Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	router.POST("/api/person", p.Create)
	router.PUT("/api/person/:personID", p.Update)
	router.DELETE("/api/person/:personID", p.Delete)
	router.GET("/api/person/:personID", p.FindByID)
	router.GET("/api/persons", p.FindAll)

	return router
}
