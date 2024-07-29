package product_repository

import (
	"06-go-api-with-unittest/entity"
	"06-go-api-with-unittest/pkg/errs"
)

var (
	CreateProduct     func(productEntity entity.Product) (*entity.Product, errs.ErrorMessage)
	GetAllProducts    func() ([]*entity.Product, errs.ErrorMessage)
	GetUserProducts   func(userId uint) ([]*entity.Product, errs.ErrorMessage)
	GetProductById    func(productId uint) (*entity.Product, errs.ErrorMessage)
	UpdateProductById func(productEntity entity.Product) (*entity.Product, errs.ErrorMessage)
	DeleteProductById func(productId uint) errs.ErrorMessage
)

type productRepositoryMock struct{}

func NewProductRepositoryMock() ProductRepository {
	return &productRepositoryMock{}
}

func (p *productRepositoryMock) CreateProduct(productEntity entity.Product) (*entity.Product, errs.ErrorMessage) {
	return CreateProduct(productEntity)
}

func (p *productRepositoryMock) GetAllProducts() ([]*entity.Product, errs.ErrorMessage) {
	return GetAllProducts()
}

func (p *productRepositoryMock) GetUserProducts(userId uint) ([]*entity.Product, errs.ErrorMessage) {
	return GetUserProducts(userId)
}

func (p *productRepositoryMock) GetProductById(productId uint) (*entity.Product, errs.ErrorMessage) {
	return GetProductById(productId)
}

func (p *productRepositoryMock) UpdateProductById(productEntity entity.Product) (*entity.Product, errs.ErrorMessage) {
	return UpdateProductById(productEntity)
}

func (p *productRepositoryMock) DeleteProductById(productId uint) errs.ErrorMessage {
	return DeleteProductById(productId)
}
