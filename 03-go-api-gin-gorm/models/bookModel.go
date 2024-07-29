package models

import "time"

type Book struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NameBook  string    `gorm:"type:varchar(191);not null" json:"name_book"`
	Author    string    `gorm:"type:varchar(191);not null" json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
