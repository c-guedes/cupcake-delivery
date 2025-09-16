package handlers

import (
	"fmt"
	"net/http"

	"cupcake-delivery/internal/models"
	"cupcake-delivery/internal/utils"
	"cupcake-delivery/internal/validators"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// Auth handlers
func (h *Handler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, utils.ErrorTypeValidation, "Dados JSON inválidos")
		return
	}

	// Validações
	var validationErrors []utils.ValidationError

	if err := validators.ValidateEmail(user.Email); err != nil {
		validationErrors = append(validationErrors, *err)
	}

	if err := validators.ValidatePassword(user.Password); err != nil {
		validationErrors = append(validationErrors, *err)
	}

	if err := validators.ValidateName(user.Name); err != nil {
		validationErrors = append(validationErrors, *err)
	}

	if err := validators.ValidateUserType(string(user.Type)); err != nil {
		validationErrors = append(validationErrors, *err)
	}

	if len(validationErrors) > 0 {
		utils.RespondWithValidationError(c, validationErrors)
		return
	}

	// Verificar se email já existe
	var existingUser models.User
	if err := h.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		utils.RespondWithError(c, http.StatusConflict, utils.ErrorTypeConflict, utils.MessageEmailExists)
		return
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, utils.ErrorTypeInternal, utils.MessageInternalError)
		return
	}
	user.Password = string(hashedPassword)

	// Criar usuário
	if err := h.db.Create(&user).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, utils.ErrorTypeInternal, "Erro ao criar usuário")
		return
	}

	// Remover senha da resposta
	user.Password = ""
	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuário criado com sucesso",
		"user":    user,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, utils.ErrorTypeValidation, "Dados JSON inválidos")
		return
	}

	// Validações básicas
	var validationErrors []utils.ValidationError

	if err := validators.ValidateEmail(credentials.Email); err != nil {
		validationErrors = append(validationErrors, *err)
	}

	if credentials.Password == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field:   "password",
			Message: "Senha é obrigatória",
		})
	}

	if len(validationErrors) > 0 {
		utils.RespondWithValidationError(c, validationErrors)
		return
	}

	// Buscar usuário
	var user models.User
	if err := h.db.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, utils.ErrorTypeAuthentication, utils.MessageInvalidCredentials)
		return
	}

	// Verificar senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, utils.ErrorTypeAuthentication, utils.MessageInvalidCredentials)
		return
	}

	// TODO: Gerar JWT token
	c.JSON(http.StatusOK, gin.H{
		"message": "Login realizado com sucesso",
		"token":   "seu_token_jwt",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"type":  user.Type,
		},
	})
}

// Product handlers
func (h *Handler) ListProducts(c *gin.Context) {
	var products []models.Product
	if err := h.db.Find(&products).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, utils.ErrorTypeInternal, "Erro ao buscar produtos")
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, utils.ErrorTypeValidation, "Dados JSON inválidos")
		return
	}

	// Validações
	var validationErrors []utils.ValidationError

	if err := validators.ValidateProductName(product.Name); err != nil {
		validationErrors = append(validationErrors, *err)
	}

	if err := validators.ValidateProductDescription(product.Description); err != nil {
		validationErrors = append(validationErrors, *err)
	}

	if err := validators.ValidateProductPrice(product.Price); err != nil {
		validationErrors = append(validationErrors, *err)
	}

	if err := validators.ValidateImageUrl(product.ImageURL); err != nil {
		validationErrors = append(validationErrors, *err)
	}

	if len(validationErrors) > 0 {
		utils.RespondWithValidationError(c, validationErrors)
		return
	}

	// Verificar se produto já existe
	var existingProduct models.Product
	if err := h.db.Where("name = ?", product.Name).First(&existingProduct).Error; err == nil {
		utils.RespondWithError(c, http.StatusConflict, utils.ErrorTypeConflict, "Produto com este nome já existe")
		return
	}

	// Criar produto
	if err := h.db.Create(&product).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, utils.ErrorTypeInternal, "Erro ao criar produto")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Produto criado com sucesso",
		"product": product,
	})
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := h.db.First(&product, id).Error; err != nil {
		utils.RespondWithError(c, http.StatusNotFound, utils.ErrorTypeNotFound, "Produto não encontrado")
		return
	}

	var updates models.Product
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, utils.ErrorTypeValidation, "Dados JSON inválidos")
		return
	}

	// Validações
	var validationErrors []utils.ValidationError

	if updates.Name != "" {
		if err := validators.ValidateProductName(updates.Name); err != nil {
			validationErrors = append(validationErrors, *err)
		}
	}

	if updates.Description != "" {
		if err := validators.ValidateProductDescription(updates.Description); err != nil {
			validationErrors = append(validationErrors, *err)
		}
	}

	if updates.Price > 0 {
		if err := validators.ValidateProductPrice(updates.Price); err != nil {
			validationErrors = append(validationErrors, *err)
		}
	}

	if updates.ImageURL != "" {
		if err := validators.ValidateImageUrl(updates.ImageURL); err != nil {
			validationErrors = append(validationErrors, *err)
		}
	}

	if len(validationErrors) > 0 {
		utils.RespondWithValidationError(c, validationErrors)
		return
	}

	// Atualizar produto
	if err := h.db.Model(&product).Updates(updates).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, utils.ErrorTypeInternal, "Erro ao atualizar produto")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produto atualizado com sucesso",
		"product": product,
	})
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := h.db.First(&product, id).Error; err != nil {
		utils.RespondWithError(c, http.StatusNotFound, utils.ErrorTypeNotFound, "Produto não encontrado")
		return
	}

	if err := h.db.Delete(&product).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, utils.ErrorTypeInternal, "Erro ao deletar produto")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produto deletado com sucesso",
	})
}

// Order handlers
func (h *Handler) CreateOrder(c *gin.Context) {
	var orderRequest struct {
		Items []struct {
			ProductID uint `json:"product_id" binding:"required"`
			Quantity  int  `json:"quantity" binding:"required,min=1"`
		} `json:"items" binding:"required,min=1"`
		Address       string `json:"address"`
		Phone         string `json:"phone"`
		PaymentMethod string `json:"paymentMethod"`
		Notes         string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, utils.ErrorTypeValidation, "Dados JSON inválidos")
		return
	}

	// Validações
	var validationErrors []utils.ValidationError

	if len(orderRequest.Items) == 0 {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field:   "items",
			Message: "Pelo menos um item é obrigatório",
		})
	}

	for i, item := range orderRequest.Items {
		if err := validators.ValidateQuantity(item.Quantity); err != nil {
			err.Field = fmt.Sprintf("items[%d].quantity", i)
			validationErrors = append(validationErrors, *err)
		}
	}

	if len(validationErrors) > 0 {
		utils.RespondWithValidationError(c, validationErrors)
		return
	}

	// Pegar ID do usuário do token JWT
	userID, exists := c.Get("user_id")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, utils.ErrorTypeAuthentication, utils.MessageUnauthorized)
		return
	}

	// Iniciar transação
	tx := h.db.Begin()

	// Criar pedido
	var order models.Order
	order = models.Order{
		CustomerID: userID.(uint),
		Status:     models.StatusPending,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		utils.RespondWithError(c, http.StatusInternalServerError, utils.ErrorTypeInternal, "Erro ao criar pedido")
		return
	}

	// Commit da transação
	if err := tx.Commit().Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, utils.ErrorTypeInternal, "Erro ao finalizar pedido")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pedido criado com sucesso",
		"order":   order,
	})
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	var status struct {
		Status models.OrderStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Atualizar status do pedido

	c.JSON(http.StatusOK, gin.H{"message": "Status atualizado"})
}
