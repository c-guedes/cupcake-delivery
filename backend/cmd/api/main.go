package main

import (
	"cupcake-delivery/internal/config"
	"cupcake-delivery/internal/database"
	"cupcake-delivery/internal/handlers"
	"cupcake-delivery/internal/middleware"
	"cupcake-delivery/internal/models"
	"cupcake-delivery/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Carregar configurações
	cfg := config.Load()

	// Conectar ao banco de dados
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Criar services e handlers
	notificationService := services.NewNotificationService(db)

	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
	productHandler := handlers.NewProductHandler(db)
	orderHandler := handlers.NewOrderHandler(db, notificationService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Configurar rotas
	r := gin.Default()

	// Middleware de CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})

	// Rota de health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Rotas de autenticação
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	// Rotas de produtos
	products := r.Group("/products")
	{
		products.GET("", productHandler.List)
		products.GET("/:id", productHandler.Get)

		// Rotas protegidas para admin
		adminProducts := products.Group("")
		adminProducts.Use(middleware.AuthMiddleware(cfg.JWTSecret), middleware.TypeMiddleware(models.AdminType))
		{
			adminProducts.POST("", productHandler.Create)
			adminProducts.PUT("/:id", productHandler.Update)
			adminProducts.DELETE("/:id", productHandler.Delete)
		}
	}

	// Rotas de pedidos (todas precisam de autenticação)
	orders := r.Group("/orders")
	orders.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Rotas para clientes
		orders.POST("", orderHandler.Create)
		orders.GET("", orderHandler.List) // Lista filtrada por tipo de usuário

		// Rota de atualização de status (admin e entregadores)
		orders.PUT("/:id/status", orderHandler.UpdateStatus)
	}

	// Rotas de notificações (todas precisam de autenticação)
	notifications := r.Group("/notifications")
	notifications.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		notifications.GET("", notificationHandler.GetNotifications)
		notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
		notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
		notifications.PUT("/mark-all-read", notificationHandler.MarkAllAsRead)

		// Rota para criar notificação de teste (apenas para desenvolvimento)
		notifications.POST("/test", notificationHandler.CreateTestNotification)
	}

	// Iniciar servidor
	log.Printf("Servidor rodando na porta %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
