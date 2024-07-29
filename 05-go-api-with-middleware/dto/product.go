package dto

import "time"

type ProductRequest struct {
	Name        string  `json:"name" valid:"required~Name cannot be empty"`
	Description string  `json:"description"`
	Price       float64 `json:"price" valid:"required~Price cannot be empty"`
	ImageUrl    string  `json:"imageUrl" valid:"required~Image Url cannot be empty"`
}

type ProductCreatedResponse struct {
	Result     string             `json:"result"`
	StatusCode int                `json:"statusCode"`
	Message    string             `json:"message"`
	Data       CreatedProductData `json:"data"`
}

type ProductsRetrievedResponse struct {
	Result     string                 `json:"result"`
	StatusCode int                    `json:"statusCode"`
	Message    string                 `json:"message"`
	Data       []RetrievedProductData `json:"data"`
}

type ProductRetrievedResponse struct {
	Result     string               `json:"result"`
	StatusCode int                  `json:"statusCode"`
	Message    string               `json:"message"`
	Data       RetrievedProductData `json:"data"`
}

type ProductUpdatedResponse struct {
	Result     string             `json:"result"`
	StatusCode int                `json:"statusCode"`
	Message    string             `json:"message"`
	Data       UpdatedProductData `json:"data"`
}

type ProductDeletedResponse struct {
	Result     string    `json:"result"`
	StatusCode int       `json:"statusCode"`
	Message    string    `json:"message"`
	DeletedAt  time.Time `json:"deletedAt,omitempty"`
}

type CreatedProductData struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImageUrl    string    `json:"imageUrl"`
	UserId      uint      `json:"userId"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

type UpdatedProductData struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImageUrl    string    `json:"imageUrl"`
	UserId      uint      `json:"userId"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

type RetrievedProductData struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImageUrl    string    `json:"imageUrl"`
	UserId      uint      `json:"userId"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}
