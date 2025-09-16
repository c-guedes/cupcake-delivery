package main

import (
	"cupcake-delivery/internal/config"
	"cupcake-delivery/internal/database"
	"cupcake-delivery/internal/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Carregar configurações
	cfg := config.Load()

	// Conectar ao banco de dados
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Criar usuário admin
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		Name:     "Admin",
		Email:    "admin@cupcakedelivery.com",
		Password: string(hashedPassword),
		Type:     models.AdminType,
	}
	db.FirstOrCreate(&admin, models.User{Email: admin.Email})

	// Criar produtos de exemplo
	products := []models.Product{
		{
			Name:        "Cupcake de Chocolate",
			Description: "Delicioso cupcake de chocolate com cobertura de chocolate",
			Price:       8.50,
			ImageURL:    "https://via.placeholder.com/300x300/8B4513/FFFFFF?text=Chocolate",
		},
		{
			Name:        "Cupcake de Baunilha",
			Description: "Cupcake clássico de baunilha com cobertura cremosa",
			Price:       7.00,
			ImageURL:    "https://via.placeholder.com/300x300/F5DEB3/000000?text=Baunilha",
		},
		{
			Name:        "Cupcake Red Velvet",
			Description: "Cupcake red velvet com cream cheese",
			Price:       9.00,
			ImageURL:    "https://via.placeholder.com/300x300/DC143C/FFFFFF?text=Red+Velvet",
		},
		{
			Name:        "Cupcake de Morango",
			Description: "Cupcake de morango com pedaços de morango real",
			Price:       8.00,
			ImageURL:    "https://via.placeholder.com/300x300/FFB6C1/000000?text=Morango",
		},
		{
			Name:        "Cupcake de Limão",
			Description: "Cupcake refrescante de limão com cobertura cítrica",
			Price:       7.50,
			ImageURL:    "https://via.placeholder.com/300x300/FFFF00/000000?text=Limao",
		},
		{
			Name:        "Cupcake de Coco",
			Description: "Cupcake tropical de coco com flocos de coco",
			Price:       8.50,
			ImageURL:    "https://via.placeholder.com/300x300/FFFFFF/000000?text=Coco",
		},
	}

	for _, product := range products {
		db.FirstOrCreate(&product, models.Product{Name: product.Name})
	}

	log.Println("Banco de dados populado com sucesso!")
	log.Printf("Admin criado: %s (senha: admin123)", admin.Email)
	log.Printf("%d produtos adicionados", len(products))
}
