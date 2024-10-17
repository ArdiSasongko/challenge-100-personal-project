package helper

import (
	"basic-crud/model/domain"
	"basic-crud/model/web"
)

func ToPersonRespone(person domain.Person) web.PersonResponse {
	return web.PersonResponse{
		ID:   person.ID,
		Name: person.Name,
		Age:  person.Age,
	}
}

func ToAllPersonResponse(persons []domain.Person) []web.PersonResponse {
	var responses []web.PersonResponse
	for _, person := range persons {
		responses = append(responses, ToPersonRespone(person))
	}
	return responses
}
