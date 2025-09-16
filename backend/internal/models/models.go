package models

import (
	"gorm.io/gorm"
)

type UserType string

const (
	CustomerType UserType = "customer"
	DeliveryType UserType = "delivery"
	AdminType    UserType = "admin"
)

type User struct {
	gorm.Model
	Name     string   `json:"name"`
	Email    string   `json:"email" gorm:"unique"`
	Password string   `json:"-"` // Não será retornado no JSON
	Type     UserType `json:"type" gorm:"type:varchar(20)"`
	Vehicle  *string  `json:"vehicle,omitempty"` // Apenas para entregadores
}

type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"imageUrl"`
}

type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusPreparing  OrderStatus = "preparing"
	StatusReady      OrderStatus = "ready"
	StatusDelivering OrderStatus = "delivering"
	StatusDelivered  OrderStatus = "delivered"
)

type Order struct {
	gorm.Model
	CustomerID uint        `json:"customerId"`
	Customer   User        `json:"customer" gorm:"foreignKey:CustomerID"`
	DeliveryID *uint       `json:"deliveryId,omitempty"`
	Delivery   *User       `json:"delivery,omitempty" gorm:"foreignKey:DeliveryID"`
	Status     OrderStatus `json:"status" gorm:"type:varchar(20)"`
	Total      float64     `json:"total"`
	Address    string      `json:"address"`
	Items      []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"orderId"`
	ProductID uint    `json:"productId"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // Preço no momento da compra
}
