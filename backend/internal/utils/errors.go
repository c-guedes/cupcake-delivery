package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Error       string            `json:"error"`
	Message     string            `json:"message"`
	Validations []ValidationError `json:"validations"`
}

// RespondWithError responde com um erro padronizado
func RespondWithError(c *gin.Context, statusCode int, errorType, message string, details ...map[string]string) {
	response := ErrorResponse{
		Error:   errorType,
		Message: message,
	}

	if len(details) > 0 {
		response.Details = details[0]
	}

	c.JSON(statusCode, response)
}

// RespondWithValidationError responde com erros de validação
func RespondWithValidationError(c *gin.Context, validations []ValidationError) {
	response := ValidationErrorResponse{
		Error:       "validation_error",
		Message:     "Dados inválidos fornecidos",
		Validations: validations,
	}

	c.JSON(http.StatusBadRequest, response)
}

// Common error types
const (
	ErrorTypeValidation     = "validation_error"
	ErrorTypeAuthentication = "authentication_error"
	ErrorTypeAuthorization  = "authorization_error"
	ErrorTypeNotFound       = "not_found"
	ErrorTypeConflict       = "conflict"
	ErrorTypeInternal       = "internal_error"
	ErrorTypeBadRequest     = "bad_request"
)

// Common error messages
const (
	MessageInvalidCredentials = "Email ou senha inválidos"
	MessageUnauthorized       = "Acesso não autorizado"
	MessageNotFound           = "Recurso não encontrado"
	MessageInternalError      = "Erro interno do servidor"
	MessageValidationFailed   = "Dados inválidos fornecidos"
	MessageEmailExists        = "Email já está em uso"
	MessageInvalidToken       = "Token inválido ou expirado"
)
