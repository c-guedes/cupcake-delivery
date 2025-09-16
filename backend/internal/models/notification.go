package models

import (
	"time"
)

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	OrderID   *uint     `json:"order_id,omitempty"`
	Type      string    `json:"type" gorm:"not null"`
	Title     string    `json:"title" gorm:"not null"`
	Message   string    `json:"message" gorm:"not null"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relações
	User  User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Order *Order `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// Tipos de notificações
const (
	NotificationTypeOrderCreated    = "order_created"
	NotificationTypeOrderConfirmed  = "order_confirmed"
	NotificationTypeOrderPreparing  = "order_preparing"
	NotificationTypeOrderReady      = "order_ready"
	NotificationTypeOrderDelivering = "order_delivering"
	NotificationTypeOrderDelivered  = "order_delivered"
	NotificationTypeOrderCancelled  = "order_cancelled"
)

// CreateNotificationData estrutura para criação de notificações
type CreateNotificationData struct {
	UserID  uint   `json:"user_id"`
	OrderID *uint  `json:"order_id,omitempty"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
}
