package user_repository

import (
	"06-go-api-with-unittest/entity"
	"06-go-api-with-unittest/pkg/errs"
)

type UserRepository interface {
	CreateUser(userEntity entity.User) errs.ErrorMessage
	GetUserByEmail(userEmail string) (*entity.User, errs.ErrorMessage)
}
