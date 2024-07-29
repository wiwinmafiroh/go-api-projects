package product_repository

import (
	"05-go-api-with-middleware/entity"
	"05-go-api-with-middleware/pkg/errs"
)

type ProductRepository interface {
	CreateProduct(productEntity entity.Product) (*entity.Product, errs.ErrorMessage)
	GetAllProducts() ([]*entity.Product, errs.ErrorMessage)
	GetUserProducts(userId uint) ([]*entity.Product, errs.ErrorMessage)
	GetProductById(productId uint) (*entity.Product, errs.ErrorMessage)
	UpdateProductById(productEntity entity.Product) (*entity.Product, errs.ErrorMessage)
	DeleteProductById(productId uint) errs.ErrorMessage
}
