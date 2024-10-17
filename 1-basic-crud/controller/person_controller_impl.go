package controller

import (
	"basic-crud/model/web"
	"basic-crud/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type PersonControllerImpl struct {
	PersonService service.PersonServiceInterface
}

func NewPersonController(personService service.PersonServiceInterface) *PersonControllerImpl {
	return &PersonControllerImpl{
		PersonService: personService,
	}
}

func (c *PersonControllerImpl) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	personRequest := new(web.CreatePerson)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(personRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	person, err := c.PersonService.Create(r.Context(), *personRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   person,
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(webResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *PersonControllerImpl) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	personRequest := new(web.UpdatePerson)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(personRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	personID := p.ByName("personID")
	id, err := strconv.Atoi(personID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatePerson, err := c.PersonService.Update(r.Context(), *personRequest, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   updatePerson,
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(webResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *PersonControllerImpl) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	personID := p.ByName("personID")
	id, err := strconv.Atoi(personID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.PersonService.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(webResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *PersonControllerImpl) FindByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	personID := p.ByName("personID")
	id, err := strconv.Atoi(personID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	person, err := c.PersonService.FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   person,
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(webResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *PersonControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	persons, err := c.PersonService.FindAll(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   persons,
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(webResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
