package handlers

import (
	"cupcake-delivery/internal/models"
	"cupcake-delivery/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// GetNotifications busca as notificações do usuário logado
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	userID := userIDRaw.(uint)

	// Parâmetro opcional para limitar quantidade
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	notifications, err := h.notificationService.GetUserNotifications(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar notificações"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"count":         len(notifications),
	})
}

// GetUnreadCount retorna o número de notificações não lidas
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	userID := userIDRaw.(uint)

	count, err := h.notificationService.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar contagem"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

// MarkAsRead marca uma notificação específica como lida
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	userID := userIDRaw.(uint)

	notificationIDStr := c.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de notificação inválido"})
		return
	}

	err = h.notificationService.MarkAsRead(uint(notificationID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao marcar notificação como lida"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notificação marcada como lida"})
}

// MarkAllAsRead marca todas as notificações como lidas
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	userID := userIDRaw.(uint)

	err := h.notificationService.MarkAllAsRead(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao marcar notificações como lidas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todas as notificações foram marcadas como lidas"})
}

// CreateTestNotification cria uma notificação de teste (apenas para desenvolvimento)
func (h *NotificationHandler) CreateTestNotification(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	userID := userIDRaw.(uint)

	var req struct {
		Title   string `json:"title" binding:"required"`
		Message string `json:"message" binding:"required"`
		Type    string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Type == "" {
		req.Type = "test"
	}

	notificationData := models.CreateNotificationData{
		UserID:  userID,
		Type:    req.Type,
		Title:   req.Title,
		Message: req.Message,
	}

	notification, err := h.notificationService.CreateNotification(notificationData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar notificação"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Notificação criada com sucesso",
		"notification": notification,
	})
}
