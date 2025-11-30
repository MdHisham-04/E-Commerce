package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Name      string    `json:"name"`
	Password  string    `json:"-" gorm:"not null"`
	Role      string    `json:"role" gorm:"default:'buyer'"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"default:0"`
	SellerID    int       `json:"seller_id" gorm:"not null"`
	Seller      User      `json:"seller,omitempty" gorm:"foreignKey:SellerID"`
	CreatedAt   time.Time `json:"created_at"`
}

type CartItem struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id" gorm:"not null"`
	ProductID int       `json:"product_id" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID         int         `json:"id" gorm:"primaryKey"`
	UserID     int         `json:"user_id" gorm:"not null"`
	Total      float64     `json:"total" gorm:"not null"`
	Status     string      `json:"status" gorm:"default:'pending'"`
	User       User        `json:"user" gorm:"foreignKey:UserID"`
	OrderItems []OrderItem `json:"order_items"`
	CreatedAt  time.Time   `json:"created_at"`
}

type OrderItem struct {
	ID        int     `json:"id" gorm:"primaryKey"`
	OrderID   int     `json:"order_id" gorm:"not null"`
	ProductID int     `json:"product_id" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
	Status    string  `json:"status" gorm:"default:'pending'"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Order     Order   `json:"-" gorm:"foreignKey:OrderID"`
}
