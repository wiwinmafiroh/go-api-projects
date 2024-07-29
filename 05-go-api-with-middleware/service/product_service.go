package service

import (
	"05-go-api-with-middleware/dto"
	"05-go-api-with-middleware/entity"
	"05-go-api-with-middleware/pkg/errs"
	"05-go-api-with-middleware/pkg/helpers"
	"05-go-api-with-middleware/repository/product_repository"
	"net/http"
	"time"
)

type ProductService interface {
	CreateProduct(userId uint, productRequest dto.ProductRequest) (*dto.ProductCreatedResponse, errs.ErrorMessage)
	GetProducts(userId uint, accessRole string) (*dto.ProductsRetrievedResponse, errs.ErrorMessage)
	GetProductById(productId uint) (*dto.ProductRetrievedResponse, errs.ErrorMessage)
	UpdateProductById(productId uint, productRequest dto.ProductRequest) (*dto.ProductUpdatedResponse, errs.ErrorMessage)
	DeleteProductById(productId uint) (*dto.ProductDeletedResponse, errs.ErrorMessage)
}

type productService struct {
	productRepository product_repository.ProductRepository
}

func NewProductService(productRepository product_repository.ProductRepository) ProductService {
	return &productService{
		productRepository: productRepository,
	}
}

func (p *productService) CreateProduct(userId uint, productRequest dto.ProductRequest) (*dto.ProductCreatedResponse, errs.ErrorMessage) {
	err := helpers.ValidateStruct(productRequest)
	if err != nil {
		return nil, err
	}

	productEntity := entity.Product{
		Name:        productRequest.Name,
		Description: productRequest.Description,
		Price:       productRequest.Price,
		ImageUrl:    productRequest.ImageUrl,
		UserId:      userId,
	}

	createdProduct, err := p.productRepository.CreateProduct(productEntity)
	if err != nil {
		return nil, err
	}

	response := dto.ProductCreatedResponse{
		Result:     "SUCCESS",
		StatusCode: http.StatusOK,
		Message:    "Product created successfully",
		Data: dto.CreatedProductData{
			Id:          createdProduct.Id,
			Name:        createdProduct.Name,
			Description: createdProduct.Description,
			Price:       createdProduct.Price,
			ImageUrl:    createdProduct.ImageUrl,
			UserId:      createdProduct.UserId,
			CreatedAt:   createdProduct.CreatedAt,
		},
	}

	return &response, nil
}

func (p *productService) GetProducts(userId uint, accessRole string) (*dto.ProductsRetrievedResponse, errs.ErrorMessage) {
	var retrievedProducts []*entity.Product
	var err errs.ErrorMessage

	switch accessRole {
	case string(entity.AdminRole):
		retrievedProducts, err = p.productRepository.GetAllProducts()
	case string(entity.UserRole):
		retrievedProducts, err = p.productRepository.GetUserProducts(userId)
	}

	if err != nil {
		return nil, err
	}

	var retrievedProductsDto []dto.RetrievedProductData

	for _, product := range retrievedProducts {
		retrievedProductsDto = append(retrievedProductsDto, product.ConvertProductEntityToDto())
	}

	if len(retrievedProductsDto) == 0 {
		retrievedProductsDto = []dto.RetrievedProductData{}
	}

	response := dto.ProductsRetrievedResponse{
		Result:     "SUCCESS",
		StatusCode: http.StatusOK,
		Message:    "Products retrieved successfully",
		Data:       retrievedProductsDto,
	}

	return &response, nil
}

func (p *productService) GetProductById(productId uint) (*dto.ProductRetrievedResponse, errs.ErrorMessage) {
	retrievedProduct, err := p.productRepository.GetProductById(productId)
	if err != nil {
		return nil, err
	}

	response := dto.ProductRetrievedResponse{
		Result:     "SUCCESS",
		StatusCode: http.StatusOK,
		Message:    "Product retrieved successfully",
		Data: dto.RetrievedProductData{
			Id:          retrievedProduct.Id,
			Name:        retrievedProduct.Name,
			Description: retrievedProduct.Description,
			Price:       retrievedProduct.Price,
			ImageUrl:    retrievedProduct.ImageUrl,
			UserId:      retrievedProduct.UserId,
			CreatedAt:   retrievedProduct.CreatedAt,
			UpdatedAt:   retrievedProduct.UpdatedAt,
		},
	}

	return &response, nil
}

func (p *productService) UpdateProductById(productId uint, productRequest dto.ProductRequest) (*dto.ProductUpdatedResponse, errs.ErrorMessage) {
	err := helpers.ValidateStruct(productRequest)
	if err != nil {
		return nil, err
	}

	productEntity := entity.Product{
		Id:          productId,
		Name:        productRequest.Name,
		Description: productRequest.Description,
		Price:       productRequest.Price,
		ImageUrl:    productRequest.ImageUrl,
	}

	updatedProduct, err := p.productRepository.UpdateProductById(productEntity)
	if err != nil {
		return nil, err
	}

	response := dto.ProductUpdatedResponse{
		Result:     "SUCCESS",
		StatusCode: http.StatusOK,
		Message:    "Product updated successfully",
		Data: dto.UpdatedProductData{
			Id:          updatedProduct.Id,
			Name:        updatedProduct.Name,
			Description: updatedProduct.Description,
			Price:       updatedProduct.Price,
			ImageUrl:    updatedProduct.ImageUrl,
			UserId:      updatedProduct.UserId,
			UpdatedAt:   updatedProduct.UpdatedAt,
		},
	}

	return &response, nil
}

func (p *productService) DeleteProductById(productId uint) (*dto.ProductDeletedResponse, errs.ErrorMessage) {
	_, err := p.productRepository.GetProductById(productId)
	if err != nil {
		return nil, err
	}

	err = p.productRepository.DeleteProductById(productId)
	if err != nil {
		return nil, err
	}

	response := dto.ProductDeletedResponse{
		Result:     "SUCCESS",
		StatusCode: http.StatusOK,
		Message:    "Product deleted successfully",
		DeletedAt:  time.Now(),
	}

	return &response, nil
}
