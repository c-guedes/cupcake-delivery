package services

import (
	"cupcake-delivery/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

// CreateNotification cria uma nova notificação
func (s *NotificationService) CreateNotification(data models.CreateNotificationData) (*models.Notification, error) {
	notification := &models.Notification{
		UserID:  data.UserID,
		OrderID: data.OrderID,
		Type:    data.Type,
		Title:   data.Title,
		Message: data.Message,
		IsRead:  false,
	}

	if err := s.db.Create(notification).Error; err != nil {
		return nil, err
	}

	return notification, nil
}

// GetUserNotifications busca notificações de um usuário
func (s *NotificationService) GetUserNotifications(userID uint, limit int) ([]models.Notification, error) {
	var notifications []models.Notification

	query := s.db.Where("user_id = ?", userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

// MarkAsRead marca uma notificação como lida
func (s *NotificationService) MarkAsRead(notificationID uint, userID uint) error {
	return s.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}

// MarkAllAsRead marca todas as notificações de um usuário como lidas
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}

// GetUnreadCount retorna o número de notificações não lidas
func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}

// NotifyOrderStatusChange notifica sobre mudança de status do pedido
func (s *NotificationService) NotifyOrderStatusChange(order *models.Order, newStatus string) error {
	// Definir mensagens baseadas no status
	statusMessages := map[string]struct {
		customerTitle   string
		customerMessage string
		deliveryTitle   string
		deliveryMessage string
		adminTitle      string
		adminMessage    string
	}{
		"confirmed": {
			customerTitle:   "Pedido Confirmado!",
			customerMessage: fmt.Sprintf("Seu pedido #%d foi confirmado e está sendo preparado.", order.ID),
			adminTitle:      "Novo Pedido Confirmado",
			adminMessage:    fmt.Sprintf("Pedido #%d foi confirmado e deve ser preparado.", order.ID),
		},
		"preparing": {
			customerTitle:   "Preparando seu Pedido",
			customerMessage: fmt.Sprintf("Seu pedido #%d está sendo preparado com carinho!", order.ID),
		},
		"ready": {
			customerTitle:   "Pedido Pronto!",
			customerMessage: fmt.Sprintf("Seu pedido #%d está pronto e será entregue em breve.", order.ID),
			deliveryTitle:   "Novo Pedido para Entrega",
			deliveryMessage: fmt.Sprintf("Pedido #%d está pronto para entrega.", order.ID),
		},
		"delivering": {
			customerTitle:   "Pedido a Caminho!",
			customerMessage: fmt.Sprintf("Seu pedido #%d saiu para entrega. Em breve estará aí!", order.ID),
		},
		"delivered": {
			customerTitle:   "Pedido Entregue!",
			customerMessage: fmt.Sprintf("Seu pedido #%d foi entregue com sucesso. Obrigado!", order.ID),
			adminTitle:      "Pedido Entregue",
			adminMessage:    fmt.Sprintf("Pedido #%d foi entregue com sucesso.", order.ID),
		},
		"cancelled": {
			customerTitle:   "Pedido Cancelado",
			customerMessage: fmt.Sprintf("Seu pedido #%d foi cancelado. Entre em contato conosco para mais informações.", order.ID),
			adminTitle:      "Pedido Cancelado",
			adminMessage:    fmt.Sprintf("Pedido #%d foi cancelado.", order.ID),
		},
	}

	messages, exists := statusMessages[newStatus]
	if !exists {
		return fmt.Errorf("status de pedido não reconhecido: %s", newStatus)
	}

	// Notificar o cliente (sempre)
	if messages.customerTitle != "" {
		customerNotification := models.CreateNotificationData{
			UserID:  order.CustomerID,
			OrderID: &order.ID,
			Type:    fmt.Sprintf("order_%s", newStatus),
			Title:   messages.customerTitle,
			Message: messages.customerMessage,
		}
		if _, err := s.CreateNotification(customerNotification); err != nil {
			return err
		}
	}

	// Buscar usuários admin para notificar
	if messages.adminTitle != "" {
		var adminUsers []models.User
		if err := s.db.Where("type = ?", "admin").Find(&adminUsers).Error; err == nil {
			for _, admin := range adminUsers {
				adminNotification := models.CreateNotificationData{
					UserID:  admin.ID,
					OrderID: &order.ID,
					Type:    fmt.Sprintf("order_%s", newStatus),
					Title:   messages.adminTitle,
					Message: messages.adminMessage,
				}
				s.CreateNotification(adminNotification)
			}
		}
	}

	// Buscar usuários entregadores para notificar (apenas para pedidos prontos)
	if messages.deliveryTitle != "" {
		var deliveryUsers []models.User
		if err := s.db.Where("type = ?", "delivery").Find(&deliveryUsers).Error; err == nil {
			for _, delivery := range deliveryUsers {
				deliveryNotification := models.CreateNotificationData{
					UserID:  delivery.ID,
					OrderID: &order.ID,
					Type:    fmt.Sprintf("order_%s", newStatus),
					Title:   messages.deliveryTitle,
					Message: messages.deliveryMessage,
				}
				s.CreateNotification(deliveryNotification)
			}
		}
	}

	return nil
}
