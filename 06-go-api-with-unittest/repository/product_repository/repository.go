package product_repository

import (
	"06-go-api-with-unittest/entity"
	"06-go-api-with-unittest/pkg/errs"
)

type ProductRepository interface {
	CreateProduct(productEntity entity.Product) (*entity.Product, errs.ErrorMessage)
	GetAllProducts() ([]*entity.Product, errs.ErrorMessage)
	GetUserProducts(userId uint) ([]*entity.Product, errs.ErrorMessage)
	GetProductById(productId uint) (*entity.Product, errs.ErrorMessage)
	UpdateProductById(productEntity entity.Product) (*entity.Product, errs.ErrorMessage)
	DeleteProductById(productId uint) errs.ErrorMessage
}
