package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cupcake-delivery/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test user struct with password included in JSON
type TestUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"` // Include password for tests
	Type     string `json:"type"`
}

// Mock handler for testing without database
type MockHandler struct{}

func (h *MockHandler) Register(c *gin.Context) {
	var user TestUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Simulate validation with real logic
	if user.Email == "invalid-email" || !strings.Contains(user.Email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation_error",
			"validations": []map[string]string{
				{"field": "email", "message": "Formato de email inválido"},
			},
		})
		return
	}

	if len(user.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation_error",
			"validations": []map[string]string{
				{"field": "password", "message": "Senha deve ter pelo menos 6 caracteres"},
			},
		})
		return
	}

	if user.Type != "customer" && user.Type != "admin" && user.Type != "delivery" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation_error",
			"validations": []map[string]string{
				{"field": "type", "message": "Tipo de usuário deve ser 'customer', 'delivery' ou 'admin'"},
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuário criado com sucesso"})
}

func (h *MockHandler) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Simulate authentication
	if credentials.Email == "test@example.com" && credentials.Password == "password" {
		c.JSON(http.StatusOK, gin.H{"message": "Login realizado com sucesso", "token": "mock_token"})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication_error", "message": "Email ou senha inválidos"})
}

func (h *MockHandler) ListProducts(c *gin.Context) {
	mockProducts := []models.Product{
		{Name: "Cupcake 1", Description: "Delicious cupcake", Price: 10.0, ImageURL: "http://example.com/image1.jpg"},
		{Name: "Cupcake 2", Description: "Another delicious cupcake", Price: 15.0, ImageURL: "http://example.com/image2.jpg"},
	}
	c.JSON(http.StatusOK, mockProducts)
}

func (h *MockHandler) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Simulate validation
	if product.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation_error",
			"validations": []map[string]string{
				{"field": "price", "message": "Preço deve ser maior que zero"},
			},
		})
		return
	}

	if product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation_error",
			"validations": []map[string]string{
				{"field": "name", "message": "Nome do produto é obrigatório"},
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Produto criado com sucesso"})
}

// Setup test router with mock handlers
func setupMockRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	handler := &MockHandler{}

	// Auth routes
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	// Product routes
	router.GET("/products", handler.ListProducts)
	router.POST("/products", handler.CreateProduct)

	return router
}

// Tests for user registration
func TestMockRegister(t *testing.T) {
	router := setupMockRouter()

	t.Run("valid registration", func(t *testing.T) {
		user := TestUser{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
			Type:     "customer",
		}

		jsonData, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Usuário criado com sucesso")
	})

	t.Run("invalid email format", func(t *testing.T) {
		user := TestUser{
			Name:     "Test User",
			Email:    "invalid-email",
			Password: "password123",
			Type:     "customer",
		}

		jsonData, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "validation_error")
		assert.Contains(t, w.Body.String(), "email")
	})

	t.Run("short password", func(t *testing.T) {
		user := TestUser{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "123",
			Type:     "customer",
		}

		jsonData, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "validation_error")
		assert.Contains(t, w.Body.String(), "password")
	})
}

// Tests for user login
func TestMockLogin(t *testing.T) {
	router := setupMockRouter()

	t.Run("valid login", func(t *testing.T) {
		credentials := map[string]string{
			"email":    "test@example.com",
			"password": "password",
		}

		jsonData, _ := json.Marshal(credentials)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Login realizado com sucesso")
	})

	t.Run("invalid credentials", func(t *testing.T) {
		credentials := map[string]string{
			"email":    "wrong@example.com",
			"password": "wrongpassword",
		}

		jsonData, _ := json.Marshal(credentials)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "authentication_error")
	})
}

// Tests for products
func TestMockProducts(t *testing.T) {
	router := setupMockRouter()

	t.Run("list products", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Cupcake 1")
		assert.Contains(t, w.Body.String(), "Cupcake 2")
	})

	t.Run("create product valid", func(t *testing.T) {
		product := models.Product{
			Name:        "New Cupcake",
			Description: "Delicious new cupcake",
			Price:       12.50,
			ImageURL:    "http://example.com/new.jpg",
		}

		jsonData, _ := json.Marshal(product)
		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Produto criado com sucesso")
	})

	t.Run("create product invalid price", func(t *testing.T) {
		product := models.Product{
			Name:        "Invalid Cupcake",
			Description: "Cupcake with invalid price",
			Price:       -5.0,
			ImageURL:    "http://example.com/invalid.jpg",
		}

		jsonData, _ := json.Marshal(product)
		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "validation_error")
		assert.Contains(t, w.Body.String(), "price")
	})
}
