package validators

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	testCases := []struct {
		name        string
		email       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid email",
			email:       "test@example.com",
			expectError: false,
		},
		{
			name:        "Empty email",
			email:       "",
			expectError: true,
			errorMsg:    "Email é obrigatório",
		},
		{
			name:        "Invalid email format",
			email:       "invalid-email",
			expectError: true,
			errorMsg:    "Formato de email inválido",
		},
		{
			name:        "Email without domain",
			email:       "test@",
			expectError: true,
			errorMsg:    "Formato de email inválido",
		},
		{
			name:        "Email with spaces",
			email:       " test@example.com ",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEmail(tc.email)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Message != tc.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tc.errorMsg, err.Message)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %s", err.Message)
				}
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	testCases := []struct {
		name        string
		password    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid password",
			password:    "Test123",
			expectError: false,
		},
		{
			name:        "Too short password",
			password:    "Test1",
			expectError: true,
			errorMsg:    "Senha deve ter pelo menos 6 caracteres",
		},
		{
			name:        "Password without uppercase",
			password:    "test123",
			expectError: true,
			errorMsg:    "Senha deve conter pelo menos uma letra maiúscula, uma minúscula e um número",
		},
		{
			name:        "Password without lowercase",
			password:    "TEST123",
			expectError: true,
			errorMsg:    "Senha deve conter pelo menos uma letra maiúscula, uma minúscula e um número",
		},
		{
			name:        "Password without number",
			password:    "TestAbc",
			expectError: true,
			errorMsg:    "Senha deve conter pelo menos uma letra maiúscula, uma minúscula e um número",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.password)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Message != tc.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tc.errorMsg, err.Message)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %s", err.Message)
				}
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	testCases := []struct {
		name        string
		inputName   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid name",
			inputName:   "João Silva",
			expectError: false,
		},
		{
			name:        "Empty name",
			inputName:   "",
			expectError: true,
			errorMsg:    "Nome é obrigatório",
		},
		{
			name:        "Too short name",
			inputName:   "J",
			expectError: true,
			errorMsg:    "Nome deve ter pelo menos 2 caracteres",
		},
		{
			name:        "Name with spaces trimmed",
			inputName:   " João ",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateName(tc.inputName)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Message != tc.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tc.errorMsg, err.Message)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %s", err.Message)
				}
			}
		})
	}
}

func TestValidateProductPrice(t *testing.T) {
	testCases := []struct {
		name        string
		price       float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid price",
			price:       10.50,
			expectError: false,
		},
		{
			name:        "Zero price",
			price:       0,
			expectError: true,
			errorMsg:    "Preço deve ser maior que zero",
		},
		{
			name:        "Negative price",
			price:       -5.0,
			expectError: true,
			errorMsg:    "Preço deve ser maior que zero",
		},
		{
			name:        "Too high price",
			price:       15000.0,
			expectError: true,
			errorMsg:    "Preço não pode ser maior que R$ 10.000,00",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateProductPrice(tc.price)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Message != tc.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tc.errorMsg, err.Message)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %s", err.Message)
				}
			}
		})
	}
}

func TestValidateQuantity(t *testing.T) {
	testCases := []struct {
		name        string
		quantity    int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid quantity",
			quantity:    5,
			expectError: false,
		},
		{
			name:        "Zero quantity",
			quantity:    0,
			expectError: true,
			errorMsg:    "Quantidade deve ser maior que zero",
		},
		{
			name:        "Negative quantity",
			quantity:    -1,
			expectError: true,
			errorMsg:    "Quantidade deve ser maior que zero",
		},
		{
			name:        "Too high quantity",
			quantity:    150,
			expectError: true,
			errorMsg:    "Quantidade não pode ser maior que 100",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateQuantity(tc.quantity)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Message != tc.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tc.errorMsg, err.Message)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %s", err.Message)
				}
			}
		})
	}
}
