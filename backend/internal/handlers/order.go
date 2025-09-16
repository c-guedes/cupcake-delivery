package handlers

import (
	"net/http"

	"cupcake-delivery/internal/models"
	"cupcake-delivery/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	db                  *gorm.DB
	notificationService *services.NotificationService
}

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required,min=1"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

func NewOrderHandler(db *gorm.DB, notificationService *services.NotificationService) *OrderHandler {
	return &OrderHandler{
		db:                  db,
		notificationService: notificationService,
	}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pegar ID do usuário do token JWT
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	// Iniciar transação
	tx := h.db.Begin()

	// Criar pedido
	order := models.Order{
		CustomerID: userID.(uint),
		Status:     models.StatusPending,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar pedido"})
		return
	}

	// Adicionar itens ao pedido
	var totalPrice float64 = 0
	for _, item := range req.Items {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Produto não encontrado"})
			return
		}

		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar item do pedido"})
			return
		}

		totalPrice += product.Price * float64(item.Quantity)
	}

	// Atualizar preço total do pedido
	order.Total = totalPrice
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar preço total"})
		return
	}

	// Commit da transação
	tx.Commit()

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) List(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	role, _ := c.Get("type")

	var orders []models.Order

	// Filtrar pedidos baseado no role
	switch role {
	case string(models.CustomerType):
		if err := h.db.Where("customer_id = ?", userID.(uint)).Order("created_at DESC").Find(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar pedidos do cliente", "details": err.Error()})
			return
		}
	case string(models.DeliveryType):
		if err := h.db.Where("delivery_id = ? OR (delivery_id IS NULL AND status IN ?)", userID.(uint), []models.OrderStatus{models.StatusReady}).Order("created_at DESC").Find(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar pedidos do entregador", "details": err.Error()})
			return
		}
	case string(models.AdminType):
		if err := h.db.Order("created_at DESC").Find(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar todos os pedidos", "details": err.Error()})
			return
		}
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Acesso não autorizado"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	orderID := c.Param("id")
	status := c.PostForm("status")

	userID, _ := c.Get("user_id")
	role, _ := c.Get("type")

	var order models.Order
	if err := h.db.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado"})
		return
	}

	// Verificar permissões
	switch role {
	case string(models.DeliveryType):
		if order.DeliveryID != nil && *order.DeliveryID != userID.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acesso não autorizado"})
			return
		}
		// Entregador só pode atualizar para "delivering" ou "delivered"
		if status != string(models.StatusDelivering) && status != string(models.StatusDelivered) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Status inválido para entregador"})
			return
		}
		if status == string(models.StatusDelivering) {
			deliveryID := userID.(uint)
			order.DeliveryID = &deliveryID
		}
	case string(models.AdminType):
		// Admin pode atualizar para qualquer status
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Acesso não autorizado"})
		return
	}

	order.Status = models.OrderStatus(status)
	if err := h.db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar status"})
		return
	}

	// Disparar notificações sobre mudança de status
	if h.notificationService != nil {
		if err := h.notificationService.NotifyOrderStatusChange(&order, status); err != nil {
			// Log do erro, mas não falha a operação
			// TODO: implementar logger adequado
		}
	}

	c.JSON(http.StatusOK, order)
}
