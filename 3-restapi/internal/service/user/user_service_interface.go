package userservice

import (
	"context"
	"restapi/internal/model/payload"
	"restapi/internal/model/web"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user payload.CreateUser) (*web.ToUser, error)
	UpdateUser(ctx context.Context, user payload.UpdateUser, userID int) (*web.ToUser, error)
	DeleteUser(ctx context.Context, userID int) error
	FindByID(ctx context.Context, userID int) (*web.ToUserDetail, error)
}
