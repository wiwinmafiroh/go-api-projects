package service

import (
	"06-go-api-with-unittest/dto"
	"06-go-api-with-unittest/entity"
	"06-go-api-with-unittest/pkg/errs"
	"06-go-api-with-unittest/repository/product_repository"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProductService_GetProducts_Success(t *testing.T) {
	productRepository := product_repository.NewProductRepositoryMock()
	productService := NewProductService(productRepository)

	currentTime := time.Now()

	adminProducts := []*entity.Product{
		{
			Id:          1,
			Name:        "Laptop Pro X",
			Description: "Powerful laptop with advanced features for seamless performance.",
			Price:       15000000,
			ImageUrl:    "https://example.com/laptop_pro_x.jpg",
			UserId:      1,
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
		},
		{
			Id:          2,
			Name:        "Home Espresso Machine",
			Description: "Brew the perfect cup of coffee at home with our state-of-the-art Espresso Machine.",
			Price:       3000000,
			ImageUrl:    "https://example.com/espresso_machine.jpg",
			UserId:      2,
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
		},
	}

	userProducts := []*entity.Product{
		{
			Id:          2,
			Name:        "Home Espresso Machine",
			Description: "Brew the perfect cup of coffee at home with our state-of-the-art Espresso Machine.",
			Price:       3000000,
			ImageUrl:    "https://example.com/espresso_machine.jpg",
			UserId:      2,
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
		},
		{
			Id:          4,
			Name:        "Portable Bluetooth Speaker",
			Description: "Bring the party anywhere with our Portable Bluetooth Speaker.",
			Price:       500000,
			ImageUrl:    "https://example.com/portable_speaker.jpg",
			UserId:      2,
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
		},
	}

	product_repository.GetAllProducts = func() ([]*entity.Product, errs.ErrorMessage) {
		return adminProducts, nil
	}

	product_repository.GetUserProducts = func(userId uint) ([]*entity.Product, errs.ErrorMessage) {
		return userProducts, nil
	}

	tableTest := []struct {
		name        string
		userId      uint
		accessRole  string
		expectation []*entity.Product
	}{
		{
			name:        "Admin Products",
			userId:      1,
			accessRole:  "admin",
			expectation: adminProducts,
		},
		{
			name:        "User Products",
			userId:      2,
			accessRole:  "user",
			expectation: userProducts,
		},
	}

	for i, eachTest := range tableTest {
		t.Run(eachTest.name, func(t *testing.T) {
			response, err := productService.GetProducts(eachTest.userId, eachTest.accessRole)

			assert.Nil(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, "SUCCESS", response.Result)
			assert.Equal(t, http.StatusOK, response.StatusCode)
			assert.Equal(t, "Products retrieved successfully", response.Message)
			assert.Equal(t, len(eachTest.expectation), len(response.Data))
			assert.Equal(t, eachTest.expectation[i].Id, response.Data[i].Id)
			assert.Equal(t, eachTest.expectation[i].Name, response.Data[i].Name)
			assert.Equal(t, eachTest.expectation[i].Description, response.Data[i].Description)
			assert.Equal(t, eachTest.expectation[i].Price, response.Data[i].Price)
			assert.Equal(t, eachTest.expectation[i].ImageUrl, response.Data[i].ImageUrl)
			assert.Equal(t, eachTest.expectation[i].UserId, response.Data[i].UserId)
			assert.Equal(t, eachTest.expectation[i].CreatedAt, response.Data[i].CreatedAt)
			assert.Equal(t, eachTest.expectation[i].UpdatedAt, response.Data[i].UpdatedAt)
		})
	}
}

func TestProductService_GetProducts_NotFoundError(t *testing.T) {
	productRepository := product_repository.NewProductRepositoryMock()
	productService := NewProductService(productRepository)

	product_repository.GetAllProducts = func() ([]*entity.Product, errs.ErrorMessage) {
		return []*entity.Product{}, nil
	}

	product_repository.GetUserProducts = func(userId uint) ([]*entity.Product, errs.ErrorMessage) {
		return []*entity.Product{}, nil
	}

	tableTest := []struct {
		name        string
		userId      uint
		accessRole  string
		expectation []dto.RetrievedProductData
	}{
		{
			name:        "Admin Products Not Found",
			userId:      1,
			accessRole:  "admin",
			expectation: []dto.RetrievedProductData{},
		},
		{
			name:        "User Products Not Found",
			userId:      2,
			accessRole:  "user",
			expectation: []dto.RetrievedProductData{},
		},
	}

	for _, eachTest := range tableTest {
		t.Run(eachTest.name, func(t *testing.T) {
			response, err := productService.GetProducts(eachTest.userId, eachTest.accessRole)

			assert.Nil(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, "SUCCESS", response.Result)
			assert.Equal(t, http.StatusOK, response.StatusCode)
			assert.Equal(t, "Products retrieved successfully", response.Message)
			assert.Equal(t, 0, len(response.Data))
			assert.Equal(t, eachTest.expectation, response.Data)
		})
	}
}

func TestProductService_GetProductById_Success(t *testing.T) {
	productRepository := product_repository.NewProductRepositoryMock()
	productService := NewProductService(productRepository)

	currentTime := time.Now()

	product := entity.Product{
		Id:          1,
		Name:        "Laptop Pro X",
		Description: "Powerful laptop with advanced features for seamless performance.",
		Price:       15000000,
		ImageUrl:    "https://example.com/laptop_pro_x.jpg",
		UserId:      1,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}

	product_repository.GetProductById = func(productId uint) (*entity.Product, errs.ErrorMessage) {
		return &product, nil
	}

	response, err := productService.GetProductById(1)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Result)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Product retrieved successfully", response.Message)
	assert.Equal(t, uint(1), response.Data.Id)
	assert.Equal(t, "Laptop Pro X", response.Data.Name)
	assert.Equal(t, "Powerful laptop with advanced features for seamless performance.", response.Data.Description)
	assert.Equal(t, float64(15000000), response.Data.Price)
	assert.Equal(t, "https://example.com/laptop_pro_x.jpg", response.Data.ImageUrl)
	assert.Equal(t, uint(1), response.Data.UserId)
	assert.Equal(t, currentTime, response.Data.CreatedAt)
	assert.Equal(t, currentTime, response.Data.UpdatedAt)
}

func TestProductService_GetProductById_NotFoundError(t *testing.T) {
	productRepository := product_repository.NewProductRepositoryMock()
	productService := NewProductService(productRepository)

	product_repository.GetProductById = func(productId uint) (*entity.Product, errs.ErrorMessage) {
		return nil, errs.NewNotFoundError("Product not found")
	}

	response, err := productService.GetProductById(1)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "NOT_FOUND", err.Error())
	assert.Equal(t, http.StatusNotFound, err.StatusCode())
	assert.Equal(t, "Product not found", err.Message())
}
