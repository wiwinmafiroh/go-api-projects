package handler

import (
	"05-go-api-with-middleware/dto"
	"05-go-api-with-middleware/pkg/errs"
	"05-go-api-with-middleware/pkg/helpers"
	"05-go-api-with-middleware/service"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) userHandler {
	return userHandler{
		userService: userService,
	}
}

func (u *userHandler) UserRegister(ctx *gin.Context) {
	err := helpers.CheckContentType(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	var userRequest dto.UserRegisterRequest

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		errBindJSON := errs.NewUnprocessableEntityError("Invalid request body")

		ctx.AbortWithStatusJSON(errBindJSON.StatusCode(), errBindJSON)
		return
	}

	result, err := u.userService.UserRegister(userRequest)
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	ctx.JSON(result.StatusCode, result)
}

func (u *userHandler) UserLogin(ctx *gin.Context) {
	err := helpers.CheckContentType(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	var userRequest dto.UserLoginRequest

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		errBindJSON := errs.NewUnprocessableEntityError("Invalid request body")

		ctx.AbortWithStatusJSON(errBindJSON.StatusCode(), errBindJSON)
		return
	}

	result, err := u.userService.UserLogin(userRequest)
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	ctx.JSON(result.StatusCode, result)
}
