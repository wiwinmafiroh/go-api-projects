package entity

import (
	"05-go-api-with-middleware/dto"
	"time"
)

type Product struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImageUrl    string    `json:"imageUrl"`
	UserId      uint      `json:"userId"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

func (p *Product) ConvertProductEntityToDto() dto.RetrievedProductData {
	return dto.RetrievedProductData{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		ImageUrl:    p.ImageUrl,
		UserId:      p.UserId,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
