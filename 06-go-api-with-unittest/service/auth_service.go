package service

import (
	"06-go-api-with-unittest/entity"
	"06-go-api-with-unittest/pkg/errs"
	"06-go-api-with-unittest/pkg/helpers"
	"06-go-api-with-unittest/repository/product_repository"
	"06-go-api-with-unittest/repository/user_repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Authentication() gin.HandlerFunc
	AuthorizationRole() gin.HandlerFunc
	AuthorizationProduct() gin.HandlerFunc
}

type authService struct {
	userRepository    user_repository.UserRepository
	productRepository product_repository.ProductRepository
}

func NewAuthService(userRepository user_repository.UserRepository, productRepository product_repository.ProductRepository) AuthService {
	return &authService{
		userRepository:    userRepository,
		productRepository: productRepository,
	}
}

func (a *authService) Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var invalidTokenErr = errs.NewUnauthenticatedError("Invalid token")
		var bearerToken = ctx.GetHeader("Authorization")
		var userData entity.User

		err := userData.ValidateToken(bearerToken)
		if err != nil {
			ctx.AbortWithStatusJSON(invalidTokenErr.StatusCode(), invalidTokenErr)
			return
		}

		_, err = a.userRepository.GetUserByEmail(userData.Email)
		if err != nil {
			ctx.AbortWithStatusJSON(invalidTokenErr.StatusCode(), invalidTokenErr)
			return
		}

		ctx.Set("userData", userData)
		ctx.Next()
	}
}

func (a *authService) AuthorizationRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessRole := ctx.MustGet("userData").(entity.User).Role
		requestMethod := ctx.Request.Method
		unauthorizedErr := errs.NewUnauthorizedError("You don't have the required permissions for this operation")

		if accessRole != entity.AdminRole && accessRole != entity.UserRole {
			ctx.AbortWithStatusJSON(unauthorizedErr.StatusCode(), unauthorizedErr)
			return
		}

		if accessRole == entity.UserRole && (requestMethod == http.MethodPut || requestMethod == http.MethodDelete) {
			ctx.AbortWithStatusJSON(unauthorizedErr.StatusCode(), unauthorizedErr)
			return
		}

		ctx.Next()
	}
}

func (a *authService) AuthorizationProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userData := ctx.MustGet("userData").(entity.User)
		unauthorizedErr := errs.NewUnauthorizedError("You're not allowed to access this product")

		productId, err := helpers.GetParamId(ctx, "productId")
		if err != nil {
			ctx.AbortWithStatusJSON(err.StatusCode(), err)
			return
		}

		retrievedProduct, err := a.productRepository.GetProductById(productId)
		if err != nil {
			ctx.AbortWithStatusJSON(err.StatusCode(), err)
			return
		}

		if userData.Role != entity.AdminRole && userData.Id != retrievedProduct.UserId {
			ctx.AbortWithStatusJSON(unauthorizedErr.StatusCode(), unauthorizedErr)
			return
		}

		ctx.Next()
	}
}
