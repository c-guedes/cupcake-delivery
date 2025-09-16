package validators

import (
	"regexp"
	"strings"
	"unicode"

	"cupcake-delivery/internal/utils"
)

// ValidateEmail valida formato do email
func ValidateEmail(email string) *utils.ValidationError {
	email = strings.TrimSpace(email)
	if email == "" {
		return &utils.ValidationError{
			Field:   "email",
			Message: "Email é obrigatório",
		}
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return &utils.ValidationError{
			Field:   "email",
			Message: "Formato de email inválido",
		}
	}

	return nil
}

// ValidatePassword valida força da senha
func ValidatePassword(password string) *utils.ValidationError {
	if len(password) < 6 {
		return &utils.ValidationError{
			Field:   "password",
			Message: "Senha deve ter pelo menos 6 caracteres",
		}
	}

	if len(password) > 100 {
		return &utils.ValidationError{
			Field:   "password",
			Message: "Senha não pode ter mais de 100 caracteres",
		}
	}

	hasUpper := false
	hasLower := false
	hasNumber := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsNumber(char) {
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return &utils.ValidationError{
			Field:   "password",
			Message: "Senha deve conter pelo menos uma letra maiúscula, uma minúscula e um número",
		}
	}

	return nil
}

// ValidateName valida nome
func ValidateName(name string) *utils.ValidationError {
	name = strings.TrimSpace(name)
	if name == "" {
		return &utils.ValidationError{
			Field:   "name",
			Message: "Nome é obrigatório",
		}
	}

	if len(name) < 2 {
		return &utils.ValidationError{
			Field:   "name",
			Message: "Nome deve ter pelo menos 2 caracteres",
		}
	}

	if len(name) > 100 {
		return &utils.ValidationError{
			Field:   "name",
			Message: "Nome não pode ter mais de 100 caracteres",
		}
	}

	return nil
}

// ValidateUserType valida tipo de usuário
func ValidateUserType(userType string) *utils.ValidationError {
	validTypes := []string{"customer", "delivery", "admin"}
	for _, validType := range validTypes {
		if userType == validType {
			return nil
		}
	}

	return &utils.ValidationError{
		Field:   "type",
		Message: "Tipo de usuário deve ser 'customer', 'delivery' ou 'admin'",
	}
}

// ValidateProductName valida nome do produto
func ValidateProductName(name string) *utils.ValidationError {
	name = strings.TrimSpace(name)
	if name == "" {
		return &utils.ValidationError{
			Field:   "name",
			Message: "Nome do produto é obrigatório",
		}
	}

	if len(name) < 2 {
		return &utils.ValidationError{
			Field:   "name",
			Message: "Nome do produto deve ter pelo menos 2 caracteres",
		}
	}

	if len(name) > 200 {
		return &utils.ValidationError{
			Field:   "name",
			Message: "Nome do produto não pode ter mais de 200 caracteres",
		}
	}

	return nil
}

// ValidateProductDescription valida descrição do produto
func ValidateProductDescription(description string) *utils.ValidationError {
	description = strings.TrimSpace(description)
	if description == "" {
		return &utils.ValidationError{
			Field:   "description",
			Message: "Descrição do produto é obrigatória",
		}
	}

	if len(description) < 10 {
		return &utils.ValidationError{
			Field:   "description",
			Message: "Descrição do produto deve ter pelo menos 10 caracteres",
		}
	}

	if len(description) > 1000 {
		return &utils.ValidationError{
			Field:   "description",
			Message: "Descrição do produto não pode ter mais de 1000 caracteres",
		}
	}

	return nil
}

// ValidateProductPrice valida preço do produto
func ValidateProductPrice(price float64) *utils.ValidationError {
	if price <= 0 {
		return &utils.ValidationError{
			Field:   "price",
			Message: "Preço deve ser maior que zero",
		}
	}

	if price > 10000 {
		return &utils.ValidationError{
			Field:   "price",
			Message: "Preço não pode ser maior que R$ 10.000,00",
		}
	}

	return nil
}

// ValidateImageUrl valida URL da imagem
func ValidateImageUrl(url string) *utils.ValidationError {
	url = strings.TrimSpace(url)
	if url == "" {
		return &utils.ValidationError{
			Field:   "imageUrl",
			Message: "URL da imagem é obrigatória",
		}
	}

	urlRegex := regexp.MustCompile(`^https?://.*\.(jpg|jpeg|png|gif|webp)(\?.*)?$`)
	if !urlRegex.MatchString(url) {
		return &utils.ValidationError{
			Field:   "imageUrl",
			Message: "URL da imagem deve ser válida e ter extensão jpg, jpeg, png, gif ou webp",
		}
	}

	return nil
}

// ValidateQuantity valida quantidade
func ValidateQuantity(quantity int) *utils.ValidationError {
	if quantity <= 0 {
		return &utils.ValidationError{
			Field:   "quantity",
			Message: "Quantidade deve ser maior que zero",
		}
	}

	if quantity > 100 {
		return &utils.ValidationError{
			Field:   "quantity",
			Message: "Quantidade não pode ser maior que 100",
		}
	}

	return nil
}

// ValidateOrderStatus valida status do pedido
func ValidateOrderStatus(status string) *utils.ValidationError {
	validStatuses := []string{"pending", "preparing", "ready", "delivering", "delivered"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}

	return &utils.ValidationError{
		Field:   "status",
		Message: "Status deve ser 'pending', 'preparing', 'ready', 'delivering' ou 'delivered'",
	}
}
