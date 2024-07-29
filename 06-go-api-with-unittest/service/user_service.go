package service

import (
	"06-go-api-with-unittest/dto"
	"06-go-api-with-unittest/entity"
	"06-go-api-with-unittest/pkg/errs"
	"06-go-api-with-unittest/pkg/helpers"
	"06-go-api-with-unittest/repository/user_repository"
	"net/http"
)

type UserService interface {
	UserRegister(userRequest dto.UserRegisterRequest) (*dto.UserRegisterResponse, errs.ErrorMessage)
	UserLogin(userRequest dto.UserLoginRequest) (*dto.UserLoginResponse, errs.ErrorMessage)
}

type userService struct {
	userRepository user_repository.UserRepository
}

func NewUserService(userRepository user_repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (u *userService) UserRegister(userRequest dto.UserRegisterRequest) (*dto.UserRegisterResponse, errs.ErrorMessage) {
	err := helpers.ValidateStruct(userRequest)
	if err != nil {
		return nil, err
	}

	userEntity := entity.User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: userRequest.Password,
	}

	err = userEntity.HashPassword()
	if err != nil {
		return nil, err
	}

	err = u.userRepository.CreateUser(userEntity)
	if err != nil {
		return nil, err
	}

	response := dto.UserRegisterResponse{
		Result:     "SUCCESS",
		StatusCode: http.StatusCreated,
		Message:    "User registered successfully",
	}

	return &response, nil
}

func (u *userService) UserLogin(userRequest dto.UserLoginRequest) (*dto.UserLoginResponse, errs.ErrorMessage) {
	err := helpers.ValidateStruct(userRequest)
	if err != nil {
		return nil, err
	}

	retrievedUser, err := u.userRepository.GetUserByEmail(userRequest.Email)
	if err != nil {
		if err.StatusCode() == http.StatusNotFound {
			return nil, errs.NewUnauthenticatedError("Invalid email or password")
		}

		return nil, err
	}

	isValidPassword := retrievedUser.ComparePassword(userRequest.Password)
	if !isValidPassword {
		return nil, errs.NewUnauthenticatedError("Invalid email or password")
	}

	response := dto.UserLoginResponse{
		Result:     "SUCCESS",
		StatusCode: http.StatusOK,
		Message:    "User logged in successfully",
		Data: dto.TokenData{
			Token: retrievedUser.GenerateToken(),
		},
	}

	return &response, nil
}
