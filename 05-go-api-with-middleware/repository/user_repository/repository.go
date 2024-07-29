package user_repository

import (
	"05-go-api-with-middleware/entity"
	"05-go-api-with-middleware/pkg/errs"
)

type UserRepository interface {
	CreateUser(userEntity entity.User) errs.ErrorMessage
	GetUserByEmail(userEmail string) (*entity.User, errs.ErrorMessage)
}
