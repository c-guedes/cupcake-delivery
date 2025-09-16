package router

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "cupcake-delivery/internal/handlers"
    "cupcake-delivery/internal/middleware"
)

func Setup(db *gorm.DB) *gin.Engine {
    r := gin.Default()
    h := handlers.NewHandler(db)

    // Rotas públicas
    public := r.Group("/api")
    {
        public.POST("/register", h.Register)
        public.POST("/login", h.Login)
        public.GET("/products", h.ListProducts)
    }

    // Rotas protegidas
    protected := r.Group("/api")
    protected.Use(middleware.AuthRequired())
    {
        // Pedidos
        protected.POST("/orders", h.CreateOrder)
        protected.PUT("/orders/:id/status", h.UpdateOrderStatus)
        
        // TODO: Adicionar mais rotas protegidas
    }

    return r
}
