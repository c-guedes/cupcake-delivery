package database

import (
	"cupcake-delivery/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(url string, autoMigrate bool) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if autoMigrate {
		err = db.AutoMigrate(
			&models.User{},
			&models.Product{},
			&models.Order{},
			&models.OrderItem{},
			&models.Notification{},
		)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
