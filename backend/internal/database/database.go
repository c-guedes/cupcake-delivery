package database

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "cupcake-delivery/internal/models"
)

func Connect(url string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // Auto Migrate os modelos
    err = db.AutoMigrate(
        &models.User{},
        &models.Product{},
        &models.Order{},
        &models.OrderItem{},
    )
    if err != nil {
        return nil, err
    }

    return db, nil
}
