package db

import (
	"context"

	"formation/model"
)

type UserStore interface {
	CreateUser(ctx context.Context, u *model.User) error
	GetUserByID(ctx context.Context, uuid string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*model.User, error)
	UpdateUser(ctx context.Context, uuid string, u *model.User) error
	DeleteUser(ctx context.Context, uuid string) error
	Close() error
}
