package repository

import "basic-rest-api/model/domain"

type UserRepoInterface interface {
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByID(id int) (*domain.User, error)
	CreateUser(user domain.User) error
}
