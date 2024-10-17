package service

import (
	"basic-crud/model/web"
	"context"
)

type PersonServiceInterface interface {
	Create(ctx context.Context, request web.CreatePerson) (web.PersonResponse, error)
	Update(ctx context.Context, request web.UpdatePerson, personID int) (web.PersonResponse, error)
	Delete(ctx context.Context, personID int) error
	FindByID(ctx context.Context, personID int) (web.PersonResponse, error)
	FindAll(ctx context.Context) ([]web.PersonResponse, error)
}
