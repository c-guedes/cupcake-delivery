package middleware

import (
	"cupcake-delivery/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		// Remover "Bearer " do token se presente
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		// Validar token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Extrair claims do token
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Verificar se o token não expirou
			if exp, ok := claims["exp"]; ok {
				if expTime := int64(exp.(float64)); expTime < time.Now().Unix() {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
					c.Abort()
					return
				}
			}

			// Adicionar informações do usuário ao contexto
			if userID, ok := claims["user_id"]; ok {
				c.Set("user_id", uint(userID.(float64)))
			}
			if userType, ok := claims["type"]; ok {
				c.Set("type", userType.(string))
			}
		}

		c.Next()
	}
}

func TypeMiddleware(requiredType models.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de usuário não encontrado"})
			c.Abort()
			return
		}

		if userType.(string) != string(requiredType) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acesso não autorizado"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		// TODO: Validar token JWT
		// TODO: Adicionar claims ao contexto

		c.Next()
	}
}
