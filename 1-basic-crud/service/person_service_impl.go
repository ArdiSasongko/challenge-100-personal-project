package service

import (
	"basic-crud/helper"
	"basic-crud/model/domain"
	"basic-crud/model/web"
	"basic-crud/repository"
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
)

type PersonServiceImpl struct {
	PersonRepository repository.PersonRepositoryInterface
	DB               *sql.DB
	Validate         *validator.Validate
}

func NewPersonService(personRepository repository.PersonRepositoryInterface, db *sql.DB, validate *validator.Validate) *PersonServiceImpl {
	return &PersonServiceImpl{
		PersonRepository: personRepository,
		DB:               db,
		Validate:         validate,
	}
}

func (s *PersonServiceImpl) Create(ctx context.Context, request web.CreatePerson) (web.PersonResponse, error) {
	err := s.Validate.Struct(request)
	if err != nil {
		return web.PersonResponse{}, err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return web.PersonResponse{}, err
	}
	defer helper.RollBackCommit(tx)

	newPerson := domain.Person{
		Name: request.Name,
		Age:  request.Age,
	}

	result, err := s.PersonRepository.Create(ctx, tx, newPerson)

	if err != nil {
		return web.PersonResponse{}, err
	}

	return helper.ToPersonRespone(*result), nil
}

func (s *PersonServiceImpl) Update(ctx context.Context, request web.UpdatePerson, personID int) (web.PersonResponse, error) {
	err := s.Validate.Struct(request)
	if err != nil {
		return web.PersonResponse{}, err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return web.PersonResponse{}, err
	}
	defer helper.RollBackCommit(tx)

	person, err := s.PersonRepository.FindByID(ctx, tx, personID)
	if err != nil {
		return web.PersonResponse{}, err
	}

	valueName := helper.DefaultValue(request.Name, person.Name)
	valueAge := helper.DefaultValue(request.Age, person.Age)

	person.Name = valueName
	person.Age = valueAge

	result, err := s.PersonRepository.Update(ctx, tx, *person)

	if err != nil {
		return web.PersonResponse{}, err
	}

	return helper.ToPersonRespone(*result), nil
}

func (s *PersonServiceImpl) Delete(ctx context.Context, personID int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.RollBackCommit(tx)

	_, err = s.PersonRepository.FindByID(ctx, tx, personID)
	if err != nil {
		return err
	}

	err = s.PersonRepository.Delete(ctx, tx, personID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PersonServiceImpl) FindByID(ctx context.Context, personID int) (web.PersonResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return web.PersonResponse{}, err
	}
	defer helper.RollBackCommit(tx)

	person, err := s.PersonRepository.FindByID(ctx, tx, personID)
	if err != nil {
		return web.PersonResponse{}, err
	}

	return helper.ToPersonRespone(*person), nil
}

func (s *PersonServiceImpl) FindAll(ctx context.Context) ([]web.PersonResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return []web.PersonResponse{}, err
	}
	defer helper.RollBackCommit(tx)

	person, err := s.PersonRepository.FindAll(ctx, tx)
	if err != nil {
		return []web.PersonResponse{}, err
	}

	return helper.ToAllPersonResponse(person), nil
}
