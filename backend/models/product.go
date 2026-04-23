package models

import "time"

type Product struct {
	ID        uint      `json:"id" gorm:"primaryKey" `
	Name      string    `json:"name" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	Stock     *int      `json:"stock" gorm:"not null" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var Products []Product
